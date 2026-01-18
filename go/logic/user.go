package logic

import (
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = bcrypt.DefaultCost

type UserForCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserForUpdate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func CreateUser(db *database.Database, userForCreate *UserForCreate) (database.User, ErrorResponse) {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	userForInsert := &database.UserForInsert{
		Username:     userForCreate.Username,
		PasswordHash: hashPassword(userForCreate.Password),
		Email:        userForCreate.Email,
	}

	user, err := db.InsertUserDB(userForInsert)
	if err != nil {
		existingUser, err := db.SelectUserByEmailDB(userForInsert.Email)
		if existingUser != nil && !existingUser.EmailVerified {
			DeleteUser(db, existingUser.ID)
			user, err = db.InsertUserDB(userForInsert)
			if err == nil {
				return user, errorResp
			}
		}
		fmt.Println("Error inserting user:", err)
		return database.User{}, ErrorResponse{
			Message: "Failed to insert user",
			Code:    http.StatusInternalServerError,
		}
	}
	return user, errorResp
}

func GetUsers(db *database.Database) ([]*database.User, ErrorResponse) {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	users, err := db.SelectAllUsersDB()
	if err != nil {
		fmt.Println("Error retrieving users:", err)
		return nil, ErrorResponse{
			Message: "Failed to retrieve users",
			Code:    http.StatusInternalServerError,
		}
	}
	return users, errorResp
}

func GetUserByID(db *database.Database, id int64) (*database.User, ErrorResponse) {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	user, err := db.SelectUserByIDDB(id)
	if err != nil {
		fmt.Println("Error retrieving user:", err)
		return nil, ErrorResponse{
			Message: "User not found",
			Code:    http.StatusNotFound,
		}
	}
	return user, errorResp

}

func UpdateUser(db *database.Database, userForUpdate *UserForUpdate, id int64) ErrorResponse {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	userForUpdateDB := &database.UserForUpdate{
		ID:       id,
		Username: userForUpdate.Username,
		Email:    userForUpdate.Email,
		Password: userForUpdate.Password,
	}

	if err := db.UpdateUserDB(userForUpdateDB); err != nil {
		fmt.Println("Error updating user:", err)
		return ErrorResponse{
			Message: "Failed to update user",
			Code:    http.StatusInternalServerError,
		}
	}
	return errorResp
}

func DeleteUser(db *database.Database, id int64) ErrorResponse {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	if err := db.DeleteUserDB(id); err != nil {
		fmt.Println("Error deleting user:", err)
		return ErrorResponse{
			Message: "Failed to delete user",
			Code:    http.StatusInternalServerError,
		}
	}

	return errorResp
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return ""
	}
	return string(hash)
}


