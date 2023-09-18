package config

type DB struct {
	URL         string `env:"DB_URL, required"`
	PORT        string `env:"DB_PORT, required"`
	HOST        string `env:"DB_HOST, required"`
	USER        string `env:"DB_USER, required"`
	SSLMODE     string `env:"DB_SSL_MODE, required"`
	PASS        string `env:"DB_PASS, required"`
	NAME        string `env:"DB_NAME, required"`
	MAXIDLE     string `env:"DB_MAX_IDLE_CONN,default=10"`
	MAXOPENCONN string `env:"DB_MAX_OPEN_CONN,default=100"`
	MAXLIFETIME string `env:"DB_MAX_LIFE_TIME,default=1"`
}
