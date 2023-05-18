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
// We cache the versions for avoiding to make constant connections and queries to the db.
var version_cache = map[string][]*VersionChangeEvent{}

func getRepoId (repo string) (int, error) {
    _, ok := id_cache[repo]
    if !ok {
        var repoId int
        database, _ := Connect()
        query := "SELECT id FROM Repos WHERE repo = '" + repo + "'"
        if err := database.QueryRow(query).Scan(&repoId); err != nil {
            if err == sql.ErrNoRows {
                fmt.Println("Failed to execute "+ query + "'", "No row that verifies condition.")
            }
            return 0, err
        }
        id_cache[repo] = repoId
    }
    return id_cache[repo], nil
}
// Inserts a row into the database with a version change event. 
// 
// It returns an error if repo does not exist or if fails to insert. It also caches the repo ids.
func InsertVersionChangeEvent(repo string, hash string) error {
    repo_id, err := getRepoId(repo)
    if err != nil {
        return err
    }
    err = connectExecuteAndClose("INSERT INTO History (hash, createdAt, repoId) VALUES ('" + hash + "','" + time.Now().String() + "'," + strconv.Itoa(repo_id) + ")")
    if err != nil {
        log.Println("Error trying to execute INSERT statement: ", err.Error())
        return err
    }
    delete(version_cache, repo)
    return nil
}

func insertRepo(database *sql.DB, repo string) {
    query := "INSERT INTO Repos (repo) VALUES ('" + repo + "')"
    statement, err := database.Prepare(query)
    if err != nil {
        log.Println("Error trying to prepare INSERT statement: ", err.Error())
    }
    _, err = statement.Exec()
    fmt.Println(query)
    if err != nil {
        log.Println("Error trying to execute INSERT statement: ", err.Error())
        return
    }
    log.Println("Repo " + repo + " inserted.")
}

type VersionChangeEvent struct {
    Hash        string
    CreatedAt   string
}

// Gets all the version change events for a given repo.
//
// The format is a struct with the following form: {Hash: string, CreatedAt: string}. It returns an error if repo does not exist or if fails to select.
func SelectVersionChangeEvents(repo string) ([]*VersionChangeEvent, error) {
    repoId, err := getRepoId(repo)
    if err != nil {
        return nil, err
    }
    if response, ok := version_cache[repo]; ok {
        return response, nil
    }
    database, _ := Connect()
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
            continue
        }
        versionChangeEvents = append(versionChangeEvents, &VersionChangeEvent{Hash: hash, CreatedAt: createdAt})
    }
    version_cache[repo] = versionChangeEvents
    return versionChangeEvents, nil
}
