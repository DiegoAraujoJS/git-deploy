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
    _, ok := id_cache[repo]
    fmt.Println("Inserting version change event for repo: " + repo, ok)
    if !ok {
        var repoId int
        database, _ := Connect()
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
    err := connectExecuteAndClose("INSERT INTO History (hash, createdAt, repoId) VALUES ('" + hash + "','" + time.Now().String() + "'," + strconv.Itoa(id_cache[repo]) + ")")
    if err != nil {
        log.Println("Error trying to execute INSERT statement: ", err.Error())
        return err
    }
    return nil
}

func insertRepo(database *sql.DB, repo string) {
    statement, err := database.Prepare("INSERT INTO Repos (repo) VALUES ('" + repo + "')")
    if err != nil {
        log.Println("Error trying to prepare INSERT statement: ", err.Error())
    }
    _, err = statement.Exec()
    if err != nil {
        log.Println("Error trying to execute INSERT statement: ", err.Error())
        return
    }
    log.Println("Repo " + repo + " inserted.")
}

type VersionChangeEvent struct {
    Hash string
    CreatedAt string
}

// Gets all the version change events for a given repo. The format is a struct with the following form: {hash: string, createdAt: string}. It returns an error if repo does not exist or if fails to select.
func SelectVersionChangeEvents(repo string) ([]*VersionChangeEvent, error) {
    var repoId int
    database, _ := Connect()
    query := "SELECT id FROM Repos WHERE repo = '" + repo + "'"
    if err := database.QueryRow(query).Scan(&repoId); err != nil {
        if err == sql.ErrNoRows {
            fmt.Println("Failed to execute "+ query + "'", "No row that verifies condition.")
        }
        return nil, err
    }
    rows, err := database.Query("SELECT hash, createdAt FROM History WHERE repoId = " + strconv.Itoa(repoId) + " ORDER BY createdAt DESC")
    if err != nil {
        log.Println(err.Error())
        return nil, err
    }
    var versionChangeEvents []*VersionChangeEvent
    for rows.Next() {
        var hash string
        var createdAt string
        if err := rows.Scan(&hash, &createdAt); err != nil {
            log.Println(err.Error())
            return nil, err
        }
        versionChangeEvents = append(versionChangeEvents, &VersionChangeEvent{Hash: hash, CreatedAt: createdAt})
    }
    return versionChangeEvents, nil
}
