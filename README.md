# telegram-expenses-bot

# Add To Telegram
Add @MyXpensesBot on telegram.

# How To Self Host

## Prerequisite

1. Create a telegram bot using [BotFather](https://telegram.me/botfather).
2. Optional. If using webhook, register the domain of the server you are hosting this bot on.
```bash
curl https://api.telegram.org/bot=<token>/setWebhook?url=<domain>
```
2. A postgres database and init the database with the [init.sql](https://github.com/aattwwss/telegram-expense-bot/blob/main/scripts/init.sql).

## Run the bot
1. Clone the repo
```bash
git clone https://github.com/aattwwss/telegram-expense-bot
```
2. Create local env file
```bash
cd telegram-expense-bot
touch .env.local
```
3. Edit and save local env file with your own configurations
```bash
vim .env.local
```
4. Run the server and bot
Without docker
```bash
export .env.local
go run main.go
```
With docker
```bash
./start.sh
```

# Privacy
This bot does not store any personal information other than your telegram user id.

# Features
- [x] Sign up as a new user from new chat with bot
  - [x] How to determine a new user is using it for the first time? Use /start to check? Use the user in db to check
-[x] Add a transaction as current user
  - [x] Selection of category 
- [x] Implement expense and income command
  - [x] Enter value -> choose expense / income -> choose category
- [ ] CRUD-ish functionalities
  - [x] Delete last entry by using /undo command
  - [ ] View transactions by using /list command (UX tbd)
- [ ] ~~Parse callback amount using regex , so i can use actual currency value instead of just the float amount~~
- [ ] User preference menu (currency and timezone)
  - [ ] ~~On /start allow user to choose their timezone~~
  - [ ] Make it customisable
- [ ] Calculate transaction per month
  - [x] Return summary in text
    - [x] Triggered from /stats, default fetch from last 3 months.
    - [x] /stats month year 
    - [x] /stats month \[current year\] 
  - [ ] Return summary in chart
- [ ] Calculate by date range
  - [ ] Return summary in text
  - [ ] Return summary in chart
- [ ] Export transactions
- [ ] Set custom categories per user

# Dev / Infra 
- [ ] Add CICD from github actions to deploy to remote server 
  - [ ] ~~Build image -> deploy to docker hub -> trigger deployment in remote server~~
  - [ ] Pull from remote and build and run :(

# Misc
- [ ] Consolidate SQL to the latest schema
- [ ] Prepare docker compose for self hosted guide
- [ ] Prepare self hosted guide
