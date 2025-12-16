package logic

import (
	"errors"
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"github.com/Leander-s/money_manager/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


func GenerateToken(userID int64) database.Token {
	expirationTime := time.Now().Add(24 * time.Hour)
	return database.Token{
		UserID: userID,
		Token:  uuid.New(),
		Expiry: expirationTime,
	}
}

func ValidateToken(db *database.Database, tokenStr string) (int64, error) {
	tokenID, err := uuid.Parse(tokenStr)
	if err != nil {
		return 0, errors.New("WrongFormat")
	}

	token, err := db.GetToken(tokenID)
	if err != nil {
		fmt.Println("Failed to get token:", tokenID.String())
		fmt.Println("Token error:", err.Error())
		tokens, err := db.ListTokens()
		if err != nil {
			fmt.Println("Cannot print available tokens:", err.Error())
		} else {
			fmt.Println("Tokens are:")
			for _, t := range tokens {
				fmt.Println(t.Token.String(), " for user ", t.UserID)
			}
		}
		return 0, errors.New("InvalidToken")
	}

	expirationTime := token.Expiry
	if time.Now().After(expirationTime) {
		err := db.DeleteToken(token.Token)
		if err != nil {
			fmt.Println("Failed to delete token:", err)
		}
		return 0, errors.New("TokenExpired")
	}

	return token.UserID, nil
}

func ValidateUser(db *database.Database, email string, password string) (int64, error) {
	user, err := db.GetUserByEmail(email)
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

func HandleLogin(db *database.Database, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := ValidateUser(db, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Login failed", http.StatusForbidden)
		return
	}

	token := GenerateToken(id)
	err = db.DeleteTokensByUserID(id)
	if err != nil {
		fmt.Println("Failed to delete existing tokens for user:", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = db.InsertToken(&token)
	if err != nil {
		fmt.Println("Failed to insert token:", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(token)
}

func HandleRegister(db *database.Database, w http.ResponseWriter, r *http.Request) {
	var user database.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// password is being converted to password hash here
	_, err = CreateUser(db, &user)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Println("Error inserting user:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
