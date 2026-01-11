package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL" env-default:"INFO"` 
	Server Server `yaml:"server"`
	PostgresDB PostgresDB `yaml:"postgres_db"`
}

type Server struct {
	Port string `yaml:"port" env:"PORT" env-default:":8080"` 
	Host string `yaml:"host" env:"HOST" env-default:"localhost"` 
}

type PostgresDB struct {
	User string `yaml:"user" env:"USER_DB" env-default:"postgres"`
	Pass string `yaml:"pass" env:"PASS_DB" env-default:"postgres"`
	Host string `yaml:"host" env:"HOST_DB" env-default:"localhost"`
	Port int `yaml:"port" env:"PORT_DB" env-default:"5432"`
	Name string `yaml:"name" env:"NAME_DB" env-default:"postgres"`
	MaxAttempts int `yaml:"max_attempts" env:"MAX_ATEMPT" env-default:"5"`
}

func(p *PostgresDB) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", p.User, p.Pass, p.Host, p.Port, p.Name)
}

func Init(path string) (*Config, error) {
	var cfg Config
    if err := cleanenv.ReadConfig(path, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}