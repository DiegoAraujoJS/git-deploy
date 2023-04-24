package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/database"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

type CheckoutResponse struct {
    Errors map[string]string
	Version string
}

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
    repo := r.URL.Query().Get("repo")
	fmt.Println("checkout branch", repo, r.URL.Query().Get("commit"))

	checkout_result, err := navigation.Checkout(repo, r.URL.Query().Get("commit"))
	if err != nil {
        fmt.Println("Checkout error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("🔴 " + err.Error()))
        return
	}

    if build_err := builddeploy.Build(repo); build_err != nil {
        log.Println("Build error", build_err.Error())
        w.Header().Set("Content-Type", "text")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("🔴 " + build_err.Error()))
        return
	}

    response := CheckoutResponse{
		Version: checkout_result.Hash().String(),
    }

    fmt.Println("checkout result", checkout_result.Hash().String())

    if query_error := database.InsertVersionChangeEvent(repo, checkout_result.Hash().String()); query_error != nil {
        log.Println("Error while inserting version change event to version history", query_error.Error())
        w.Header().Set("Content-Type", "text")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("🔴 " + query_error.Error()))
    }

	if res, err := json.Marshal(&response); err != nil {
		log.Println(err.Error())
	    w.Header().Set("Content-Type", "text")
	    w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
	} else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(res)
    }
}
