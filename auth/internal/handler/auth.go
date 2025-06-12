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

var (
	db        *gorm.DB
	tokenAuth = jwtauth.New("HS256", []byte("secret-key"), nil)
)

func Init(database *gorm.DB) {
	db = database
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
	}

	if _, found := repository.GetUserByUsername(user.Username); found {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}
	user.Password, _ = utils.HashPassword(user.Password)
	if !repository.SaveUser(user) {
		http.Error(w, "Failed to save user ", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})

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

	_, err := utils.GenerateJWT(storedUser.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	json.NewEncoder(w).Encode(map[string]interface{}{"claims": claims})
}
