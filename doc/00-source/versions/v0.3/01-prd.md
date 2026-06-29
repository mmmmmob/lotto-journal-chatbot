# Product Requirements Document — Lotto Journal (LINE-based)

**Version:** v0.3
**Date:** 2026-06-29
**Status:** Approved
**Based on:** ADR-001 (Accepted) and PRD v0.3 Open Questions
**Supersedes:** v0.2 LINE-based PRD

---

## 1. Product Goal

Lotto Journal is a **LINE chatbot service** that lets Thai lottery players:

1. Record the lottery ticket numbers they own by sending them through LINE chat.
2. Be automatically notified via LINE when any of their tickets win a prize.
3. Interact in either **Thai (TH)** or **English (EN)**, with automatic locale detection and manual toggle options.

Users interact entirely through LINE. There is no web app or mobile app to install.

---

## 2. Target Users

Thai and foreign lottery players in Thailand who:

- Already have a LINE account.
- Purchase Thai Government Lottery (สลากกินแบ่งรัฐบาล) tickets regularly.
- Want a convenient way to track tickets and get notified of winnings in their preferred language (Thai or English).

---

## 3. User Flows

### 3.1 User Onboarding (First Interaction)

**Trigger:** User adds the Lotto Journal LINE Official Account as a friend.

**Steps:**

1. LINE delivers a `follow` event to the backend webhook.
2. Backend checks if this `line_user_id` exists in the `users` table.
3. If new user:
   - Create a `users` record.
   - Detect the user's language by querying the LINE Profile API (`profile.Language`).
   - If the profile language is `"th"`, save `language = 'th'`. Otherwise, default to `language = 'en'`.
   - Send the **welcome message** in the detected language, including a direction to switch languages:
     - **English Welcome (if language = 'en'):**
       > "Welcome to Lotto Journal! Send me your ticket numbers (e.g., `123456` or `123456x2` / `456 789`). Type `list` to see your registered tickets.
       > 
       > (Type `ไทย` to switch to Thai language)"
     - **Thai Welcome (if language = 'th'):**
       > "👋 สวัสดีคุณ {displayName}!
       > 
       > 🎟️ ยินดีต้อนรับสู่ Lotto Journal!
       > ตัวอย่าง: 123456 หรือ 456
       > ส่งหลายเลขได้ เช่น 123456 789012
       > ระบุจำนวนตั๋วด้วย x เช่น 123456x2
       > 
       > 📝 หากต้องการดูสลากที่บันทึกไว้ พิมพ์ 'โพย'
       > 
       > (พิมพ์ `english` เพื่อเปลี่ยนเป็นภาษาอังกฤษ)"
4. If returning user (already exists in the database with status `inactive`):
   - Reactivate the user (`status = 'active'`).
   - Send the welcome message in their **previously saved language**.
   - **Do not include** the language switch direction in the welcome message (since they have already onboarded).

---

### 3.2 Ticket Submission

**Trigger:** User sends a message containing ticket numbers.

**Message input format:**
- Plain numbers (e.g. `123456` or `456`).
- Multiple numbers separated by spaces or commas (e.g., `123456, 654321` or `456 789`).
- Quantities indicated with `xN` suffix (e.g. `123456x2`).
- Bilingual commands are supported at all times:
  - **List tickets:** `โพย` or `list` or `tickets`.
  - **Switch to Thai:** `ไทย` or `thai`.
  - **Switch to English:** `english` or `en`.

**Steps:**
1. System receives the webhook `message` event.
2. Checks for commands:
   - If `ไทย` or `thai` → Update user language to `th`. Reply in Thai confirming the language switch, and attach Thai **Quick Replies** (`["โพย", "เพิ่ม", "แจ้งเตือน"]`).
   - If `english` or `en` → Update user language to `en`. Reply in English confirming the language switch, and attach English **Quick Replies** (`["list", "add", "notify"]`).
   - If a list command (`โพย`, `list`, `tickets`) → Retrieve and display tickets for the upcoming draw in the user's set language.
   - If `เพิ่ม` or `add` → Reply with instructions on how to add tickets in the user's set language.
   - If `แจ้งเตือน` or `notify` → Reply with details on how automatic win/loss notifications work in the user's set language.
3. If not a command, parse ticket numbers using `ParseTicketInput` and persist them. Reply with a confirmation message in the user's set language.

---

### 3.3 Automatic Win Notification (Cronjob Flow)

**Trigger:** Scheduled cronjob fires on lottery draw days (1st and 16th of each month), after official results are published.

**Steps:**
1. Fetch results from GLO API.
2. Insert results and verify the draw.
3. Compare results against all unchecked tickets.
4. Insert `user_winnings` for winning tickets.
5. Send a **LINE push message** to each user in their set language:
   - **English (Winning):** "Congratulations! Your ticket {number} won {prize_category} ({prize_amount} Baht)."
   - **English (Non-Winning):** "Unfortunately, your ticket {number} did not win this time. Better luck next time!"
   - **Thai (Winning):** "ยินดีด้วย! สลากเลข {number} ของคุณ ถูกรางวัล {prize_category} ({prize_amount} บาท)"
   - **Thai (Non-Winning):** "ขออภัย สลากเลข {number} ของคุณ ไม่ถูกรางวัลในงวดนี้ พยายามใหม่อีกครั้งนะ!"
6. Mark processed tickets as `is_checked = true`.

---

### 3.4 Quick Replies UI

Quick Replies will be displayed as suggestions at the bottom of the chat interface whenever the user interacts with commands or changes settings.

- **Thai Quick Replies:**
  - `โพย`: Triggers the ticket list command.
  - `เพิ่ม`: Sends an instruction on how to add tickets.
  - `แจ้งเตือน`: Sends information on automated win notifications.
- **English Quick Replies:**
  - `list`: Triggers the ticket list command.
  - `add`: Sends an instruction on how to add tickets.
  - `notify`: Sends information on automated win notifications.

---

## 4. System Architecture

No major architectural changes are introduced. The core Fiber app receives the webhook events, queries the database, and communicates with the LINE Messaging API. The translation dictionaries are loaded dynamically in memory.

---

## 5. Data Model

### 5.1 Migration 000006 — User Language Settings

Add a `language` column to the `users` table.

**`users` table:**

| Column     | Type        | Constraints                | Default |
| ---------- | ----------- | -------------------------- | ------- |
| `language` | `varchar(10)`| `NOT NULL`                 | `'en'`  |

---

## 6. Message Localization Dictionaries

### 6.1 English (EN)

- **Welcome (First-time):**
  > "Welcome to Lotto Journal! Send me your ticket numbers (e.g., `123456` or `123456x2` / `456 789`). Type `list` to see your registered tickets.
  > 
  > (Type `ไทย` to switch to Thai language)"
- **Welcome (Returning):**
  > "Welcome back to Lotto Journal! Send me your ticket numbers to get started, or type `list` to see your registered tickets."
- **Submit Confirm:**
  > "Successfully recorded ticket(s) ✅:\n{tickets_list}"
- **Submit Invalid:**
  > "No valid ticket numbers found ❌. Please send 3-digit or 6-digit numbers only."
- **Submit Mixed (Valid + Invalid):**
  > "Successfully recorded ticket(s) ✅:\n{tickets_list}\n\nInvalid numbers (skipped): {invalid_list}"
- **List Header:**
  > "Your tickets for the upcoming draw ({date}) 📝:"
- **List Empty:**
  > "You have no tickets registered for this draw."
- **Language Switched:**
  > "Language changed to English 🇺🇸"
- **Add Help:**
  > "🎟️ To record ticket numbers, simply send them in chat.
  > 
  > - Single number: `123456` or `456`
  > - Multiple numbers: `123456, 654321` or `123 456`
  > - Quantity: Append `xN` (e.g. `123456x2` for 2 tickets)"
- **Notify Help:**
  > "🔔 You will be automatically notified here shortly after the lottery draw finishes (typically around 16:00 Bangkok time on the 1st and 16th of each month)."
- **Win Notification:**
  > "Congratulations! Your ticket {number} won {prize_category} ({prize_amount} Baht)."
- **Loss Notification:**
  > "Unfortunately, your ticket {number} did not win this time. Better luck next time!"

### 6.2 Thai (TH)

- **Welcome (First-time):**
  > "👋 สวัสดีคุณ {displayName}!
  > 
  > 🎟️ ยินดีต้อนรับสู่ Lotto Journal!
  > พิมพ์เลขสลากที่คุณซื้อไว้เพื่อรอตรวจผลอัตโนมัติได้เลย
  > ตัวอย่าง: 123456 หรือ 456
  > ส่งหลายเลขได้ เช่น 123456 789012
  > ระบุจำนวนตั๋วด้วย x เช่น 123456x2
  > 
  > 📝 หากต้องการดูสลากที่บันทึกไว้ พิมพ์ 'โพย'
  > 
  > (พิมพ์ `english` เพื่อเปลี่ยนเป็นภาษาอังกฤษ)"
- **Welcome (Returning):**
  > "👋 ยินดีต้อนรับกลับสู่ Lotto Journal! ส่งเลขสลากของคุณเพื่อเริ่มต้นบันทึกได้เลย หรือพิมพ์ 'โพย' เพื่อดูสลากที่บันทึกไว้"
- **Submit Confirm:**
  > "บันทึกสลากเรียบร้อย ✅:\n{tickets_list}"
- **Submit Invalid:**
  > "ไม่พบเลขสลากที่ถูกต้อง ❌\nกรุณาส่งเลข 3 หรือ 6 หลักเท่านั้น"
- **Submit Mixed (Valid + Invalid):**
  > "บันทึกสลากเรียบร้อย ✅:\n{tickets_list}\n\nเลขที่ไม่ถูกต้อง (ข้ามไป): {invalid_list}"
- **List Header:**
  > "สลากที่คุณบันทึกไว้ในงวดนี้ ({date}) 📝:"
- **List Empty:**
  > "คุณยังไม่ได้บันทึกสลากในงวดนี้"
- **Language Switched:**
  > "เปลี่ยนภาษาเป็นภาษาไทยเรียบร้อย 🇹🇭"
- **Add Help:**
  > "🎟️ คุณสามารถบันทึกสลากได้ง่ายๆ โดยพิมพ์ส่งเข้ามาในแชท:
  > 
  > - เลขตัวเดียว: `123456` หรือ `456`
  > - ส่งหลายเลขพร้อมกัน: `123456, 654321` หรือ `123 456`
  > - ระบุจำนวน: ต่อท้ายด้วย `xN` เช่น `123456x2` (บันทึก 2 ใบ)"
- **Notify Help:**
  > "🔔 ระบบจะส่งผลรางวัลให้คุณทราบโดยอัตโนมัติ ทันทีหลังจากประกาศผลรางวัลเสร็จสิ้น (ปกติประมาณ 16:00 น. ของวันที่ 1 และ 16 ของทุกเดือน)"
- **Win Notification:**
  > "ยินดีด้วย! สลากเลข {number} ของคุณ ถูกรางวัล {prize_category} ({prize_amount} บาท) 🎉"
- **Loss Notification:**
  > "ขออภัย สลากเลข {number} ของคุณ ไม่ถูกรางวัลในงวดนี้ พยายามใหม่อีกครั้งนะ! ✌️"
