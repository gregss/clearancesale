package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type AppCfg struct {
	IsDev     bool   `env:"IS_DEV"`
	RabbitURI string `env:"RABBIT_URI" envDefault:"amqp://localhost:5672/"`
	DbURI     string `env:"DB_URI" envDefault:"postgres://localhost:5432/"`
}

func LoadConfig(cfg interface{}, fileNames ...string) {
	if len(fileNames) == 0 {
		fileNames = []string{".env", ".env.local"}
	}

	valid := []string{""}
	for _, f := range fileNames {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			continue
		}
		valid = append(valid, f)
	}
	_ = godotenv.Overload(valid...)
	_ = env.Parse(cfg)
	return
}
