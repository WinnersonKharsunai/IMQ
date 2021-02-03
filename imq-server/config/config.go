package config

// Settings contain app environment configuration settngs
type Settings struct {
	ImqServerHost string `envconfig:"IMQ_SERVER_HOST"`
	ImqServerPort int    `envconfig:"IMQ_SERVER_PORT"`

	DbUserName string `envconfig:"DB_USERNAME"`
	DbPassword string `envconfig:"DB_PASSWORD"`
	DbHost     string `envconfig:"DB_HOST"`
	DbPort     int    `envconfig:"DB_PORT"`
	DbName     string `envconfig:"DB_NAME"`
}
