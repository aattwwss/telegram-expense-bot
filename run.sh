docker build --tag telegram-expense-bot .
docker run -e APP_ENV=prod --env-file ./.env.local telegram-expense-bot