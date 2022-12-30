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
- [x] Add a transaction as current user
- [x] Selection of category when adding transaction
- [x] Delete last entry by using /undo command
- [X] Calculate transaction per month
- [x] Triggered from /stats, default fetch from current month.
- [x] /stats [month] [year]
- [ ] View transactions by using /list command (UX tbd)
- [ ] Allow user to change timezone. (default Asia/Singapore)
- [ ] Allow user to change currency. (default SGD)
- [ ] Export transactions to file

# Dev / Infra 
- [ ] Fix image deployed on github container repository not being reached by telegram server

# Misc
- [x] Consolidate SQL to the latest schema
- [x] Prepare docker compose for self hosted guide
- [x] Prepare self hosted guide
