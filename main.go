package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	LoadEnv()
	InitDB()
	defer DB.Close()

	http.HandleFunc("/users", UsersHandler)
	http.HandleFunc("/users/", UserHandler) // for GET/PUT/DELETE by id

	port := os.Getenv("PORT")
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
