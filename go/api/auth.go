package api

import (
	"context"
	"encoding/json"
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

	errorResp := logic.Register(ctx.Db, ctx.MailConfig, ctx.HostAddress, &userForCreate)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctx *Context) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tokenStr := strings.TrimPrefix(r.URL.Path, "/verify-email/")

	if tokenStr == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	errorResp := logic.VerifyEmail(ctx.Db, tokenStr)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	// TODO: check if errors are truly unreachable and remove 
	if ctx.NoUsers {
		users, errorResp := logic.GetUsers(ctx.Db, nil)
		// There should be no way to reach this error
		if len(users) == 0 || errorResp.Code != http.StatusOK {
			http.Error(w, errorResp.Message, errorResp.Code)
			return
		}

		ctx.NoUsers = false

		// This should not be reached either
		if len(users) > 1 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Email successfully verified."))
			return
		}

		// There is exactly one user, grant admin rights
		userID := users[0].ID
		errorResp = logic.GrantAdminRights(ctx.Db, &userID)
		if errorResp.Code != http.StatusOK {
			http.Error(w, errorResp.Message, errorResp.Code)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email successfully verified."))
}
