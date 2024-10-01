package postgres

import (
	"database/sql"
	"fmt"
	"github.com/2pizzzza/TestTask/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

type Storage struct {
	Db *sql.DB
}

func New(con *config.Config) (*Storage, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		con.DBConn.Host, con.DBConn.Port, con.DBConn.Username, con.DBConn.Database, con.DBConn.Password)

	connDb, err := sql.Open("postgres", conn)

	if err != nil {
		log.Printf("Error connection db: %s", err)
		return nil, err
	}
	log.Printf("Succses connect database")

	m, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres:postgres@localhost:5432/testtask?sslmode=disable")
	if err != nil {
		log.Printf("%s", err)
	}
	if err := m.Up(); err != nil {
		log.Printf("%s", err)
	}

	return &Storage{
		Db: connDb,
	}, nil
}
