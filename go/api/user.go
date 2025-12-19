package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/Leander-s/money_manager/logic"
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

	var id int64
	if idStr == "self" {
		id = r.Context().Value("userID").(int64)
		fmt.Println("Resolved 'self' to user ID:", id)
	} else {
		idParsed, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		id = idParsed
	}

	switch r.Method {
	case http.MethodGet:
		ctx.GetUserByIDHandler(w, r, id)
	case http.MethodPut:
		ctx.UpdateUserHandler(w, r, id)
	case http.MethodDelete:
		ctx.DeleteUserHandler(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var ufc logic.UserForCreate
	if err := json.NewDecoder(r.Body).Decode(&ufc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, errorResp := logic.CreateUser(ctx.Db, &ufc)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	fmt.Println("Created user with ID:", user.ID)
}

func (ctx *Context) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, errorResp := logic.GetUsers(ctx.Db)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	fmt.Println("Retrieved all users")
}

func (ctx *Context) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, id int64) {
	user, errorResp := logic.GetUserByID(ctx.Db, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	fmt.Println("Retrieved user with ID:", id)
}

func (ctx *Context) UpdateUserHandler(w http.ResponseWriter, r *http.Request, id int64) {
	var userForUpdate logic.UserForUpdate
	if err := json.NewDecoder(r.Body).Decode(&userForUpdate); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	errorResp := logic.UpdateUser(ctx.Db, &userForUpdate, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Updated user with ID:", id)
}

func (ctx *Context) DeleteUserHandler(w http.ResponseWriter, r *http.Request, id int64) {
	errorResp := logic.DeleteUser(ctx.Db, w, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Deleted user with ID:", id)
}
