package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP     HTTP     `yaml:"server"`
	Database Database `yaml:"database"`
	Kafka    Kafka    `yaml:"kafka"`

	CacheCapacity int `yaml:"cache_capacity" env:"CACHE_CAPACITY" env-default:"100"`
}

type HTTP struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"SERVER_PORT" env-default:"8888"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Username string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
	Scheme   string `yaml:"scheme" env:"DB_SCHEME" env-default:"public"`

	MaxConnectionAttempts int           `yaml:"maxConnectionAttempts"  env:"DB_MAX_ATTEMPTS" env-default:"5"`
	RetryDelay            time.Duration `yaml:"retryDelay" env:"DB_RETRY_DELAY" env-default:"2s"`
	ConnectionTimeout     time.Duration `yaml:"connectionTimeout" env:"DB_CONNECTION_TIMEOUT" env-default:"30s"`
}

type Kafka struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS"`
	Topic   string   `yaml:"topic" env:"KAFKA_TOPIC"`
	GroupID string   `yaml:"group_id" env:"KAFKA_GROUP_ID"`
}

func New(path string) (*Config, error) {
	var cfg Config

	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			return nil, err
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
