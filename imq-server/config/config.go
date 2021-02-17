package config

// Settings contain app environment configuration settngs
type Settings struct {
	ImqServerHost string `env:"IMQ_SERVER_HOST" envDefault:"localhost"`
	ImqServerPort int    `env:"IMQ_SERVER_PORT" envDefault:"80"`
	ShutdownGrace int    `env:"SHUTDOWN_GRACE" envDefault:"10"`

	DbUserName string `env:"DB_USERNAME" envDefault:"root"`
	DbPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DbHost     string `env:"DB_HOST" envDefault:"localhost"`
	DbPort     int    `env:"DB_PORT" envDefault:"3308"`
	DbName     string `env:"DB_NAME" envDefault:"messagingQueueDev"`
}
