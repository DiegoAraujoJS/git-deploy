package builddeploy

import (
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

    ActiveTimers[config.Repo] = &struct{Timer *time.Ticker; Config *AutobuildConfig}{
        Timer: new_chan,
        Config: config,
    }

    go func () {
        for t := range new_chan.C {
            fmt.Println(t, config)
            if (config.Status == ready) {fetchAndSendAction(config)}
        }
    }()

    return new_chan
}

func fetchAndSendAction(config *AutobuildConfig) error {
    config.Status = fetching
    repo := utils.GetRepository(config.Repo)

    branch, err := utils.GetBranch(repo, config.Branch)
    if err != nil {
        config.Status = down
        return err
    }
    last_commit := branch.Hash().String()

    err = utils.ForceUpdateAllBranches(repo)
    if err == git.NoErrAlreadyUpToDate {
        config.Status = ready
        return err
    }
    if err != nil {
        fmt.Println(err)
        config.Status = down
        return err
    }

    branch, _ = utils.GetBranch(repo, config.Branch)
    new_commit := branch.Hash().String()

    if last_commit != new_commit {
        CheckoutBuildInsertChan <- &Action{
            ID: GenerateActionID(),
            Repo: config.Repo,
            Hash: new_commit,
        }
    }
    config.Status = ready

    return nil
}
