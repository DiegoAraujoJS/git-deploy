package builddeploy

import (
	"bytes"

	"github.com/DiegoAraujoJS/webdev-git-server/database"
)

type Status struct {
    Finished    bool
    Moment      int8
    Stdout      *bytes.Buffer
    Stderr      *bytes.Buffer
}

type Action struct {
    ID      int
    Repo    string
    Hash    string
    Status  *Status
}

var CheckoutBuildInsertChan = make(chan *Action)

var ActionStatus = map[int]*Action{}

var StatusDict = map[int8]string {
    0: "Inactive",
    1: "Checkout branch",
    2: "Building",
    3: "Registering build",
}

func init () {
    go func () {
        for action := range CheckoutBuildInsertChan {
            go checkoutBuildInsert(action)
        }
    }()
}

const (
    inactive = iota
    checkout
    building
    registering
)

func checkoutBuildInsert(action *Action) error {
    ActionStatus[action.ID] = action
    if action.Status == nil { action.Status = &Status{} }
    for _, v := range ActionStatus {
        if v.Repo == action.Repo && v.Status.Moment != inactive && v.ID != action.ID {
            return nil
        }
    }
    if action.Status.Stdout == nil { action.Status.Stdout = &bytes.Buffer{} }
    if action.Status.Stderr == nil { action.Status.Stderr = &bytes.Buffer{} }
    action.Status.Moment = checkout
	checkout_result, err := Checkout(action.Repo, action.Hash, action.Status.Stdout)
    if err != nil {
        action.Status.Moment = inactive
        action.Status.Stderr.WriteString(err.Error())
        return err
    }
    action.Status.Moment = building
    if build_err := Build(action.Repo, action.Status.Stdout, action.Status.Stderr); build_err != nil {
        action.Status.Moment = inactive
        return build_err
    }
    action.Status.Moment = registering
    if query_error := database.InsertVersionChangeEvent(action.Repo, checkout_result.Hash().String()); query_error != nil {
        action.Status.Moment = inactive
        action.Status.Stderr.WriteString(err.Error())
        return query_error
    }
    action.Status.Moment = inactive
    action.Status.Finished = true
    // We could free the memory occupied by the buffers like below, but data may be used for further fetching.
    // action.Status.Stdout, action.Status.Stderr = nil, nil
    return nil
}


var GenerateActionID = func () func() int {
    var i int
    return func () int {
        i++
        return i
    }
}()
