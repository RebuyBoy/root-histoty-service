package db

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"root-histoty-service/config"
)

func ConnectPostgresDB(c config.DBConfig) (*sqlx.DB, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		"localhost", c.Port, c.User, c.Password, c.Database)

	conn, err := sqlx.Open(c.Driver, connString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
