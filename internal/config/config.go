package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	DBConn   DatabaseConfig
	HttpConn HttpConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	DbUrl    string
}

type HttpConfig struct {
	HttpPort int
}

func NewConfig() (db *Config, err error) {
	err = godotenv.Load()

	if err != nil {
		fmt.Println("Error is occurred  on .env file please check", err)
		return nil, err
	}

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")
	dburl := os.Getenv("DB_URL")

	env := os.Getenv("ENV")

	httpPort, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))

	return &Config{
		Env: env,
		DBConn: DatabaseConfig{
			Host:     host,
			Port:     port,
			Database: dbname,
			Username: pass,
			Password: pass,
			DbUrl:    dburl,
		},
		HttpConn: HttpConfig{
			HttpPort: httpPort,
		},
	}, nil
}
