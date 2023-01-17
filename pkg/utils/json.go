package utils

import (
	"encoding/json"
	"os"
)

func WriteIntoJson() error {
    new_json, err := json.MarshalIndent(ConfigValue, "", "    ")
    if err != nil {
        return err
    }
    file, err := os.Create("./config.json")
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(new_json)
    return err
}
