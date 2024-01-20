package main

import (
	"net/http"

	"erp.com/erp/user"
	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// corsHandler := handlers.CORS(
	// 	handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Add your frontend URL here
	// 	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	// 	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	// )

	// http.Handle("/", corsHandler(router))

	// Define your routes and handlers here
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/about", AboutHandler).Methods("GET")

	// // Example of a route with a URL parameter
	// router.HandleFunc("/user/{id}", GetUserHandler).Methods("GET")

	// // user endpoints
	router.HandleFunc("/api/user/register", user.CreateUserHandlar).Methods("POST")
	// router.HandleFunc("/api/user/login", user.LoginHandler).Methods("POST")
	// router.HandleFunc("/api/user/register", user.RegisterUserHandler).Methods("POST")
	// router.HandleFunc("/api/user/logout", user.LogoutUserHandler).Methods("POST")

	// // product endpoints
	// router.HandleFunc("/api/product/create", product.CreateProductHandler).Methods("POST")

	user.User_init()

	return router
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the home route
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the home page!"))
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the about route
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is the about page!"))
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get and handle the user ID from the URL parameter
	vars := mux.Vars(r)
	userID := vars["id"]

	// Handle the user route with the given ID
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Getting user with ID: " + userID))
}
