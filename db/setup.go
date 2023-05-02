package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var ErrNoMatch = fmt.Errorf("no matching record")

type Database struct {
	Conn *sql.DB
}

func Initialize(name string) (Database, error) {
	db := Database{}
	conn, err := sql.Open("sqlite3", name)

	if err != nil {
		return db, err
	}

	db.Conn = conn

	sqlStm := `
	create table if not exists users (id integer not null primary key autoincrement, username text, password text);
	`
	_, err = db.Conn.Exec(sqlStm)
	if err != nil {
		return db, err
	}
	log.Println("Database connection established")

	return db, nil
}
