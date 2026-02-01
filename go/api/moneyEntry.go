package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Leander-s/money_manager/db"
	"github.com/Leander-s/money_manager/logic"
	"github.com/google/uuid"
)

func (ctx *Context) BalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	switch r.Method {
	case http.MethodGet:
		ctx.HandleBalanceGet(w, &userID)
	case http.MethodPost:
		ctx.HandleBalanceInsert(w, r, &userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) BalanceHandlerByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/balance/id/")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	balanceID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ctx.HandleBalanceGetByID(w, &balanceID)
	case http.MethodDelete:
		ctx.HandleBalanceDelete(w, &balanceID)
	case http.MethodPut:
		ctx.HandleBalanceUpdate(w, r, &balanceID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) BalanceHandlerByCount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(uuid.UUID)
	countStr := strings.TrimPrefix(r.URL.Path, "/balance/count/")
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
		ctx.HandleBalanceGetByCount(w, &userID, count)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (ctx *Context) BudgetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received balance", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Budget Path Accessed with method:", r.Method)
}

func (ctx *Context) HandleBalanceUpdate(w http.ResponseWriter, r *http.Request, balanceID *uuid.UUID) {
	// Get actor ID from context
	actorID := r.Context().Value("userID").(uuid.UUID)

	// Decode request body
	var entryForUpdate logic.EntryForUpdate
	err := json.NewDecoder(r.Body).Decode(&entryForUpdate)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	entryForUpdate.ID = *balanceID

	// Call logic to update balance
	entries, errorResp := logic.UpdateBalance(ctx.Db, &actorID, &entryForUpdate)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	// Respond with updated entries
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
	fmt.Println("Updated entry with ID:", balanceID)
}

func (ctx *Context) HandleBalanceDelete(w http.ResponseWriter, balanceID *uuid.UUID) {
	entries, errorResp := logic.DeleteBalance(ctx.Db, balanceID)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
	fmt.Println("Deleted entry with ID:", balanceID)
}

func (ctx *Context) HandleBalanceGetByID(w http.ResponseWriter, balanceID *uuid.UUID) {
	balance, errorResp := logic.GetBalanceByID(ctx.Db, balanceID)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
	fmt.Println("Retrieved balance with ID:", balanceID)
}

func (ctx *Context) HandleBalanceGet(w http.ResponseWriter, id *uuid.UUID) { 
	balances, errorResp := logic.GetAllBalances(ctx.Db, id)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balances)
	fmt.Println("Retrieved balances for user ID:", id)
}

func (ctx *Context) HandleBalanceGetByCount(w http.ResponseWriter, userID *uuid.UUID, count int64) {
	balances, errorResp := logic.GetBalanceByCount(ctx.Db, userID, count)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balances)
	fmt.Println("Retrieved", len(balances), "balances for user ID:", userID)
}

func (ctx *Context) HandleBalanceInsert(w http.ResponseWriter, r *http.Request, userID *uuid.UUID) {
	var entry database.MoneyEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	entry.UserID = *userID
	newEntry, errorResp := logic.InsertBalance(ctx.Db, &entry, userID)
	if errorResp.Code != http.StatusOK {
		http.Error(w, errorResp.Message, errorResp.Code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEntry)
	fmt.Println("Inserted balance with ID:", entry.ID)
}
