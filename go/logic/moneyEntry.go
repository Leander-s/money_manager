package logic

import (
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
	"github.com/google/uuid"
)

func InsertBalance(store database.MoneyStore, entry *database.MoneyEntry, userID *uuid.UUID) (*database.MoneyEntry, ErrorResponse) {
	lastEntry, _ := GetLastBalance(store, userID)

	entry.Budget = calculateBudget(entry, lastEntry)

	id, err := store.InsertMoneyDB(entry)
	if err != nil {
		fmt.Println("Error inserting balance:", err)
		return nil, ErrorResponse{
			Message: "Failed to insert balance",
			Code:    http.StatusInternalServerError,
		}
	}

	entry.ID = id
	return entry, ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}
}

func calculateBudget(currentBalance *database.MoneyEntry, lastBalance *database.MoneyEntry) float64 {
	// If this is the first balance entry, set budget based on ratio
	var budget = currentBalance.Balance * currentBalance.Ratio
	if (lastBalance == nil) {
		return budget
	}

	// Otherwise, start from last budget
	budget = lastBalance.Budget

	diff := currentBalance.Balance - lastBalance.Balance

	// If balance decreased, subtract from budget
	if diff < 0 {
		return budget + diff
	}
	// If balance increased, add to budget based on ratio
	return lastBalance.Budget + diff*currentBalance.Ratio
}

func GetLastBalance(store database.MoneyStore, userID *uuid.UUID) (*database.MoneyEntry, ErrorResponse) {
	balances, err := store.SelectUserMoneyByCountDB(userID, 1)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return nil, ErrorResponse{
			Message: "Failed to retrieve balance",
			Code:    http.StatusInternalServerError,
		}
	}

	if len(balances) == 0 {
		return nil, ErrorResponse{
			Message: "No balance entries found",
			Code:    http.StatusNotFound,
		}
	}

	return balances[0], ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}
}

func GetBalanceByCount(store database.MoneyStore, userID *uuid.UUID, count int64) ([]*database.MoneyEntry, ErrorResponse) {
	balances, err := store.SelectUserMoneyByCountDB(userID, count)
	if err != nil {
		fmt.Println("Error retrieving balances:", err)
		return nil, ErrorResponse{
			Message: "Failed to retrieve balances",
			Code:    http.StatusInternalServerError,
		}
	}

	return balances, ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}
}

func GetAllBalances(store database.MoneyStore, userID *uuid.UUID) ([]*database.MoneyEntry, ErrorResponse) {
	balances, err := store.SelectUserMoneyDB(userID)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return nil, ErrorResponse{
			Message: "Failed to retrieve balances",
			Code:    http.StatusInternalServerError,
		}
	}

	return balances, ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}
}
