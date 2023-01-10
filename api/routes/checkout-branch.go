package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

type CheckoutResponse struct {
	Version string
}

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println("checkout branch")

	tag := r.URL.Query().Get("tag")

	checkout_result, err := navigation.Checkout(tag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while moving to reference"))
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
