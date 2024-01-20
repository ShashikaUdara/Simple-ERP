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

type UserLoginAPIResponse struct {
	Status       int         `json:"status"`
	Message      string      `json:"message"`
	UserID       interface{} `json:"user_id"`
	SessionToken string      `json:"session_token"`
}

type UserLogoutAPIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type UserProfileUpdateAPIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Success bool   `json:"success"`
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

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body into User struct
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isValid, err := VerifyUser(user)
	if err != nil {
		response := UserLoginAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalied user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if isValid {
		session_token, err := CreateUserSession(user, r)
		if err != nil {
			response := UserLoginAPIResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error creating user session",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := UserLoginAPIResponse{
			Status:       http.StatusCreated,
			Message:      "Session created successfully",
			UserID:       user.Email,
			SessionToken: session_token,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	} else {
		response := UserLoginAPIResponse{
			Status:  http.StatusNoContent,
			Message: "User credentials are wrong",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}

// user logout handlar - session- deletion
// user token
func LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := GetUserBearerToken(r)
	if err != nil {
		response := UserLogoutAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalied user token",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	isValid, err := VerifyUserSession(token)
	if err != nil {
		response := UserLogoutAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error verifiying user session",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if isValid {
		isLogout, err := UpdateUserSession(token)
		if err != nil {
			response := UserLogoutAPIResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error logging out from the user session",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := UserLogoutAPIResponse{
			Status:  http.StatusOK,
			Message: "Session deactivated successfully",
			Success: isLogout,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		response := UserLoginAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalied session",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
}

// user edit profile data
// user token + profile data
func UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	token, err := GetUserBearerToken(r)
	if err != nil {
		response := UserLogoutAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalied user token",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	isValid, err := VerifyUserSession(token)
	if err != nil {
		response := UserLogoutAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error verifiying user session",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if isValid {
		var user User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		isUpdated, err := UpdateUserProfileData(user, token)
		if err != nil {
			response := UserLogoutAPIResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error updating the user profile",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := UserProfileUpdateAPIResponse{
			Status:  http.StatusOK,
			Message: "User profile updated successfully",
			Success: isUpdated,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		response := UserLoginAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalied session",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
}

// user edit profile data
// user token + new password
func UserPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	token, err := GetUserBearerToken(r)
	if err != nil {
		response := UserLogoutAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalied user token",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	isValid, err := VerifyUserSession(token)
	if err != nil {
		response := UserLogoutAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error verifiying user session",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if isValid {
		var user User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		isUpdated, err := UserPasswordReset(user, token)
		if err != nil {
			response := UserLogoutAPIResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error resetting the password",
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := UserProfileUpdateAPIResponse{
			Status:  http.StatusOK,
			Message: "Password reset successful",
			Success: isUpdated,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		response := UserLoginAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalied session",
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
}
