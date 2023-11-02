package builddeploy

import (
	"bytes"
	"log"

	"github.com/DiegoAraujoJS/webdev-git-server/database"
	"github.com/DiegoAraujoJS/webdev-git-server/globals"
	"github.com/go-git/go-git/v5/plumbing"
)

type Status struct {
    Finished    bool
    Moment      int8
    Stdout      *bytes.Buffer
    Stderr      *bytes.Buffer
}

type Action struct {
    ID      int
    App     string
    Hash    plumbing.Hash
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
            go func (action *Action) {
                defer globals.GenericRecover()
                checkoutBuildInsert(action)
            }(action)
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
    if action.ID == 0 { action.ID = GenerateActionID() }
    log.Println("Received action", action.ID, action.App, action.Hash.String())
    ActionStatus[action.ID] = action
    if action.Status == nil { action.Status = &Status{} }
    for _, v := range ActionStatus {
        if v.App == action.App && v.Status.Moment != inactive && v.ID != action.ID {
            return nil
        }
    }
    if action.Status.Stdout == nil { action.Status.Stdout = &bytes.Buffer{} }
    if action.Status.Stderr == nil { action.Status.Stderr = &bytes.Buffer{} }

    action.Status.Moment = checkout
	checkout_result, err := Checkout(action)
    if err != nil {
        action.Status.Moment = inactive
        action.Status.Stderr.WriteString(err.Error())
        return err
    }

    action.Status.Moment = building
    if build_err := Build(action); build_err != nil {
        action.Status.Moment = inactive
        return build_err
    }

    action.Status.Moment = registering
    if query_error := database.InsertVersionChangeEvent(action.App, checkout_result.Hash().String()); query_error != nil {
        action.Status.Moment = inactive
        action.Status.Stderr.WriteString(err.Error())
        return query_error
    }
    action.Status.Moment = inactive
    action.Status.Finished = true
    return nil
}


var GenerateActionID = func () func() int {
    var i int
    return func () int {
        i++
        return i
    }
}()
