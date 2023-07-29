package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-git/go-git/v5"
)

// Config is the struct that holds the configuration for the git repositories.
//
// The configuration is read from a json file that is located in the root of the project.
// For deployment, the config.json is located on the same folder as the binary as this:
// .
// ├── config.json
// └── main

type ApplicationsType map[string]*struct{
    Directory string
    Repo *git.Repository
}

func (a ApplicationsType) GetAllRepos() []*git.Repository {
    var repos = map[*git.Repository]struct{}{}
    for _, app := range a {
        repos[app.Repo] = struct{}{}
    }
    var reposSlice = []*git.Repository{}
    for repo := range repos {
        reposSlice = append(reposSlice, repo)
    }
    return reposSlice
}

var (
    Applications = ApplicationsType{}
    ConfigValue struct {
        Port string
        Directories []struct {
            Name      string
            Directory string
        }
        Database struct {
            Server      string
            Port        string
            User        string
            Password    string
            Name        string
        }
        Env string
        Credentials struct {
            Password    string
            Ssh         string
        }
        CliBinaryForCheckout string
    }

)

func Connect() {
	content, err := ioutil.ReadFile("./config.json")

	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(content, &ConfigValue)

	if err != nil {
		log.Fatal(err.Error())
	}

    // The for loop below implements the following logic: An application with the same directory as another has also the same repository. (Numerical identity, same pointer)
	for _, dir := range ConfigValue.Directories {
        var repo *git.Repository
        for _, app := range Applications {
            if app.Directory == dir.Directory {
                repo = app.Repo
            }
        }

        if repo != nil {
            Applications[dir.Name] = &struct {
                Directory string
                Repo *git.Repository
            }{
                Directory: dir.Directory,
                Repo: repo,
            }
            continue
        }

		r, err := git.PlainOpen(dir.Directory)
		if err != nil {
			log.Fatal(err.Error())
		}
        Applications[dir.Name] = &struct {
            Directory string
            Repo *git.Repository
        }{
            Directory: dir.Directory,
            Repo: r,
        }
	}
}
