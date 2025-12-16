package api

import (
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
	logic.HandleLogin(ctx.Db, w, r)
}

func (ctx *Context) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	logic.HandleRegister(ctx.Db, w, r)
}
