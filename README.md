# telegram-expenses-bot

# Text Format
- \<item\> \<amount\> [date]

# Features
- [x] Sign up as a new user from new chat with bot
    - [x] How to determine a new user is using it for the first time? Use /start to check? Use the user in db to check
- [ ] Add CICD from github actions to deploy to remote server 
  - [ ] ~~Build image -> deploy to docker hub -> trigger deployment in remote server~~
  - [ ] Pull from remote and build and run :(
-[x] Add a transaction as current user
    - [x] Selection of category 
    - [ ] ~~Input text format to record or step by step input ???~~
- [ ] Calculate transaction per month/year
  - [ ] Return summary in text
  - [ ] Return summary in chart
- [ ] Calculate by date range
  - [ ] Return summary in text
  - [ ] Return summary in chart