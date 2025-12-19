package logic

import (
	"errors"
	"fmt"
	"time"
	"net/http"

	"github.com/Leander-s/money_manager/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code	int    `json:"code"`
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
		PrintAvailableTokens(db)
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

func PrintAvailableTokens(db *database.Database) {
	tokens, err := db.ListTokens()
	if err != nil {
		fmt.Println("Cannot print available tokens:", err.Error())
	} else {
		fmt.Println("Tokens are:")
		for _, t := range tokens {
			fmt.Println(t.Token.String(), " for user ", t.UserID)
		}
	}
}

func ValidateUser(db *database.Database, email string, password string) (int64, error) {
	user, err := db.SelectUserByEmailDB(email)
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

func Login(db *database.Database, loginReq *LoginRequest) (database.Token, ErrorResponse) {
	var token database.Token
	var errorResp ErrorResponse = ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	id, err := ValidateUser(db, loginReq.Email, loginReq.Password)
	if err != nil {
		errorResp = ErrorResponse{
			Message: "Login failed",
			Code:    http.StatusUnauthorized,
		}
		return token, errorResp
	}

	token = GenerateToken(id)
	err = db.DeleteTokensByUserID(id)
	if err != nil {
		fmt.Println("Failed to delete existing tokens for user:", err.Error())
		errorResp = ErrorResponse{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
		return token, errorResp
	}
	err = db.InsertToken(&token)
	if err != nil {
		fmt.Println("Failed to insert token:", err.Error())
		errorResp = ErrorResponse{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
		return token, errorResp
	}

	return token, errorResp
}

func Register(db *database.Database, registerReq *UserForCreate) ErrorResponse {
	var errorResp ErrorResponse = ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	_, err := CreateUser(db, registerReq)

	if err.Code != http.StatusOK {
		errorResp = ErrorResponse{
			Message: "Registration failed",
			Code:    http.StatusInternalServerError,
		}
		fmt.Println("Error creating user:", err)
		return errorResp
	}

	return errorResp
}
