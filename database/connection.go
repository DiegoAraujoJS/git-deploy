package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	_ "github.com/denisenkom/go-mssqldb"
)

var sql_database *sql.DB
var no_database bool

func InitDatabase() {
    if no_database {return}

	noDb := flag.Bool("noDb", false, "a bool")
	flag.Parse()
    no_database = *noDb

    if no_database {
        log.Println("-noDb flag detected. Using git deploy server without database")
        return
    }

	db, err := Connect()
	if err != nil {
		log.Println("Error while opening database connection:", err.Error())
        return
	}

	createTableReposQuery := `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='Repos' and xtype='U') CREATE TABLE Repos (
		id INTEGER IDENTITY(1,1) PRIMARY KEY,		
        repo VARCHAR(50) NOT NULL UNIQUE
	  );`

    createTableHistoryQuery := `IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='History' and xtype='U') CREATE TABLE History (
		id INTEGER IDENTITY(1,1) PRIMARY KEY,		
		hash VARCHAR(255) NOT NULL,
        createdAt VARCHAR(255) NOT NULL,
        repoId INTEGER NOT NULL,
        CONSTRAINT fk_repoId FOREIGN KEY(repoId)
        REFERENCES Repos(id)
	  );`

	_, err = db.Exec(createTableReposQuery)
	if err != nil {
		log.Fatal("Error while creating table:", err.Error())
	}

	fmt.Println(createTableReposQuery)

	_, err = db.Exec(createTableHistoryQuery)
	if err != nil {
		log.Fatal("Error while creating table:", err.Error())
	}

	fmt.Println(createTableHistoryQuery)

    for _, dir := range utils.ConfigValue.Directories {
        insertRepo(db, dir.Name)
    }

}

func Connect() (*sql.DB, error) {
    if no_database {return nil, fmt.Errorf("-noDb flag detected; server is running disconnected from database")}
    if sql_database != nil {
        return sql_database, nil
    }
    var conn_string string
    if utils.ConfigValue.Env == "dev" {
        conn_string = "server="+utils.ConfigValue.Database.Server+";user id="+utils.ConfigValue.Database.User+";password="+ utils.ConfigValue.Database.Password+";database="+ utils.ConfigValue.Database.Name+";"
    } else {
        conn_string = "server=" + utils.ConfigValue.Database.Server + ";user id=" + ";database=" + utils.ConfigValue.Database.Name + ";trusted_connection=yes;"
    }
    fmt.Println("Attempting to connect to database with connection string:", conn_string)
    new_sql_database, err := sql.Open("sqlserver", conn_string)
    if err != nil {
        no_database = true
        log.Println("Error while opening database connection:", err.Error())
        return nil, err
    }
	err = new_sql_database.Ping()
	if err != nil {
        no_database = true
		log.Println("Error while pinging database:", err.Error())
        defer new_sql_database.Close()
        return nil, err
	}
    sql_database = new_sql_database
    fmt.Println("Successfully connected to database")
    return sql_database, err
}

func connectExecuteAndClose(query string) error {
	sql_database, err := Connect()
    if err != nil {
        return err
    }
	statement, err := sql_database.Prepare(query)
	if err != nil {
		log.Println(err.Error())
        return err
	}
    _, err = statement.Exec() // Execute SQL Statements
    fmt.Println(query)
    return err
}
