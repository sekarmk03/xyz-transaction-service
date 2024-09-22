package config

import (
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	ServiceName       string `env:"SERVICE_NAME,default=xyz-grpc"`
	Port              Port
	MySQL             MySQL
	JWT               JWTConfig
	ClientURL         ClientURL
}

type Port struct {
	GRPC string `env:"PORT_GRPC,default=8081"`
}

type MySQL struct {
	Host     string `env:"MYSQL_HOST,default=localhost"`
	Port     string `env:"MYSQL_PORT,default=3306"`
	User     string `env:"MYSQL_USER"`
	Password string `env:"MYSQL_PASSWORD"`
	Name     string `env:"MYSQL_NAME"`
}

type JWTConfig struct {
	JwtSecretKey  string        `env:"JWT_SECRET_KEY"`
	TokenDuration time.Duration `env:"JWT_DURATION,default=30m"`
}

type ClientURL struct {
	Consumer string `env:"CLIENT_URL_CONSUMER"`
}

func NewConfig(env string) (*Config, error) {
	_ = godotenv.Load(env)

	var config Config
	if err := envdecode.Decode((&config)); err != nil {
		return nil, errors.Wrap(err, "ERROR: [NewConfig] Failed to decode environment variables")
	}

	return &config, nil
}
