package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/Leander-s/money_manager/logic"
	"github.com/Leander-s/money_manager/db"
)

func (ctx *Context) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	switch r.Method {
	case http.MethodGet:
		ctx.HandleBalanceGet(w, userID)
	case http.MethodPost:
		ctx.HandleBalanceInsert(w, r, userID)
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
		ctx.HandleBalanceGetByCount(w, userID, count)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) BudgetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received balance", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Budget Path Accessed with method:", r.Method)
}

func (ctx *Context) HandleBalanceGet(w http.ResponseWriter, id int64) { 
	balances, errorResp := logic.GetAllBalances(ctx.Db, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balances)
	fmt.Println("Retrieved balance for user ID:", id)
}

func (ctx *Context) HandleBalanceGetByCount(w http.ResponseWriter, userID int64, count int64) {
	balances, errorResp := logic.GetBalanceByCount(ctx.Db, userID, count)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balances)
	fmt.Println("Retrieved ", count, " balances for user ID:", userID)
}

func (ctx *Context) HandleBalanceInsert(w http.ResponseWriter, r *http.Request, userID int64) {
	var entry *database.MoneyEntry
	if err := json.NewDecoder(r.Body).Decode(entry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	entry.UserID = userID
	entry, errorResp := logic.InsertBalance(ctx.Db, entry, userID)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
	fmt.Println("Inserted balance with ID:", entry.ID)
}
