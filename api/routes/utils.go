package routes

import (
	"encoding/json"
	"net/http"
)

func WriteError(w *http.ResponseWriter, err string, status int) {
    w_ref := *w
    w_ref.Header().Set("Content-Type", "application/json")
    w_ref.WriteHeader(status)
    json_error, parse_error := json.Marshal(map[string]string{
        "error": err,
    })
    if parse_error != nil {
        w_ref.Write([]byte(`{"error": "Error parsing error"}`))
        return
    }
    w_ref.Write(json_error)
}
