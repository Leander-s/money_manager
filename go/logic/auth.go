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

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordExecutionRequest struct {
	Token       uuid.UUID `json:"token"`
	NewPassword string    `json:"new_password"`
}

func GenerateToken(userID *uuid.UUID) database.Token {
	if userID == nil {
		return database.Token{}
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	return database.Token{
		UserID: *userID,
		Token:  uuid.New(),
		Expiry: expirationTime,
	}
}

func ValidateToken(store database.TokenStore, tokenStr string) (uuid.UUID, error) {
	tokenID, err := uuid.Parse(tokenStr)
	if err != nil {
		return uuid.Nil, errors.New("WrongFormat")
	}

	token, err := store.GetToken(&tokenID)
	if err != nil {
		fmt.Println("Token error:", err.Error())
		return uuid.Nil, errors.New("InvalidToken")
	}

	expirationTime := token.Expiry
	if time.Now().After(expirationTime) {
		err := store.DeleteToken(&token.Token)
		if err != nil {
			fmt.Println("Failed to delete token:", err)
		}
		return uuid.Nil, errors.New("TokenExpired")
	}

	store.DeleteExpiredTokens()
	return token.UserID, nil
}

func PrintAvailableTokens(store database.TokenStore) {
	tokens, err := store.ListTokens()
	if err != nil {
		fmt.Println("Cannot print available tokens:", err.Error())
	} else {
		fmt.Println("Tokens are:")
		for _, t := range tokens {
			fmt.Println(t.Token.String(), " for user ", t.UserID)
		}
	}
}

func ValidateUser(store database.UserStore, email string, password string) (uuid.UUID, error) {
	user, err := store.SelectUserByEmailDB(email)
	if err != nil {
		fmt.Println("Error getting user with email:", email, ",", err)
		return uuid.Nil, err
	}

	if user.EmailVerified == false {
		fmt.Println("Unverified email login attempt for", email)
		return uuid.Nil, errors.New("EmailNotVerified")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Unsuccessful validation attempt for", user.Email)
		return uuid.Nil, errors.New("InvalidCredentials")
	}
	return user.ID, nil
}

func Login(store database.AuthStore, loginReq *LoginRequest) (database.Token, ErrorResponse) {
	var token database.Token
	var errorResp ErrorResponse = ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	id, err := ValidateUser(store, loginReq.Email, loginReq.Password)
	if err != nil {
		errorResp = ErrorResponse{
			Message: "Login failed: " + err.Error(),
			Code:    http.StatusUnauthorized,
		}
		return token, errorResp
	}

	token = GenerateToken(&id)
	err = store.DeleteTokensByUserID(&id)
	if err != nil {
		fmt.Println("Failed to delete existing tokens for user:", err.Error())
		errorResp = ErrorResponse{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
		return token, errorResp
	}
	err = store.InsertToken(&token)
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

func Register(store database.AuthStore, mailConfig EmailSender, hostAddress string, registerReq *UserForCreate) ErrorResponse {
	var errorResp ErrorResponse = ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	user, err := CreateUser(store, nil, registerReq)

	if err.Code != http.StatusOK {
		errorResp = ErrorResponse{
			Message: "Registration failed",
			Code:    http.StatusInternalServerError,
		}
		fmt.Println("Error creating user:", err)
		return errorResp
	}

	// Don't send email in test environment, immediately verify
	if hostAddress == "http://localhost:8080" {
		user.EmailVerified = true
		userForUpdate := database.UserForUpdate{
			ID:            user.ID,
			Username:      &user.Username,
			Password:      &user.Password,
			Email:         &user.Email,
			EmailVerified: &user.EmailVerified,
		}
		dbErr := store.UpdateUserDB(&userForUpdate)
		if dbErr != nil {
			fmt.Println("Failed to update user email verification status:", dbErr.Error())
			return ErrorResponse{
				Message: "Failed to verify email",
				Code:    http.StatusInternalServerError,
			}
		}

		// Give test users admin rights
		GrantAdminRights(store, &user.ID)

		// No email sent in test environment
		return errorResp
	}

	errorResp = SendEmailVerification(store, mailConfig, hostAddress, user)

	return errorResp
}

func SendPasswordResetEmail(store database.AuthStore, mailConfig EmailSender, frontendAddress string, request *ResetPasswordRequest) ErrorResponse {
	store.DeleteExpiredTokens()
	user, err := store.SelectUserByEmailDB(request.Email)
	if err != nil {
		fmt.Println("Failed to find user for password reset:", err.Error())
		return ErrorResponse{
			Message: "If the email is registered, a password reset link has been sent.",
			Code:    http.StatusOK,
		}
	}

	resetToken := GenerateToken(&user.ID)
	err = store.InsertToken(&resetToken)
	if err != nil {
		fmt.Println("Failed to insert password reset token:", err.Error())
		return ErrorResponse{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}

	err = mailConfig.SendEmail(user.Email, "Password Reset",
		fmt.Sprintf("Please reset your password using this link: %s", frontendAddress+"/reset-password/"+resetToken.Token.String()),
		"")
	if err != nil {
		fmt.Println("Failed to send password reset email:", err.Error())
		return ErrorResponse{
			Message: "Failed to send password reset email",
			Code:    http.StatusInternalServerError,
		}
	}
	return ErrorResponse{Message: "If the email is registered, a password reset link has been sent.", Code: http.StatusOK}
}

func ResetPassword(store database.AuthStore, token *uuid.UUID, newPassword string) ErrorResponse {
	hashedPassword := hashPassword(newPassword)
	userID, err := ValidateToken(store, token.String())
	if err != nil {
		fmt.Println("Password reset token validation failed:", err.Error())
		return ErrorResponse{
			Message: "Invalid or expired token",
			Code:    http.StatusUnauthorized,
		}
	}

	user, err := store.SelectUserByIDDB(&userID)
	if err != nil && user.EmailVerified == false {
		fmt.Println("Password reset attempt for unverified email:", user.Email)
		return ErrorResponse{
			Message: "Email not verified",
			Code:    http.StatusUnauthorized,
		}
	}

	userForUpdate := database.UserForUpdate{
		ID:       userID,
		Username: &user.Username,
		Password: &hashedPassword,
		Email:    &user.Email,
		EmailVerified: &user.EmailVerified,
	}
	errUpdate := store.UpdateUserDB(&userForUpdate)
	if errUpdate != nil {
		fmt.Println("Failed to update user password:", errUpdate.Error())
		return ErrorResponse{
			Message: "Failed to reset password",
			Code:    http.StatusInternalServerError,
		}
	}
	err = store.DeleteToken(token)
	if err != nil {
		fmt.Println("Failed to delete used password reset token:", err.Error())
	}
	return ErrorResponse{Message: "", Code: http.StatusOK}
}

func SendEmailVerification(store database.TokenStore, mailConfig EmailSender, hostAddress string, user *database.User) ErrorResponse {
	store.DeleteExpiredTokens()
	verificationToken := GenerateToken(&user.ID)
	err := store.InsertToken(&verificationToken)
	if err != nil {
		fmt.Println("Failed to insert email verification token:", err.Error())
		return ErrorResponse{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}

	// Placeholder for email sending logic
	err = mailConfig.SendEmail(user.Email, "Email Verification",
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

func VerifyEmail(store database.AuthStore, tokenStr string) ErrorResponse {
	tokenID, err := uuid.Parse(tokenStr)
	if err != nil {
		return ErrorResponse{
			Message: "Invalid token format",
			Code:    http.StatusBadRequest,
		}
	}

	token, err := store.GetToken(&tokenID)
	if err != nil {
		return ErrorResponse{
			Message: "Invalid token",
			Code:    http.StatusUnauthorized,
		}
	}

	expirationTime := token.Expiry
	if time.Now().After(expirationTime) {
		store.DeleteToken(&token.Token)
		return ErrorResponse{
			Message: "Token expired",
			Code:    http.StatusUnauthorized,
		}
	}

	user, err := store.SelectUserByIDDB(&token.UserID)
	if err != nil {
		return ErrorResponse{
			Message: "User not found",
			Code:    http.StatusNotFound,
		}
	}

	user.EmailVerified = true
	userForUpdate := database.UserForUpdate{
		ID:            user.ID,
		Username:      &user.Username,
		Password:      &user.Password,
		Email:         &user.Email,
		EmailVerified: &user.EmailVerified,
	}
	err = store.UpdateUserDB(&userForUpdate)
	if err != nil {
		fmt.Println("Failed to update user email verification status:", err.Error())
		return ErrorResponse{
			Message: "Failed to verify email",
			Code:    http.StatusInternalServerError,
		}
	}

	store.DeleteToken(&token.Token)
	return ErrorResponse{Message: "", Code: http.StatusOK}
}

func GrantModeratorRights(store database.RoleStore, userID *uuid.UUID) ErrorResponse {
	var errorResp ErrorResponse = ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	err := store.AssignRoleToUserDB(userID, "moderator")
	if err != nil {
		fmt.Println("Error granting moderator rights:", err)
		return ErrorResponse{
			Message: "Failed to grant moderator rights",
			Code:    http.StatusInternalServerError,
		}
	}

	return errorResp
}

func GrantAdminRights(store database.RoleStore, userID *uuid.UUID) ErrorResponse {
	var errorResp ErrorResponse = ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	errorResp = GrantModeratorRights(store, userID)
	if errorResp.Code != http.StatusOK {
		return errorResp
	}

	err := store.AssignRoleToUserDB(userID, "admin")
	if err != nil {
		fmt.Println("Error granting admin rights:", err)
		return ErrorResponse{
			Message: "Failed to grant admin rights",
			Code:    http.StatusInternalServerError,
		}
	}

	return errorResp
}

func CheckRole(db database.RoleStore, userID *uuid.UUID, role string) (bool, error) {
	if userID == nil {
		return true, nil
	}

	hasRole, err := db.CheckUserRoleDB(userID, role)
	if err != nil {
		fmt.Println("Error checking user role:", err)
		return false, err
	}
	return hasRole, nil
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return ""
	}
	return string(hash)
}
