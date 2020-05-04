package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/rossi1/blogapp/users"
)

func dbHandler() {
	db, err := gorm.Open("sqlite3", "blog.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&users.User{})
	defer db.Close()

}

func main() {
	
	dbHandler()

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})

	})
	r.HandleFunc("/register-user", users.RegisterUser).Methods("POST")
	r.HandleFunc("/login-user", users.LoginUser).Methods("POST")

	protectedRoutes := r.PathPrefix("/auth").Subrouter()
	protectedRoutes.HandleFunc("/profile/{id:[0-9]+}/", users.UserProfile).Methods("GET")
	protectedRoutes.HandleFunc("/profile-update/{id:[0-9]+}/", users.UpdateProfile).Methods("PUT")
	protectedRoutes.Use(users.JwtMiddleware)
	http.ListenAndServe(":8000", r)
}
