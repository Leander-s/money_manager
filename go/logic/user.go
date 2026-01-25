package logic

import (
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = bcrypt.DefaultCost

type UserForCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserForUpdate struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

func CreateUser(store database.UserRoleStore, actorID *uuid.UUID, userForCreate *UserForCreate) (*database.User, ErrorResponse) {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	isModerator, err := CheckRole(store, actorID, "moderator")
	if err != nil {
		fmt.Println("Error checking actor role:", err)
		return nil, ErrorResponse{
			Message: "Failed to check actor role",
			Code:    http.StatusInternalServerError,
		}
	}
	if !isModerator && actorID != nil {
		return nil, ErrorResponse{
			Message: "Forbidden: insufficient permissions",
			Code:    http.StatusForbidden,
		}
	}

	userForInsert := &database.UserForInsert{
		Username:     userForCreate.Username,
		PasswordHash: hashPassword(userForCreate.Password),
		Email:        userForCreate.Email,
	}

	user, err := store.InsertUserDB(userForInsert)
	if err != nil {
		existingUser, err := store.SelectUserByEmailDB(userForInsert.Email)
		if existingUser != nil && !existingUser.EmailVerified {
			DeleteUser(store, &existingUser.ID, &existingUser.ID)
			user, err = store.InsertUserDB(userForInsert)
			if err == nil {
				return &user, errorResp
			}
		}
		fmt.Println("Error inserting user:", err)
		return nil, ErrorResponse{
			Message: "Failed to insert user",
			Code:    http.StatusInternalServerError,
		}
	}

	store.AssignRoleToUserDB(&user.ID, "user")

	return &user, errorResp
}

func GetUsers(store database.UserRoleStore, actorID *uuid.UUID) ([]*database.User, ErrorResponse) {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	isModerator, err := CheckRole(store, actorID, "moderator")
	if err != nil {
		fmt.Println("Error checking actor role:", err)
		return nil, ErrorResponse{
			Message: "Failed to check actor role",
			Code:    http.StatusInternalServerError,
		}
	}
	if !isModerator {
		return nil, ErrorResponse{
			Message: "Forbidden: insufficient permissions",
			Code:    http.StatusForbidden,
		}
	}

	users, err := store.SelectAllUsersDB()
	if err != nil {
		fmt.Println("Error retrieving users:", err)
		return nil, ErrorResponse{
			Message: "Failed to retrieve users",
			Code:    http.StatusInternalServerError,
		}
	}
	return users, errorResp
}

func GetUserByID(store database.UserRoleStore, actorID *uuid.UUID, id *uuid.UUID) (*database.User, ErrorResponse) {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	isModerator, err := CheckRole(store, actorID, "moderator")
	if err != nil {
		fmt.Println("Error checking actor role:", err)
		return nil, ErrorResponse{
			Message: "Failed to check actor role",
			Code:    http.StatusInternalServerError,
		}
	}
	if !isModerator && *id != *actorID {
		return nil, ErrorResponse{
			Message: "Forbidden: insufficient permissions",
			Code:    http.StatusForbidden,
		}
	}

	user, err := store.SelectUserByIDDB(id)
	if err != nil {
		fmt.Println("Error retrieving user:", err)
		return nil, ErrorResponse{
			Message: "User not found",
			Code:    http.StatusNotFound,
		}
	}
	return user, errorResp

}

func UpdateUser(store database.UserRoleStore, userForUpdate *UserForUpdate, actorID *uuid.UUID, id *uuid.UUID) ErrorResponse {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	if userForUpdate.Password != nil || userForUpdate.Email != nil {
		errorResp = ErrorResponse{
			Message: "Forbidden: cannot update password or email",
			Code:    http.StatusForbidden,
		}
		return errorResp
	}

	isModerator, err := CheckRole(store, actorID, "moderator")
	isAdmin, err := CheckRole(store, actorID, "admin")

	if !isModerator && *id != *actorID {
		return ErrorResponse{
			Message: "Forbidden: insufficient permissions",
			Code:    http.StatusForbidden,
		}
	}

	userToUpdateIsMod, err := CheckRole(store, id, "moderator")
	userToUpdateIsAdmin, err := CheckRole(store, id, "admin")
	if err != nil {
		fmt.Println("Error checking user:", *id, "'s role:", err)
		return ErrorResponse{
			Message: "Failed to check role",
			Code:    http.StatusInternalServerError,
		}
	}

	if userToUpdateIsMod && !isAdmin {
		return ErrorResponse{
			Message: "Forbidden: only admins can update moderator users",
			Code:    http.StatusForbidden,
		}
	}

	if userToUpdateIsAdmin && *id != *actorID {
		return ErrorResponse{
			Message: "Forbidden: cannot update an admin user",
			Code:    http.StatusForbidden,
		}
	}

	userForUpdateDB := &database.UserForUpdate{
		ID:       *id,
		Username: userForUpdate.Username,
		Email:    userForUpdate.Email,
		Password: userForUpdate.Password,
	}

	if err := store.UpdateUserDB(userForUpdateDB); err != nil {
		fmt.Println("Error updating user:", err)
		return ErrorResponse{
			Message: "Failed to update user",
			Code:    http.StatusInternalServerError,
		}
	}
	return errorResp
}

func DeleteUser(store database.UserRoleStore, actorID *uuid.UUID, id *uuid.UUID) ErrorResponse {
	errorResp := ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}

	isModerator, err := CheckRole(store, actorID, "moderator")
	isAdmin, err := CheckRole(store, actorID, "admin")

	if !isModerator && *id != *actorID {
		return ErrorResponse{
			Message: "Forbidden: insufficient permissions",
			Code:    http.StatusForbidden,
		}
	}

	userToUpdateIsMod, err := CheckRole(store, id, "moderator")
	userToUpdateIsAdmin, err := CheckRole(store, id, "admin")
	if err != nil {
		fmt.Println("Error checking user:", *id, "'s role:", err)
		return ErrorResponse{
			Message: "Failed to check role",
			Code:    http.StatusInternalServerError,
		}
	}

	if userToUpdateIsMod && !isAdmin {
		return ErrorResponse{
			Message: "Forbidden: only admins can delete moderator users",
			Code:    http.StatusForbidden,
		}
	}

	if userToUpdateIsAdmin && *id != *actorID {
		return ErrorResponse{
			Message: "Forbidden: cannot delete an admin user",
			Code:    http.StatusForbidden,
		}
	}
	if err := store.DeleteUserDB(id); err != nil {
		fmt.Println("Error deleting user:", err)
		return ErrorResponse{
			Message: "Failed to delete user",
			Code:    http.StatusInternalServerError,
		}
	}

	return errorResp
}
