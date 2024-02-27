package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/toffernator/elf/store"
)

var s *store.ArrayStore = &store.ArrayStore{}

func main() {
	s.Seed()

	http.HandleFunc("GET /api/user/{id}", HttpGetUser)
	http.HandleFunc("POST /api/user/", HttpPostUser)

	println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HttpGetUser(w http.ResponseWriter, r *http.Request) {
	slog.Info("GET /api/user")

	id, err := strconv.Atoi(r.PathValue("id"))
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

func HttpPostUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data []byte
	r.Body.Read(data)

	var newUser NewUserRequest
	json.Unmarshal(data, &newUser)

	user := s.CreateUser(newUser.FirstName, newUser.LastName)
	body, _ := json.Marshal(user)
	w.Write(body)
}
