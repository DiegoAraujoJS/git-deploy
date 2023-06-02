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
    // Test that the commits have the correct format
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
    if len(body) == 3 {
        t.Error("Response body length not ok. Should be 3", len(body))
    }
    for _, item := range response_list {
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
        fmt.Println(branches)
        if !ok || len(branches) < 1 {
            t.Error("branches not found in response list item. or length not ok.")
        }
    }
}
