package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func WriteResponseOk(w *http.ResponseWriter, body interface{}) {
    w_ref := *w
    w_ref.Header().Set("Content-Type", "application/json")
    w_ref.WriteHeader(http.StatusOK)
    json_body, parse_error := json.Marshal(body)
    if parse_error != nil {
        WriteError(w, "Error parsing response", http.StatusInternalServerError)
        return
    }
    w_ref.Write(json_body)
}

func NormalizeSliceIndexes (response_length int, r *http.Request) (int, int) {
	i, err_i := strconv.Atoi(r.URL.Query().Get("i"))
	j, err_j := strconv.Atoi(r.URL.Query().Get("j"))

	if err_i != nil {
		i = 0
	}
	if err_j != nil {
		j = response_length
	}
	if i > response_length {
		i = response_length
	}
	if j > response_length {
		j = response_length
	}
	if i < 0 {
		i = 0
	}
	if j < 0 {
		j = 0
	}

    return i, j
}
