package logic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
)

func HandleGetBalanceByCount(db *database.Database, w http.ResponseWriter, userID int64, c int64) {
	balance, err := db.GetUserMoneyByCount(userID, c)
	if err != nil {
		http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
		fmt.Println("Error retrieving balance:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
	fmt.Println("Retrieved balance for count ID:", c)
}

func HandleGetBalance(db *database.Database, w http.ResponseWriter, userID int64) {
	balance, err := db.GetUserMoney(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve balance", http.StatusInternalServerError)
		fmt.Println("Error retrieving balance:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
	fmt.Println("Retrieved balance for user ID:", userID)
}

func HandleInsertBalance(db *database.Database, w http.ResponseWriter, r *http.Request, userID int64) {
	var entry database.MoneyEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	entry.UserID = userID

	lastEntry, err := db.GetUserMoneyByCount(userID, 1)
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

	id, err := db.InsertMoneyEntry(&entry)
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
