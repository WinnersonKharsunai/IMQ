package config

// Settings contain app environment configuration settngs
type Settings struct {
	ImqClientHost string `envconfig:"IMQ_CLIENT_HOST"`
	ImqClientPort int    `envconfig:"IMQ_CLIENT_PORT"`
}
