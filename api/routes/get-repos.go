package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

type Response struct {
	Repos []string
}

func GetRepos(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	content, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(content, &utils.ConfigValue)

	var repos []string
	for _, v := range utils.ConfigValue.Directories {
		repos = append(repos, v.Name)
	}
	response, err := json.Marshal(&Response{
		Repos: repos,
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
