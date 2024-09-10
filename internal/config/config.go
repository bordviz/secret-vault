package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `yaml:"env" env-required:"true"`
	RootToken      string
	Secret         string
	MigrationsPath string `yaml:"migrations_path" env-required:"true"`
	Database       `yaml:"database" env-required:"true"`
	HTTPServer     `yaml:"http-server" env-required:"true"`
}

type Database struct {
	Host     string        `yaml:"host" env-required:"true"`
	Port     int           `yaml:"port" env-required:"true"`
	User     string        `yaml:"user" env-required:"true"`
	Password string        `yaml:"password" env-required:"true"`
	Name     string        `yaml:"name" env-required:"true"`
	SSLMode  string        `yaml:"sslmode" env-required:"true"`
	Attemps  int           `yaml:"attemps" env-required:"true"`
	Delay    time.Duration `yaml:"delay" env-required:"true"`
	Timeout  time.Duration `yaml:"timeout" env-required:"true"`
}

type HTTPServer struct {
	Port        int           `yaml:"port" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"true"`
}

func MustLoad() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("failed to load environment file, error: ", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	rootToken := os.Getenv("ROOT_TOKEN")
	secret := os.Getenv("SECRET")

	if configPath == "" {
		configPath = "config/default.yaml"
	}

	if rootToken == "" {
		log.Fatal("failed to load root token from environment file")
	}

	if secret == "" {
		log.Fatal("failed to load secret from environment file")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not exists, path: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("failed to read config file, error: ", err)
	}

	cfg.RootToken = rootToken
	cfg.Secret = secret
	return &cfg
}
