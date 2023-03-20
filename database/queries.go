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

// Inserts a row into the database with a version change event. It returns an error if repo does not exist or if fails to insert. It also caches the repo ids.
func InsertVersionChangeEvent(repo string, hash string) error {
    repoId, ok := id_cache[repo]
    if !ok {
        var repoId int
        database := Connect()
        query := "SELECT id FROM Repos WHERE repo = '" + repo + "'"
        if err := database.QueryRow(query).Scan(&repoId); err != nil {
            if err == sql.ErrNoRows {
                fmt.Println("Failed to execute "+ query + "'", "No row that verifies condition.")
            }
            return err
        }
        id_cache[repo] = repoId
        database.Close()
    }
    err := connectExecuteAndClose("INSERT INTO History (hash, createdAt, repoId) VALUES ('" + hash + "','" + time.Now().String() + "'," + strconv.Itoa(repoId) + ")")
    if err != nil {
        log.Println()
        return err
    }
    return nil
}

func insertRepo(database *sql.DB, repo string) {
    fmt.Println("Inserting repo: " + repo)
    select_statement := "SELECT id FROM Repos WHERE repo = '" + repo + "'"
    query, err := database.Query(select_statement)
    if err != nil {
        log.Println(err.Error())
    }
    // We check for duplicates
    if !query.Next() {
        statement, err := database.Prepare("INSERT INTO Repos (repo) VALUES ('" + repo + "')")
        if err != nil {
            log.Println(err.Error())
        }
        _, err = statement.Exec()
        if err != nil {
            log.Println(err.Error())
        } else {
            log.Println("Repo " + repo + " inserted.")
        }
        return
    }
    log.Println("Repo already exists.")
}
