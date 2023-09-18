package config

type App struct {
	Tracing struct {
		HONEYCOMB_SERVICE_NAME string `env:"HONEYCOMB_SERVICE_NAME,default=test"`
		HONEYCOMB_WRITEKEY     string `env:"HONEYCOMB_WRITEKEY"`
		HONEYCOMB_DATASET      string `env:"HONEYCOMB_DATASET"`
	}
	DEV  string `env:"IS_DEV, default=true"`
	PORT string `env:"PORT, default=8080"`
}
