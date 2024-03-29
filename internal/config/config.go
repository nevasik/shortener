package config

import (
	"Url-shortner/internal/lib/logger/handlers/slogpretty"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	TimeOut     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// для генерации переменной alias, если ее не указал пользователь
const AliasLength = 6

// читает файл с конфигом и создаёт и заполняет объект с конфигом
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		slog.Warn("config file does not exist:", configPath)
	}

	var cfg Config

	err1 := cleanenv.ReadConfig(configPath, &cfg)
	if err1 != nil {
		slog.Warn("cannot read config: ", err1)
	}

	return &cfg
}

func SetupLogger() *slog.Logger {
	var log *slog.Logger

	log = setupPrettySlog()
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
