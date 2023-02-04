package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func CreateDatabase() {
	if _, err := os.Stat("git-history.db"); err != nil && os.IsNotExist(err) {
		file, err := os.Create("git-history.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		sqliteDatabase := Connect()
		defer sqliteDatabase.Close()
		createTable(sqliteDatabase)
		log.Println("sqlite-database.db created")
	}
}

func Connect() *sql.DB {
	sqliteDatabase, _ := sql.Open("sqlite3", "./git-history.db") // Open the created SQLite File
	return sqliteDatabase
}

func ConnectExecuteAndClose(query string) {
	sqliteDatabase := Connect()
	defer sqliteDatabase.Close() // Defer Closing the database
	statement, err := sqliteDatabase.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE IF NOT EXISTS History (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		hash TEXT,
        createdAt TEXT
	  );`

	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("student table created")
}
