package builddeploy

import (
	"bytes"
	"fmt"
	"time"

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
}

var ActiveTimers = map[string]*struct{
    Timer   *time.Ticker
    Config  *AutobuildConfig
}{}

var timers = map[string]*time.Timer{}

const (
    ready = iota
    fetching
    down
)

// The AddTimer function takes an action as parameter, and will add a timer that delivers the action to the checkoutBuildInsert channel continuously.
func AddTimer(config *AutobuildConfig) *time.Ticker{
    new_chan := time.NewTicker(time.Duration(config.Seconds) * time.Second)

    config.Stdout, config.Stderr = &bytes.Buffer{}, &bytes.Buffer{}
    ActiveTimers[config.Repo] = &struct{Timer *time.Ticker; Config *AutobuildConfig}{
        Timer: new_chan,
        Config: config,
    }

    stop_at := int((60 * 9) * (float32(60) / float32(config.Seconds)))
    var ticks int

    go func () {
        for t := range new_chan.C {
            fmt.Println(t, "Tick", config.Repo)
            if (config.Status == ready) {fetchAndSendAction(config)}

            ticks++
            if ticks == stop_at {
                fmt.Println("Automatically stopping timer for", config.Repo)
                new_chan.Stop()
                delete(ActiveTimers, config.Repo)
                return
            }
        }
    }()

    return new_chan
}

func fetchAndSendAction(config *AutobuildConfig) error {
    config.Status = fetching
    repo := utils.Repositories[config.Repo]

    branch, err := utils.GetBranch(repo, config.Branch)
    if err != nil {
        config.Status = down
        return err
    }
    last_commit := branch.Hash().String()

    err = utils.ForceUpdateAllBranches(repo)
    if err == git.NoErrAlreadyUpToDate {
        config.Status = ready
        config.Stdout.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - Already up to date\n")
        return err
    }
    if err != nil {
        fmt.Println(err)
        config.Status = down
        config.Stderr.WriteString(time.Now().Format("2006-01-02 15:04:05") + " - Error fetching\n" + err.Error() + "\n")
        return err
    }

    branch, _ = utils.GetBranch(repo, config.Branch)
    new_commit := branch.Hash().String()

    if last_commit != new_commit {
        CheckoutBuildInsertChan <- &Action{
            ID: GenerateActionID(),
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
