package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

const tableRepos = `CREATE TABLE IF NOT EXISTS Repos (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
        repo TEXT NOT NULL,
	  );`

const tableHistory = `CREATE TABLE IF NOT EXISTS History (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		hash TEXT NOT NULL,
        createdAt TEXT,
        repoId INTEGER NOT NULL,
        FOREIGN KEY(app) REFERENCES Apps(id)
	  );`

func CreateDatabase() {
	if _, err := os.Stat("git-history.db"); err != nil && os.IsNotExist(err) {
		file, err := os.Create("git-history.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		sqliteDatabase := Connect()
		defer sqliteDatabase.Close()
		createTable(sqliteDatabase, tableRepos)
		createTable(sqliteDatabase, tableHistory)

        for _, dir := range utils.ConfigValue.Directories {
            InsertRepo(sqliteDatabase, dir.Name)
        }

		log.Println("sqlite-database.db created")
	}
}

func Connect() *sql.DB {
	sqliteDatabase, _ := sql.Open("sqlite3", "./git-history.db") // Open the created SQLite File
	return sqliteDatabase
}

func connectExecuteAndClose(query string) {
	sqliteDatabase := Connect()
	defer sqliteDatabase.Close() // Defer Closing the database
	statement, err := sqliteDatabase.Prepare(query)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

func createTable(db *sql.DB, query string) {
	createStudentTableSQL := query
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("student table created")
}
