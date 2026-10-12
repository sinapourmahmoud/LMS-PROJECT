package config

var Config = &config{}


type config struct {
	Application struct {
		Name        string `env:"APPLICATION_NAME"     envDefault:"LMS"`
		Environment string `env:"APPLICATION_ENV"     envDefault:"local"`
		Version     string `env:"APPLICATION_VERSION"     envDefault:"0.1"`
	}

	Database struct {
		DB_HOST     string `env:"DB_HOST" envDefault:"localhost"`
		DB_PORT     string `env:"DB_PORT" envDefault:"5432"`
		DB_USER     string `env:"DB_USER" envDefault:"username"`
		DB_PASSWORD string `env:"DB_PASSWORD" envDefault:"password"`
		DB_NAME     string `env:"DB_NAME" envDefault:"lms"`
		DB_SSLMODE  string `env:"DB_SSLMODE" envDefault:"disable"`
	}
}
