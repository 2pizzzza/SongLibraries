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

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("PASSWORD")

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
		},
		HttpConn: HttpConfig{
			HttpPort: httpPort,
		},
	}, nil
}
