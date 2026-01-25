package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leander-s/money_manager/logic"
	"github.com/google/uuid"
)

func (ctx *Context) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ctx.CreateHandler(w, r)
	case http.MethodGet:
		ctx.GetAllUsersHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) UserHandlerByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/user/")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(uuid.UUID)

	var id uuid.UUID
	if idStr == "self" {
		fmt.Println("Resolved 'self' to user ID:", userID)
		id = userID
	} else {
		admin, err := logic.CheckRole(ctx.Db, &userID, "admin")
		moderator, err := logic.CheckRole(ctx.Db, &userID, "moderator")
		if err != nil {
			http.Error(w, "Error checking user roles", http.StatusInternalServerError)
			return
		}
		if !admin && !moderator {
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}
		idParsed, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		id = idParsed
	}

	switch r.Method {
	case http.MethodGet:
		ctx.GetUserByIDHandler(w, r, &id)
	case http.MethodPut:
		ctx.UpdateUserHandler(w, r, &id)
	case http.MethodDelete:
		ctx.DeleteUserHandler(w, r, &id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var actorID uuid.UUID = r.Context().Value("userID").(uuid.UUID)

	var ufc logic.UserForCreate
	if err := json.NewDecoder(r.Body).Decode(&ufc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, errorResp := logic.CreateUser(ctx.Db, &actorID, &ufc)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	fmt.Println("Created user with ID:", user.ID)
}

func (ctx *Context) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	var actorID uuid.UUID = r.Context().Value("userID").(uuid.UUID)

	users, errorResp := logic.GetUsers(ctx.Db, &actorID)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	fmt.Println("Retrieved all users")
}

func (ctx *Context) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, id *uuid.UUID) {
	var actorID uuid.UUID = r.Context().Value("userID").(uuid.UUID)

	user, errorResp := logic.GetUserByID(ctx.Db, &actorID, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	fmt.Println("Retrieved user with ID:", id)
}

func (ctx *Context) UpdateUserHandler(w http.ResponseWriter, r *http.Request, id *uuid.UUID) {
	var userForUpdate logic.UserForUpdate
	if err := json.NewDecoder(r.Body).Decode(&userForUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var actorID uuid.UUID = r.Context().Value("userID").(uuid.UUID)

	errorResp := logic.UpdateUser(ctx.Db, &userForUpdate, &actorID, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Updated user with ID:", *id)
}

func (ctx *Context) DeleteUserHandler(w http.ResponseWriter, r *http.Request, id *uuid.UUID) {
	var actorID uuid.UUID = r.Context().Value("userID").(uuid.UUID)

	errorResp := logic.DeleteUser(ctx.Db, &actorID, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Deleted user with ID:", id)
}
