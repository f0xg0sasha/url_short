package config

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string
	HTTPServer HTTPServer `yaml:"http_server"`
	RedisDB    RedisDB    `yaml:"redis_db"`
}

type RedisDB struct {
	Address      string `yaml:"address" env-default:"localhost:6379"`
	PoolSize     int    `yaml:"pool_size" env-default:"100"`
	MinIdleConns int    `yaml:"min_idle_conns" env-default:"30"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"deleopment"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"5s"`
}

func NewConfig() *Config {
	configPath := "config/config.yaml"

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	return cfg
}
