package user

import (
	"encoding/json"
	"net/http"
)

type UserRegisterAPIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}

// user creation hanlar
// name, email, password
func CreateUserHandlar(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		response := UserRegisterAPIResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request payload",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer r.Body.Close()

	exists, err := IsUserExists(user.Email)
	if err != nil {
		response := UserRegisterAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error checking user existence",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if exists {
		response := UserRegisterAPIResponse{
			Status:  http.StatusConflict,
			Message: "User already exists",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	rspn, err := InsertUserData(user)
	if err != nil {
		response := UserRegisterAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error inserting user data",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// User inserted successfully, return a success response
	response := UserRegisterAPIResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Id:      rspn,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// user login handlar - session creation
// email, password -> user token

// user logout handlar - session- deletion
// user token

// user edit profile data
// user token
