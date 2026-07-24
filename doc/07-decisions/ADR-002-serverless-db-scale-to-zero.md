# ADR-002: Neon DB Scale-to-Zero & External Cron Scheduling Trigger

**Date:** 2026-07-24
**Status:** Accepted
**Proposed by:** AI Pair Programming — 2026-07-24
**Accepted by:** Project owner — 2026-07-24
**Related Tasks:** T-024 (Neon DB Compute Leak & Optimization)

---

## Context

Lotto Journal uses **Neon Postgres** (Serverless Database) and **Fly.io** for compute. 

Under the default configuration:
1. Fly.io kept at least 1 machine running at all times (`min_machines_running = 1`).
2. An active HTTP health check (`[[http_service.checks]]`) pinged the server `/health` route every 15 seconds.
3. An in-process cron scheduler (`github.com/robfig/cron/v3`) was running inside the Go application to handle daily draw schedule sync and result verification.
4. The database connection pool did not aggressively close idle connections.

This combination of settings created a **compute leak** in Neon Postgres:
* Neon counts active Compute Units (CU) as long as there is at least one active connection (even if Postgres query count is 0).
* Since the Fly.io VM was kept running 24/7, the Go application kept persistent idle connections open.
* As a result, Neon DB was active 24/7, exhausting 100% of the Free Compute monthly quota (100 CU-hours) within days.

---

## Options Considered

### Option A — Keep In-Process Cron Scheduler (Always-On VM)

Keep the VM always running, but optimize database pool connections to close aggressively after idle periods (e.g. `SetMaxIdleConns(0)`).

* **Pros:**
  * No external scheduler dependency (GitHub Actions).
  * Keeps response times low since the VM is always warm.
* **Cons:**
  * Exposes us to risks if database connections leak for other reasons.
  * Fly.io hosting costs are higher (always-on compute resources).
  * Does not fully leverage serverless pay-as-you-go architecture.

---

### Option B — Enable VM Auto-Stop & Migrate Crons to GitHub Actions ✅ CHOSEN

Modify the system to let both the compute (Fly.io) and database (Neon Postgres) sleep when there is no traffic.
1. Allow Fly.io VM to autosleep (`min_machines_running = 0`, comment out the health check).
2. Optimize DB connection pooling in Go (`SetMaxIdleConns(0)`, `SetMaxOpenConns(5)`, `SetConnMaxLifetime(3 * time.Minute)`) so that when the app is active, connections are quickly returned and terminated.
3. Replace the in-process cron scheduler with **GitHub Actions schedule workflows**. These workflows send a webhook request with a custom token header (`Authorization: Bearer <CRON_SECRET>`) to wake up the app and trigger jobs over HTTP.

* **Pros:**
  * **Zero connection leak**: When there is no HTTP traffic, Fly.io stops the VM. The database connections drop to 0, allowing Neon DB to go into suspend mode (scale-to-zero).
  * **Cost-effective**: Zero consumption on Neon when inactive, and lower resource hours on Fly.io.
  * **Secured triggers**: The trigger endpoints are secured with a token, preventing unauthorized invocation.
  * **Timezone precision**: Workflows are scheduled precisely in UTC matching Bangkok Time (ICT).
* **Cons:**
  * Cold start latency: The first user message after a period of inactivity may experience a slight delay (2–5 seconds) while Fly.io wakes up the VM and establishes the DB connection.
  * Dependency on GitHub Actions scheduling (GitHub Actions crons can occasionally experience delays).

---

## Decision

**Option B** was chosen.

**Rationale:**
- Neon's Free Compute tier is critical to keep free-of-charge during development and initial launch. Scale-to-zero is the only way to stay within 100 CU-hours/month.
- A LINE Bot has extremely sporadic traffic (mostly on the 1st and 16th of the month, or occasionally when users add ticket numbers). Cold start latencies of 3-4 seconds are fully acceptable for a chatbot environment compared to a real-time web application.
- Exposed job endpoints are secure and simple.

## Consequences

- The in-process cron scheduler will not start in production (`APP_ENV == "production"`).
- We must configure a secret token `CRON_SECRET` on both GitHub Secrets and Fly.io Secrets to authenticate the cron runners.
- The Fly VM will shut down after a few minutes of inactivity. When a webhook event (LINE) arrives, Fly will automatically restart the machine to handle it.
