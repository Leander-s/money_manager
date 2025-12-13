package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Leander-s/money_manager/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateToken(userID int64) model.Token {
	expirationTime := time.Now().Add(24 * time.Hour)
	return model.Token{
		UserID: userID,
		Token:  uuid.New(),
		Expiry: expirationTime,
	}
}

func (app *App) validateToken(tokenStr string) (int64, error) {
	tokenID, err := uuid.Parse(tokenStr)
	if err != nil {
		return 0, errors.New("WrongFormat")
	}

	token, err := app.db.GetToken(tokenID)
	if err != nil {
		fmt.Println("Token error:", err.Error())
		tokens, err := app.db.ListTokens()
		if err != nil {
			fmt.Println("Cannot print available tokens:", err.Error())
		} else {
			fmt.Println("Tokens are:", tokens)
		}
		return 0, errors.New("InvalidToken")
	}

	expirationTime := token.Expiry
	if time.Now().After(expirationTime) {
		err := app.db.DeleteToken(token.Token)
		if err != nil {
			fmt.Println("Failed to delete token:", err)
		}
		return 0, errors.New("TokenExpired")
	}

	return token.UserID, nil
}

func (app *App) ValidateUser(email string, password string) (int64, error) {
	user, err := app.db.GetUserByEmail(email)
	if err != nil {
		fmt.Println("Error getting user with email:", email, ",", err)
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Unsuccessful validation attempt for", user.Email)
		return 0, errors.New("InvalidCredentials")
	}
	return user.ID, nil
}

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := app.ValidateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Login failed", http.StatusForbidden)
		return
	}

	token := generateToken(id)
	err = app.db.InsertToken(&token)
	if err != nil {
		fmt.Println("Failed to insert token:", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(token)
}

func (app *App) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// password is being converted to password hash here
	_, err = app.CreateUser(&user)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Println("Error inserting user:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
