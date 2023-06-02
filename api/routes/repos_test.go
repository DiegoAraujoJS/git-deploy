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
    fmt.Println(response.Header.Get("Content-Type"))
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
    fmt.Println(response.Header.Get("Content-Type"))
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
