package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Leander-s/money_manager/logic"
)

func (ctx *Context) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		logic.HandleCreateUser(ctx.Db, w, r)
	case http.MethodGet:
		logic.HandleGetUsers(ctx.Db, w)
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
		logic.HandleGetUserByID(ctx.Db, w, id)
	case http.MethodPut:
		logic.HandleUpdateUser(ctx.Db, w, r, id)
	case http.MethodDelete:
		logic.HandleDeleteUser(ctx.Db, w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
