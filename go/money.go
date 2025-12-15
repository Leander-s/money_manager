package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/model"
)

func (app *App) handleGetBalanceByCount(w http.ResponseWriter, userID int64, c int64) {
	balance, err := app.db.GetUserMoneyByCount(userID, c)
	if err != nil {
		http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
		fmt.Println("Error retrieving balance:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
	fmt.Println("Retrieved balance for count ID:", c)
}

func (app *App) handleGetBalance(w http.ResponseWriter, userID int64) {
	balance, err := app.db.GetUserMoney(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
		fmt.Println("Error retrieving balance:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
	fmt.Println("Retrieved balance for user ID:", userID)
}

func (app *App) handleInsertBalance(w http.ResponseWriter, r *http.Request, userID int64) {
	var entry model.MoneyEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	entry.UserID = userID

	lastEntry, err := app.db.GetUserMoneyByCount(userID, 1)
	if err != nil {
		http.Error(w, "Failed to retrieve last balance", http.StatusInternalServerError)
		fmt.Println("Error retrieving last balance:", err)
		return
	}

	entry.Budget = entry.Balance * entry.Ratio

	if len(lastEntry) > 0 {
		diff := entry.Balance - lastEntry[0].Balance
		entry.Budget = lastEntry[0].Budget + diff*entry.Ratio
	}

	id, err := app.db.InsertMoneyEntry(&entry)
	if err != nil {
		http.Error(w, "Failed to insert balance", http.StatusInternalServerError)
		fmt.Println("Error inserting balance:", err)
		return
	}

	entry.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
	fmt.Println("Inserted balance with ID:", id)
}
