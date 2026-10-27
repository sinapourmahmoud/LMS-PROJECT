package config

var Config = &config{}

type config struct {
	Application struct {
		Name        string `env:"APPLICATION_NAME"     envDefault:"LMS"`
		Environment string `env:"APPLICATION_ENV"     envDefault:"local"`
		Version     string `env:"APPLICATION_VERSION"     envDefault:"0.1"`
		PORT        string `env:"APPLICATION_PORT"     envDefault:"8080"`
	}

	Database struct {
		DB_HOST     string `env:"DB_HOST" envDefault:"localhost"`
		DB_PORT     string `env:"DB_PORT" envDefault:"5432"`
		DB_USER     string `env:"DB_USER" envDefault:"username"`
		DB_PASSWORD string `env:"DB_PASSWORD" envDefault:"password"`
		DB_NAME     string `env:"DB_NAME" envDefault:"lms"`
		DB_SSLMODE  string `env:"DB_SSLMODE" envDefault:"disable"`
	}

	Auth struct {
		JWTAccessSecret   string `env:"JWT_ACCESS_SECRET"`
		JWTRefreshSecret  string `env:"JWT_REFRESH_SECRET"`
		AccessTokenTTLMin int    `env:"ACCESS_TOKEN_TTL_MIN" envDefault:"15"`
		RefreshTokenTTLHr int    `env:"REFRESH_TOKEN_TTL_HR" envDefault:"168"`
	}
}
