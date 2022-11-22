package config

type Config struct {
	DbUsername string `env:"DB_USERNAME"`
	DbPassword string `env:"DB_PASSWORD"`
	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbDatabase string `env:"DB_DATABASE"`
	DbSchema   string `env:"DB_SCHEMA"`

	TelegramApiToken string `env:"TELEGRAM_API_TOKEN"`
}
