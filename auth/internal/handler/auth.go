package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bhati00/go-learning/auth/internal/model"
	"github.com/bhati00/go-learning/auth/internal/repository"
	"github.com/bhati00/go-learning/auth/internal/utils"
	"github.com/go-chi/jwtauth/v5"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database *gorm.DB) {
	db = database
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if _, found := repository.GetUserByUsername(user.Username); found {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}
	user.Password, _ = utils.HashPassword(user.Password)
	if !repository.SaveUser(user) {
		http.Error(w, "Failed to save user ", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"username": user.Username,
	})

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	storedUser, found := repository.GetUserByUsername(user.Username)
	if !found || !utils.CheckPasswordHash(user.Password, storedUser.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(storedUser.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	json.NewEncoder(w).Encode(map[string]interface{}{"claims": claims})
}
