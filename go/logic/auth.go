package logic

import (
	"errors"
	"fmt"
	"net/http"
	"time"

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
	Code    int    `json:"code"`
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

	db.DeleteExpiredTokens()
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

	if user.EmailVerified == false {
		fmt.Println("Unverified email login attempt for", email)
		return 0, errors.New("EmailNotVerified")
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
			Message: "Login failed: " + err.Error(),
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

func Register(db *database.Database, mailConfig *BrevoConfig, hostAddress string, registerReq *UserForCreate) ErrorResponse {
	var errorResp ErrorResponse = ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	user, err := CreateUser(db, registerReq)

	if err.Code != http.StatusOK {
		errorResp = ErrorResponse{
			Message: "Registration failed",
			Code:    http.StatusInternalServerError,
		}
		fmt.Println("Error creating user:", err)
		return errorResp
	}

	errorResp = SendEmailVerification(db, mailConfig, hostAddress, &user)

	return errorResp
}

func SendEmailVerification(db *database.Database, mailConfig *BrevoConfig, hostAddress string, user *database.User) ErrorResponse {
	db.DeleteExpiredTokens()
	verificationToken := GenerateToken(user.ID)
	err := db.InsertToken(&verificationToken)
	if err != nil {
		fmt.Println("Failed to insert email verification token:", err.Error())
		return ErrorResponse{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}

	// Placeholder for email sending logic
	err = SendEmailBrevo(mailConfig, user.Email, "Email Verification",
		fmt.Sprintf("Please verify your email using this link: %s", hostAddress+"/verify-email/"+verificationToken.Token.String()),
		"")
	if err != nil {
		fmt.Println("Failed to send verification email:", err.Error())
		return ErrorResponse{
			Message: "Failed to send verification email",
			Code:    http.StatusInternalServerError,
		}
	}
	return ErrorResponse{Message: "", Code: http.StatusOK}
}

func VerifyEmail(db *database.Database, tokenStr string) ErrorResponse {
	tokenID, err := uuid.Parse(tokenStr)
	if err != nil {
		return ErrorResponse{
			Message: "Invalid token format",
			Code:    http.StatusBadRequest,
		}
	}

	token, err := db.GetToken(tokenID)
	if err != nil {
		return ErrorResponse{
			Message: "Invalid token",
			Code:    http.StatusUnauthorized,
		}
	}

	expirationTime := token.Expiry
	if time.Now().After(expirationTime) {
		db.DeleteToken(token.Token)
		return ErrorResponse{
			Message: "Token expired",
			Code:    http.StatusUnauthorized,
		}
	}

	user, err := db.SelectUserByIDDB(token.UserID)
	if err != nil {
		return ErrorResponse{
			Message: "User not found",
			Code:    http.StatusNotFound,
		}
	}

	user.EmailVerified = true
	userForUpdate := database.UserForUpdate{
		ID:            user.ID,
		Username:      user.Username,
		Password:      user.Password,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}
	err = db.UpdateUserDB(&userForUpdate)
	if err != nil {
		fmt.Println("Failed to update user email verification status:", err.Error())
		return ErrorResponse{
			Message: "Failed to verify email",
			Code:    http.StatusInternalServerError,
		}
	}

	db.DeleteToken(token.Token)
	return ErrorResponse{Message: "", Code: http.StatusOK}
}
