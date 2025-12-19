package api

import (
	"encoding/json"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Leander-s/money_manager/logic"
)

func (ctx *Context) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			fmt.Println("Token did not have correct prefix")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		auth = strings.TrimPrefix(auth, "Bearer ")
		userID, err := logic.ValidateToken(ctx.Db, auth)
		if err != nil {
			fmt.Println("Could not validate token:", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// If authentication succeeds, proceed to the next handler
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ctx *Context) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var loginReq logic.LoginRequest
	json.NewDecoder(r.Body).Decode(&loginReq)

	token, err := logic.Login(ctx.Db, &loginReq)
	if err.Code != http.StatusOK {
		http.Error(w, err.Message, err.Code)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (ctx *Context) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var userForCreate logic.UserForCreate
	err := json.NewDecoder(r.Body).Decode(&userForCreate)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	errorResp := logic.Register(ctx.Db, &userForCreate)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
