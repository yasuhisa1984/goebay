package migrations

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() {
    db, err := sql.Open("sqlite3", "./blog.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    createTableSQL := `CREATE TABLE IF NOT EXISTS posts (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "title" TEXT,
        "content" TEXT
    );`

    _, err = db.Exec(createTableSQL)
    if err != nil {
        panic(err)
    }
}

