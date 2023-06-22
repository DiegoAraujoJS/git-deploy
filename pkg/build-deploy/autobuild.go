package builddeploy

import (
	"bytes"
	"fmt"
	"time"

	"github.com/DiegoAraujoJS/webdev-git-server/globals"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
	git "github.com/go-git/go-git/v5"
)

type AutobuildConfig struct {
    Repo    string
    Seconds int
    Branch  string
    Status  int8
    Stdout  *bytes.Buffer
    Stderr  *bytes.Buffer
    LastFetch time.Time
}

type AutobuildTimers struct{
    Timer   *time.Ticker
    Config  *AutobuildConfig
}

var ActiveTimers = map[string]*AutobuildTimers{}

var timers = map[string]*time.Timer{}

const (
    ready = iota
    fetching
    down
)

// The AddTimer function takes an action as parameter, and will add a timer that delivers the action to the checkoutBuildInsert channel continuously.
func AddTimer(config *AutobuildConfig) *time.Ticker{
    new_chan := time.NewTicker(time.Duration(config.Seconds) * time.Second)

    if config.Stdout == nil { config.Stdout = &bytes.Buffer{} }
    if config.Stderr == nil { config.Stderr = &bytes.Buffer{} }
    if config.LastFetch.IsZero() { config.LastFetch = time.Now() }
    ActiveTimers[config.Repo] = &AutobuildTimers{
        Timer: new_chan,
        Config: config,
    }

    go func () {
        for range new_chan.C {
            if config.Status == ready {
                go func () {
                    defer globals.GenericRecover()
                    fetchAndSendAction(config)
                } ()
            }

            if config.LastFetch.Add(time.Duration(24) * time.Hour).Before(time.Now()) {
                fmt.Println("Automatically stopping timer for", config.Repo)
                DeleteTimer(config.Repo)
                return
            }
        }
    }()

    return new_chan
}

func DeleteTimer(repo string) {
    if timer, ok := ActiveTimers[repo]; ok {
        timer.Timer.Stop()
        delete(ActiveTimers, repo)
    }
}

func fetchAndSendAction(config *AutobuildConfig) error {
    config.Status = fetching
    repo := utils.Repositories[config.Repo]
    globals.Get_commits_rw_mutex.Lock()
    defer globals.Get_commits_rw_mutex.Unlock()

    branch, err := utils.GetBranch(repo, config.Branch)
    if err != nil {
        config.Status = down
        return err
    }
    last_commit := branch.Hash()

    err = utils.ForceUpdateAllBranches(repo)
    if err == git.NoErrAlreadyUpToDate {
        fmt.Println("Already up to date")
        config.Status = ready
        config.Stdout.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - Already up to date\n")
        return err
    }
    if err != nil {
        fmt.Println("Error fetching", err.Error())
        config.Status = down
        config.Stderr.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - Error fetching\n" + err.Error() + "\n")
        return err
    }
    // We reset the map below to re-populate the commits.

    branch, err = utils.GetBranch(repo, config.Branch)
    if err != nil {
        fmt.Println("Error getting branch", err.Error())
        config.Status = down
        config.Stderr.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - Error getting branch\n" + err.Error() + "\n")
        return err
    }
    new_commit := branch.Hash()

    if last_commit != new_commit {
        register := time.Now().Format("2006-01-02 15:04:05 ") + last_commit.String() + " --> " + new_commit.String() + "\n"
        config.LastFetch = time.Now()
        fmt.Println(register)
        config.Stdout.WriteString(register)

        CheckoutBuildInsertChan <- &Action{
            Repo: config.Repo,
            Hash: new_commit,
            Status: &Status{
                Stdout: config.Stdout,
                Stderr: config.Stderr,
            },
        }
    }
    config.Status = ready

    return nil
}
