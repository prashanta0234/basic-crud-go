package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := DB.Query("SELECT id, name, email FROM users")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			rows.Scan(&u.ID, &u.Name, &u.Email)
			users = append(users, u)
		}
		jsonResponse(w, users, 200)

	case http.MethodPost:
		var u User
		json.NewDecoder(r.Body).Decode(&u)
		err := DB.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResponse(w, u, 201)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", 400)
		return
	}

	switch r.Method {
	case http.MethodGet:
		var u User
		err := DB.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			http.Error(w, "User not found", 404)
			return
		}
		jsonResponse(w, u, 200)

	case http.MethodPut:
		var u User
		json.NewDecoder(r.Body).Decode(&u)
		_, err := DB.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", u.Name, u.Email, id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		u.ID = id
		jsonResponse(w, u, 200)

	case http.MethodDelete:
		_, err := DB.Exec("DELETE FROM users WHERE id=$1", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResponse(w, map[string]string{"message": "Deleted"}, 200)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
