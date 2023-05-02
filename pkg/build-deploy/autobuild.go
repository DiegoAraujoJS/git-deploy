package builddeploy

import (
	"fmt"
	"time"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

type AutobuildConfig struct {
    Repo string
    Seconds int
    Branch string
    Status int8
}

var ActiveTimers = map[string]*AutobuildConfig{}

var timers = map[string]*time.Timer{}

// The AddTimer function takes an action as parameter, and will add a timer that delivers the action to the checkoutBuildInsert channel continuously.
func AddTimer(config *AutobuildConfig) *time.Ticker{
    new_chan := time.NewTicker(time.Duration(config.Seconds) * time.Second)

    ActiveTimers[config.Repo] = config

    go func () {
        for t := range new_chan.C {
            fmt.Println(t, config)
            if (config.Status == 0) {fetchAndSendAction(config)}
        }
    }()

    return new_chan
}

func fetchAndSendAction(config *AutobuildConfig) error {
    config.Status = 1
    repo := utils.GetRepository(config.Repo)
    head, _ := repo.Head()
    actual_head := head.Hash()

    fmt.Println("fetch and send action")

    branch, err := utils.GetBranch(repo, config.Branch)
    if err != nil {
        return err
    }
    fmt.Println("update branch")
    ok, err := utils.UpdateBranch(repo, branch)
    if err != nil {
        fmt.Println(err)
        return err
    }

    if actual_head == head.Hash() {return nil}

    fmt.Println("ok",ok)
    if ok {
        CheckoutBuildInsertChan <- &Action{
            ID: GenerateActionID(),
            Repo: config.Repo,
            Hash: head.Hash().String(),
        }
    }

    config.Status = 0
    return nil
}
