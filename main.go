package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	// u := user.NewUser(1, "john_doe", "john@example.com")

	// fmt.Println("User ID:", u.ID)
	// fmt.Println("Username:", u.Username)
	// fmt.Println("Email:", u.Email)

	fmt.Println("hi from erp :)")

	router := CreateRouter()

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Add your frontend URL here
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	http.Handle("/", corsHandler(router))

	// Start the HTTP server with the created router
	// http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
