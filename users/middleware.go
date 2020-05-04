package users

import (
	"encoding/json"
	"net/http"
	
)
func JwtMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ValidateJwtToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err.Error())
		} else {
			next.ServeHTTP(w, r)
		}
        
    })
}