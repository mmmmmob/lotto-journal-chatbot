# Load env vars from repo root .env.local (optional — falls back to system environment)
-include ../../.env.local
export

MIGRATION_PATH := ./migrations
DB_URL         := postgres://$(DB_USERNAME):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable

# ── App ──────────────────────────────────────────────────────────────────────

.PHONY: run build clean

run:
	@go tool air

build:
	@mkdir -p dist
	@go build -o dist/main app/main.go

clean:
	@rm -rf tmp dist

# ── Database ─────────────────────────────────────────────────────────────────

.PHONY: db-start db-stop

db-start:
	@docker compose --env-file ../../.env.local -f ../../docker-compose.yml up -d

db-stop:
	@docker compose -f ../../docker-compose.yml down

# ── Migrations ───────────────────────────────────────────────────────────────

.PHONY: migrate-up migrate-up-one migrate-down migrate-down-all migrate-version migrate-force

migrate-up:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-up-one:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up 1

migrate-down:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down 1

migrate-down-all:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down

migrate-version:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" version

migrate-force:
	@migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" force $(N)
