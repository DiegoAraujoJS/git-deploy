package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var password string

func init () {
    // load config values
    var config = struct{
        Credentials struct {
            Password string
        }
    }{}
	content, _ := ioutil.ReadFile("../../main/config.json")

	_ = json.Unmarshal(content, &config)
    password = config.Credentials.Password
}

func TestGetRepos(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:3001" +"/getRepos", nil)
    req.Header.Add("Authorization", password)
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if response.StatusCode != http.StatusOK {
        t.Error("Status code not ok")
    }
    // Check for response length, and type. Type should be a list of strings, length should be at least 1.
    if response.Header.Get("Content-Type") != "application/json" {
        t.Error("Response type not ok")
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Error(err)
    }
    if len(body) < 1 {
        t.Error("Response body length not ok. Should be at least 1.", len(body))
    }
}

func TestGetRepoHistory(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:3001" +"/repoHistory?repo=PWA-VUEJS", nil)
    req.Header.Add("Authorization", password)
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if response.StatusCode != http.StatusOK {
        t.Error("Status code not ok")
    }
    if response.Header.Get("Content-Type") != "application/json" {
        t.Error("Response type not ok")
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Error(err)
    }
    // Check for the length and the format of the response. Should be a list of objects {Hash: string, CreatedAt: string, Commit: {Message: string, Author: string}}
    if len(body) < 1 {
        t.Error("Response body length not ok. Should be at least 1.", len(body))
    }
    var response_list []map[string]interface{}
    err = json.Unmarshal(body, &response_list)
    if err != nil {
        t.Error(err)
    }
    if len(response_list) < 1 {
        t.Error("Response list length not ok. Should be at least 1.", len(response_list))
    }
    for _, item := range response_list {
        if _, ok := item["Hash"]; !ok {
            t.Error("Hash not found in response list item.")
        }
        if _, ok := item["CreatedAt"]; !ok {
            t.Error("CreatedAt not found in response list item.")
        }
        if _, ok := item["Commit"]; !ok {
            t.Error("Commit not found in response list item.")
        }
        commit := item["Commit"].(map[string]interface{})
        if _, ok := commit["Message"]; !ok {
            t.Error("Message not found in response list item commit.")
        }
        if _, ok := commit["Author"]; !ok {
            t.Error("Author not found in response list item commit.")
        }
    }
}

func TestGetTags(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:3001" +"/getTags?repo=PWA-VUEJS", nil)
    req.Header.Add("Authorization", password)
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if response.StatusCode != http.StatusOK {
        t.Error("Status code not ok")
    }
    if response.Header.Get("Content-Type") != "application/json" {
        t.Error("Response type not ok")
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Error(err)
    }
    /* Check for the length and the format of the response. 
    The response should respect the following struct:
        type BranchResponse struct {
            CurrentVersion  string          `json:"current_version"`
            Head            *object.Commit  `json:"head"`
            Branches        []string        `json:"branches"`
        }
    */
    if len(body) < 1 {
        t.Error("Response body length not ok. Should be at least 1.", len(body))
    }
    var response_struct map[string]interface{}
    err = json.Unmarshal(body, &response_struct)
    if err != nil {
        t.Error(err)
    }
    if response_struct["current_version"] == nil {
        t.Error("CurrentVersion not found in response struct.")
    }
    if response_struct["head"] == nil {
        t.Error("CurrentVersion not found in response struct.")
    }
    if response_struct["branches"] == nil {
        t.Error("Head not found in response struct.")
    }
}
// func GetCommits(w http.ResponseWriter, r *http.Request) {
// 	_, ok := utils.Repositories[r.URL.Query().Get("repo")]
//     i, err_i := strconv.Atoi(r.URL.Query().Get("i"))
//     j, err_j := strconv.Atoi(r.URL.Query().Get("j"))
// 
//     if !ok {
//         WriteError(&w, "Repository not found", 403)
//         return
//     }
// 
//     commits := navigation.GetAllCommits(r.URL.Query().Get("repo"))
//     if err_i != nil {
//         i = 0
//     }
//     if err_j != nil {
//         j = len(commits)
//     }
//     if i > len(commits) {
//         i = len(commits)
//     }
//     if j > len(commits) {
//         j = len(commits)
//     }
// 
//     response, err := json.Marshal(commits[i:j])
//     if err != nil {
//         WriteError(&w, "Error while getting release versions", 403)
//         return
//     }
//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusOK)
//     w.Write(response)
// }

func TestGetCommits(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:3001" +"/getCommits?repo=PWA-VUEJS&i=2&j=5", nil)
    req.Header.Add("Authorization", password)
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if response.StatusCode != http.StatusOK {
        t.Error("Status code not ok")
    }
    if response.Header.Get("Content-Type") != "application/json" {
        t.Error("Response type not ok")
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Error(err)
    }
    response_list := []interface{}{}
    err = json.Unmarshal(body, &response_list)
    if err != nil {
        t.Error(err)
    }
    if len(response_list) != 3 {
        t.Error("Response body length not ok. Should be 3", len(response_list))
    }
    messsage := response_list[0].(map[string]interface{})["Message"].(string)
    for i, item := range response_list {
        commit := item.(map[string]interface{})
        if _, ok := commit["Hash"]; !ok {
            t.Error("Hash not found in response list item.")
        }
        if _, ok := commit["Author"]; !ok {
            t.Error("Author not found in response list item.")
        }
        if _, ok := commit["Committer"]; !ok {
            t.Error("Committer not found in response list item.")
        }
        if _, ok := commit["Message"].(string); !ok {
            t.Error("Message not found in response list item.")
        }
        branches, ok := commit["branches"].([]interface{})
        if !ok || len(branches) < 1 {
            t.Error("branches not found in response list item. or length not ok.")
        }
        if i != 0 && messsage == commit["Message"].(string) {
            t.Error("Commits should be different.")
        }

    }
}
// func AddTimer(w http.ResponseWriter, r *http.Request) {
//     repo, ok := utils.Repositories[r.URL.Query().Get("repo")]
//     if !ok {
//         WriteError(&w, "Repository " + r.URL.Query().Get("repo") + " not found", http.StatusNotAcceptable)
//         return
//     }
// 
//     if _, err := utils.GetBranch(repo, r.URL.Query().Get("branch")); err != nil {
//         WriteError(&w, "Branch " + r.URL.Query().Get("branch") + " not found", http.StatusNotAcceptable)
//         return
//     }
// 
//     if secs, err := strconv.Atoi(r.URL.Query().Get("seconds")); err == nil && secs >= 60 {
//         w.Header().Set("Content-Type", "text")
//         w.WriteHeader(http.StatusOK)
//         w.Write([]byte("ok"))
// 
//         builddeploy.AddTimer(&builddeploy.AutobuildConfig{
//             Repo: r.URL.Query().Get("repo"),
//             Seconds: secs,
//             Branch: r.URL.Query().Get("branch"),
//         })
//         return
//     }
// 
//     WriteError(&w, "Format of \"seconds\" not correct or either has to be >= 60", http.StatusNotAcceptable)
// }
func TestAddTimer(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:3001" +"/addTimer?repo=PWA-VUEJS&branch=master&seconds=60", nil)
    req.Header.Add("Authorization", password)
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if response.StatusCode != http.StatusOK {
        t.Error("Status code not ok")
    }
    if response.Header.Get("Content-Type") != "text" {
        t.Error("Response type not ok")
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Error(err)
    }
    if string(body) != "ok" {
        t.Error("Response body not ok. Should be \"ok\"")
    }
}

// func GetTimers(w http.ResponseWriter, r *http.Request) {
// 
//     var configs = []*builddeploy.AutobuildConfig{}
//     for _, timer := range builddeploy.ActiveTimers {
//         configs = append(configs, timer.Config)
//     }
//     response, _ := json.Marshal(configs)
// 
//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusOK)
//     w.Write(response)
// }
func TestGetTimers(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:3001" +"/getTimers", nil)
    req.Header.Add("Authorization", password)
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if response.StatusCode != http.StatusOK {
        t.Error("Status code not ok")
    }
    if response.Header.Get("Content-Type") != "application/json" {
        t.Error("Response type not ok")
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Error(err)
    }
    response_list := []interface{}{}
    err = json.Unmarshal(body, &response_list)
    if err != nil {
        t.Error(err)
    }
    if len(response_list) != 1 {
        t.Error("Response body length not ok. Should be 1", len(body))
    }
    for _, item := range response_list {
        config := item.(map[string]interface{})
        if _, ok := config["Repo"]; !ok {
            t.Error("Repo not found in response list item.")
        }
        if _, ok := config["Branch"]; !ok {
            t.Error("Branch not found in response list item.")
        }
        if _, ok := config["Seconds"]; !ok {
            t.Error("Seconds not found in response list item.")
        }
    }
}
// func DeleteTimer(w http.ResponseWriter, r *http.Request) {
//     repo := r.URL.Query().Get("repo")
// 
//     if _, ok := builddeploy.ActiveTimers[repo]; ok {
//         builddeploy.DeleteTimer(repo)
//         w.Header().Set("Content-Type", "application/json")
//         w.WriteHeader(http.StatusOK)
//         w.Write([]byte("ok"))
//         return
//     }
// 
//     WriteError(&w, "Timer not found", http.StatusNotAcceptable)
// }
func TestDeleteTimer(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:3001" +"/deleteTimer?repo=PWA-VUEJS", nil)
    req.Header.Add("Authorization", password)
    response, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Error(err)
    }
    if response.StatusCode != http.StatusOK {
        t.Error("deleteTimer Status code not ok")
    }
    if response.Header.Get("Content-Type") != "application/json" {
        t.Error("Response type not ok")
    }
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Error(err)
    }
    if string(body) != "ok" {
        t.Error("Response body not ok. Should be \"ok\"")
    }
    req_get, _ := http.NewRequest("GET", "http://localhost:3001" +"/getTimers", nil)
    req_get.Header.Add("Authorization", password)
    response_get, err := http.DefaultClient.Do(req_get)
    if err != nil {
        t.Error(err)
    }
    if response_get.StatusCode != http.StatusOK {
        t.Error("getTimers Status code not ok")
        fmt.Println(response_get)
    }
    body_get, err := ioutil.ReadAll(response_get.Body)
    if err != nil {
        t.Error(err)
    }
    response_list := []interface{}{}
    err = json.Unmarshal(body_get, &response_list)
    if err != nil {
        t.Error(err)
    }
    if len(response_list) != 0 {
        t.Error("Response body length not ok. Should be 0", len(body))
    }
}
