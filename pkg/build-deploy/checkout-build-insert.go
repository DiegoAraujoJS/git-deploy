package builddeploy

import (
	"bytes"
	"fmt"

	"github.com/DiegoAraujoJS/webdev-git-server/database"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

type Status struct {
    Errors []error
    Finished bool
    Moment int8
    Stdout *bytes.Buffer
}

type Action struct {
    ID  int
    Repo string
    Hash string
    Status *Status
}

var CheckoutBuildInsertChan chan *Action = make(chan *Action)

var ActionStatus = map[int]*Action{}

var StatusDict = map[int8]string {
    0: "Inactive",
    1: "Checkout branch",
    2: "Building",
    3: "Registering build",
    4: "Superposition of actions",
}

func init () {
    fmt.Println("setup chan")
    go func () {
        for action := range CheckoutBuildInsertChan {
            fmt.Println("action recieved", action)
            go checkoutBuildInsert(action)
        }
    }()
}

func checkoutBuildInsert(action *Action) error {
    ActionStatus[action.ID] = action
    action.Status = &Status{}
    for _, v := range ActionStatus {
        if v.Repo == action.Repo && !v.Status.Finished && v.Status.Moment != 4 && v.ID != action.ID {
            action.Status.Moment = 4
            return nil
        }
    }
    status := ActionStatus[action.ID].Status
    status.Moment = 1
	checkout_result, err := navigation.Checkout(action.Repo, action.Hash)
    if err != nil {
        status.Errors = []error{err}
        return err
    }
    status.Moment = 2
    var out bytes.Buffer
    status.Stdout = &out
    build_err := Build(action.Repo, &out)
    if build_err != nil {
        status.Errors = []error{build_err}
        return build_err
    }
    status.Moment = 3
    if query_error := database.InsertVersionChangeEvent(action.Repo, checkout_result.Hash().String()); query_error != nil {
        status.Errors = []error{query_error}
        return query_error
    }
    status.Moment = 0
    status.Finished = true
    return nil
}


var GenerateActionID = func () func() int {
    i := 1
    return func () int {
        i = i + 1
        return i
    }
}()
