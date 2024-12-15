package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host            string        `yaml:"Host" env:"Host"`
	Port            string        `yaml:"Port" env:"Port"`
	DatabaseUrl     string        `yaml:"DatabaseUrl" env:"DatabaseUrl"`
	JwtSecret       string        `yaml:"JwtSecret" env:"JwtSecret"`
	AccessTokenTtl  time.Duration `yaml:"AccessTokenTtl" env:"AccessTokenTtl"`
	RefreshTokenTtl time.Duration `yaml:"RefreshTokenTtl" env:"RefreshTokenTtl"`
	EmailAddres     string        `yaml:"EmailAddres" env:"EmailAddres"`
	EmailPassword   string        `yaml:"EmailPassword" env:"EmailPassword"`
	EmailHost       string        `yaml:"EmailHost" env:"EmailHost"`
	EmailCodeTtl    time.Duration `yaml:"EmailCodeTtl" env:"EmailCodeTtl"`
}

func ReadConfig(getenv func(string) string) (*Config, error) {
	var config Config
	err := cleanenv.ReadConfig("config.yml", &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
