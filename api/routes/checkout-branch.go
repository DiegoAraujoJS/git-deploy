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
	enableCors(&w)

	fmt.Println("checkout branch", r.URL.Query().Get("repo"), r.URL.Query().Get("commit"))
    repo := r.URL.Query().Get("repo")

	checkout_result, err := navigation.Checkout(repo, r.URL.Query().Get("commit"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while moving to reference"))
	}

    if build_err := builddeploy.Build(repo); build_err != nil {
        log.Println(build_err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while building web application"))
	}

    response := CheckoutResponse{
		Version: checkout_result.Hash().String(),
    }

    if query_error := database.InsertVersionChangeEvent(repo, checkout_result.Hash().String()); query_error != nil {
        response.Errors["db_error"] = "An error ocurred while inserting the version change to version history."
    }

	if res, err := json.Marshal(&response); err != nil {
		log.Println(err.Error())
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte("Internal error"))
	} else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(res)
    }
}
