package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Server struct {
	ServerHost string `env:"HOST" env-default:"0.0.0.0"`
	ServerPort int    `env:"PORT" env-default:"8000"`
}

type Database struct {
	DbHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	DbPort     int    `env:"POSTGRES_PORT" env-default:"5432"`
	DbUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	DbPassword string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	DbName     string `env:"POSTGRES_NAME" env-default:"postgres"`
}

func (d *Database) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.DbUser, d.DbPassword, d.DbHost, d.DbPort, d.DbName)
}

type Config struct {
	Server
	Database
	JwtSecret   string `env:"JWT_SECRET"`
	SecretToken string `env:"SECRET_TOKEN"`
}

func GetConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
