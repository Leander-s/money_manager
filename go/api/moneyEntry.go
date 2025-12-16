package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Leander-s/money_manager/logic"
)

func (ctx *Context) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	switch r.Method {
	case http.MethodGet:
		logic.HandleGetBalance(ctx.Db, w, userID)
	case http.MethodPost:
		logic.HandleInsertBalance(ctx.Db, w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) BalanceHandlerByCount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	countStr := strings.TrimPrefix(r.URL.Path, "/balance/")
	if countStr == "" {
		http.Error(w, "Count is required", http.StatusBadRequest)
		return
	}

	var count int64
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid count", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		logic.HandleGetBalanceByCount(ctx.Db, w, userID, count)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) BudgetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received balance", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Budget Path Accessed with method:", r.Method)
}

