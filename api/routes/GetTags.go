package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/navigation"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func GetTags(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get tags")
	enableCors(&w)
	response_tags, err := json.Marshal(navigation.GetTags())
	if err != nil {
		log.Fatal(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response_tags)
}
