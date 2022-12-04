package config

type EnvConfig struct {
	TelegramApiToken string `env:"TELEGRAM_API_TOKEN"`

	NumRoutines int `env:"NUM_ROUTINES"`

	AppHost string `env:"APP_HOST"`
	AppPort string `env:"APP_PORT"`

	DbUsername string `env:"DB_USERNAME"`
	DbPassword string `env:"DB_PASSWORD"`
	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbDatabase string `env:"DB_DATABASE"`
	DbSchema   string `env:"DB_SCHEMA"`

	AppEnv string `env:"APP_ENV"`

	WebhookHost string `env:"WEBHOOK_HOST"`
}
