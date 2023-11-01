package routes

import (
	"net/http"

	"github.com/DiegoAraujoJS/webdev-git-server/pkg/utils"
)

func EnableCorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set the appropriate headers
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Origin, X-Requested-With")

        // Call the next handler in the chain
        next.ServeHTTP(w, r)
    })
}

func VerifyPasswordMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header["Authorization"] != nil && r.Header["Authorization"][0] == utils.ConfigValue.Credentials.Password {
            // Call the next handler in the chain
            next.ServeHTTP(w, r)
            return
        }
    })
}

func Ping(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/ping" {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("ok"))
            return
        }
        next.ServeHTTP(w, r)
    })
}
