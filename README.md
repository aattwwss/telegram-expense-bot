# telegram-expenses-bot

# Privacy
This bot does not store any personal information other than your telegram user id.

# Features
- [x] Sign up as a new user from new chat with bot
  - [x] How to determine a new user is using it for the first time? Use /start to check? Use the user in db to check
-[x] Add a transaction as current user
  - [x] Selection of category 
  - [ ] ~~Input text format to record or step by step input ???~~
- [ ] Implement expense and income command
  - [ ] Enter value -> choose expense / income -> choose category
- [ ] User preference menu (currency and timezone)
  - [ ] On /start allow user to choose their currency and timezone
  - [ ] On /currency \<curreny\> to set their currency
  - [ ] On /timezone \<timezone\> to set their timezone
- [ ] Calculate transaction per month
  - [x] Return summary in text
  - [ ] Return summary in chart
- [ ] Calculate by date range
  - [ ] Return summary in text
  - [ ] Return summary in chart
- [ ] Set currency (during /start and /currency to change)
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
