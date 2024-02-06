package main

import (
	"encoding/json"
	"html"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/toffernator/elf/store"
)

var s *store.ArrayStore = &store.ArrayStore{}

func main() {
	s.Seed()

	http.HandleFunc("/api/user/", user)

	println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func user(w http.ResponseWriter, r *http.Request) {
	slog.Info(r.Method)
	if r.Method == http.MethodGet {
		slog.Info("GET /api/user")
		GETUser(w, r)
	} else if r.Method == http.MethodPost {
		slog.Info("POST /api/user")
		POSTUser(w, r)
	}
}

func GETUser(w http.ResponseWriter, r *http.Request) {
	path := html.EscapeString(r.URL.Path)
	parts := strings.Split(path, "/")[1:]

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.GetUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	body, _ := json.Marshal(user)
	w.Write(body)
}

type NewUserRequest struct {
	FirstName string
	LastName  string
}

func POSTUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data []byte
	r.Body.Read(data)

	var newUser NewUserRequest
	json.Unmarshal(data, &newUser)

	user := s.CreateUser(newUser.FirstName, newUser.LastName)
	body, _ := json.Marshal(user)
	w.Write(body)
}
