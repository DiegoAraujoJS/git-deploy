package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/build-deploy"
	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

type CheckoutResponse struct {
	Version string
}

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println("checkout branch", r.URL.Query().Get("commit"))

	checkout_result, err := navigation.Checkout(r.URL.Query().Get("commit"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while moving to reference"))
	}

	build_err := builddeploy.Build()

	if build_err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while building web application"))
	}

    deploy_err := builddeploy.DeployIIS()

	if deploy_err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while deploying web application to IIS"))
	}

	response, err := json.Marshal(&CheckoutResponse{
		Version: checkout_result.Hash().String(),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
