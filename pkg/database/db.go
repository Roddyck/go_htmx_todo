package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDB(url string) error {
    db, err := sql.Open("sqlite3", url)

    if err != nil {
        return err
    }

    db.Exec(`CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        body TEXT NOT NULL
    )`)

    Db = db

    return nil
}
