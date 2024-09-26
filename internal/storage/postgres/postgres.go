package postgres

import (
	"database/sql"
	"fmt"
	"github.com/2pizzzza/TestTask/internal/config"
	"github.com/2pizzzza/TestTask/internal/lib/logger/sl"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"os"
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

	err = ExecuteSQLFile(connDb)
	if err != nil {
		slog.Error("Failed create db", sl.Err(err))
		panic("good bye")
	}

	return &Storage{
		Db: connDb,
	}, nil
}

func ExecuteSQLFile(db *sql.DB) error {
	content, err := os.ReadFile("storage/init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(content))
	return err
}
