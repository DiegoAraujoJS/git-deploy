package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
	"github.com/go-git/go-git/v5/plumbing"
)

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println("checkout branch")

	tag := r.URL.Query().Get("tag")

	fmt.Println("tag --> ", tag)

	checkout_result := navigation.Checkout(navigation.TagNameToHash(tag))

	fmt.Println("result --> ", checkout_result.String())

	response, err := json.Marshal(struct {
		CurrentVersion plumbing.Hash
	}{
		CurrentVersion: checkout_result,
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
