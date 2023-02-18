package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"
)

// We cache the ids for avoiding to make constant connections and queries to the db.
var id_cache = map[string]int{}

// Inserts a row into the database with a version change event.
func InsertVersionChangeEvent(repo string, hash string) error {
    repoId, ok := id_cache[repo]
    if !ok {
        var repoId int
        database := Connect()
        if err := database.QueryRow("SELECT id FROM Repos WHERE repo = '" + repo + "'").Scan(&repoId); err != nil {
            if err == sql.ErrNoRows {
                fmt.Println("Failed to execute "+"SELECT id FROM Repos WHERE repo = '" + repo + "'", "No row that verifies condition.")
                return err
            }
            return err
        }
        id_cache[repo] = repoId
        database.Close()
    }
	query := "INSERT INTO History (hash, createdAt, repoId) VALUES (" + hash + "," + time.Now().String() + "," + strconv.Itoa(repoId) + ")"
    err := connectExecuteAndClose(query)
    if err != nil {
        log.Println()
        return err
    }
    return nil
}

func insertRepo(database *sql.DB, repo string) {
    query, err := database.Query("SELECT id FROM Repos WHERE repo = '" + repo + "'")
    if err != nil {
        fmt.Println(err.Error())
    }
    // We check for duplicates
    if !query.Next() {
        statement, err := database.Prepare("INSERT INTO Repos (repo) VALUES (" + repo + ")")
        if err != nil {
            log.Println(err.Error())
        }
        statement.Exec()
    }
}
