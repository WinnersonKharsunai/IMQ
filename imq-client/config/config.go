package config

// Settings contain app environment configuration settngs
type Settings struct {
	ImqClientHost string `env:"IMQ_CLIENT_HOST" envDefault:"localhost"`
	ImqClientPort int    `env:"IMQ_CLIENT_PORT" envDefault:"80"`
	ShutdownGrace int    `env:"SHUTDOWN_GRACE" envDefault:"10"`
}
