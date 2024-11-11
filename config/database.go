package config

import (
	"database/sql"
	"fmt"
	"log"
	"sqlc-test/env"

	"github.com/samber/do"
)

type DbConn struct {
	*sql.DB
}

func NewDbConnection(i *do.Injector) (*DbConn, error) {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.GetAsString("DB_USER", "postgres"),
		env.GetAsString("DB_PASSWORD", "mysecretpassword"),
		env.GetAsString("DB_HOST", "localhost"),
		env.GetAsInt("DB_PORT", 5432),
		env.GetAsString("DB_NAME", "postgres"),
	)

	// Open the database
	database, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}

	// Connectivity check
	if err := database.Ping(); err != nil {
		log.Fatalln("Error from database ping:", err)
		return nil, err
	}

	return &DbConn{database}, nil
}
