package logic

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/Leander-s/money_manager/db"
	"github.com/google/uuid"
)

type EntryForUpdate struct {
	ID      uuid.UUID `json:"id"`
	Balance float64   `json:"balance"`
	Ratio   float64   `json:"ratio"`
}

func InsertBalance(store database.MoneyStore, entry *database.MoneyEntry, userID *uuid.UUID) (*database.MoneyEntry, ErrorResponse) {
	lastEntry, _ := GetLastBalance(store, userID)

	entry.Budget = calculateBudget(entry, lastEntry)

	newEntryID, err := store.InsertMoneyDB(entry)
	if err != nil {
		fmt.Println("Error inserting balance:", err)
		return nil, ErrorResponse{
			Message: "Failed to insert balance",
			Code:    http.StatusInternalServerError,
		}
	}

	newEntry, err := store.SelectMoneyByIDDB(&newEntryID)
	if err != nil {
		fmt.Println("Error retrieving new balance:", err)
		return nil, ErrorResponse{
			Message: "Failed to retrieve new balance",
			Code:    http.StatusInternalServerError,
		}
	}

	return newEntry, ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}
}

func calculateBudget(currentBalance *database.MoneyEntry, lastBalance *database.MoneyEntry) float64 {
	// If no last balance, start from 0
	var budget float64 = 0
	if lastBalance == nil {
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

func recalculateBudgets(entries []*database.MoneyEntry) {
	slices.Reverse(entries)

	for i := 1; i < len(entries); i++ {
		entries[i].Budget = calculateBudget(entries[i], entries[i-1])
	}

	slices.Reverse(entries)
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

func updateBalanceEntry(entries []*database.MoneyEntry, updatedEntry *database.MoneyEntry) ([]*database.MoneyEntry, []*database.MoneyEntry) {
	var entriesToUpdate []*database.MoneyEntry
	var index int
	for i, entry := range entries {
		if entry.ID == updatedEntry.ID {
			index = i
		}
	}

	entries[index].Balance = updatedEntry.Balance
	entries[index].Ratio = updatedEntry.Ratio
	entries[index].Budget = 0
	if index != len(entries)-1 {
		entries[index].Budget = calculateBudget(entries[index], entries[index+1])
	}
	entriesToUpdate = entries[:index+1]
	recalculateBudgets(entriesToUpdate)

	remainingEntries := []*database.MoneyEntry{}
	if index != len(entries)-1 {
		remainingEntries = entries[index+1:]
	}

	return append(entriesToUpdate, remainingEntries...), entriesToUpdate
}

func UpdateBalance(store database.MoneyStore, actorID *uuid.UUID, updatedEntry *EntryForUpdate) ([]*database.MoneyEntry, ErrorResponse) {
	fmt.Println("Updating balance entry with ID:", updatedEntry.ID, "to new Balance:", updatedEntry.Balance, "and Ratio:", updatedEntry.Ratio)
	entryToUpdate, errResp := GetBalanceByID(store, &updatedEntry.ID)
	entryToUpdate.Balance = updatedEntry.Balance
	entryToUpdate.Ratio = updatedEntry.Ratio
	if errResp.Code != http.StatusOK {
		return nil, errResp
	}

	if entryToUpdate.UserID != *actorID {
		return nil, ErrorResponse{
			Message: "Forbidden: cannot update another user's balance",
			Code:    http.StatusForbidden,
		}
	}

	entries, errResp := GetAllBalances(store, &entryToUpdate.UserID)

	newEntries, entriesToUpdate := updateBalanceEntry(entries, entryToUpdate)

	err := store.UpdateMoneyBatchDB(entriesToUpdate)
	if err != nil {
		fmt.Println("Error updating balances:", err)
		return nil, ErrorResponse{
			Message: "Failed to update balances",
			Code:    http.StatusInternalServerError,
		}
	}

	fmt.Println("Returning",len(newEntries) ,"updated entries:")
	for _, e := range newEntries {
		fmt.Printf("ID: %s, Balance: %.2f, Budget: %.2f, Ratio: %.2f\n", e.ID, e.Balance, e.Budget, e.Ratio)
	}
	return newEntries, ErrorResponse{
		Message: "",
		Code:    http.StatusOK,
	}
}

func deleteBalanceEntry(entries []*database.MoneyEntry, balanceID *uuid.UUID) ([]*database.MoneyEntry, []*database.MoneyEntry) {
	var index int
	for i, entry := range entries {
		if entry.ID == *balanceID {
			index = i
			break
		}
	}
	if index == 0 {
		return entries[1:], nil
	}

	entriesToUpdate := entries[:index]
	remainingEntries := []*database.MoneyEntry{}
	if index != len(entries)-1 {
		remainingEntries = entries[index+1:]
	}

	entriesToUpdate[len(entriesToUpdate)-1].Budget = 0
	if index != len(entries)-1 {
		entriesToUpdate[len(entriesToUpdate)-1].Budget = calculateBudget(entriesToUpdate[len(entriesToUpdate)-1], entries[index+1])
	}
	recalculateBudgets(entriesToUpdate)

	return append(entriesToUpdate, remainingEntries...), entriesToUpdate
}

func DeleteBalance(store database.MoneyStore, balanceID *uuid.UUID) ([]*database.MoneyEntry, ErrorResponse) {
	errResp := ErrorResponse{}
	entryToDelete, errResp := GetBalanceByID(store, balanceID)
	entries, errResp := GetAllBalances(store, &entryToDelete.UserID)

	newEntries, entriesToUpdate := deleteBalanceEntry(entries, balanceID)

	err := store.DeleteMoneyDB(balanceID)
	if err != nil {
		fmt.Println("Error deleting balance:", err)
		return nil, ErrorResponse{
			Message: "Failed to delete balance",
			Code:    http.StatusInternalServerError,
		}
	}

	if entriesToUpdate == nil {
		return newEntries, errResp
	}

	err = store.UpdateMoneyBatchDB(entriesToUpdate)
	if err != nil {
		fmt.Println("Error updating balances after deletion:", err)
		return nil, ErrorResponse{
			Message: "Failed to update balances after deletion",
			Code:    http.StatusInternalServerError,
		}
	}

	fmt.Println("Returning",len(newEntries) ,"updated entries:")
	for _, e := range newEntries {
		fmt.Printf("ID: %s, Balance: %.2f, Budget: %.2f, Ratio: %.2f\n", e.ID, e.Balance, e.Budget, e.Ratio)
	}

	return newEntries, errResp
}

func GetBalanceByID(store database.MoneyStore, balanceID *uuid.UUID) (*database.MoneyEntry, ErrorResponse) {
	balance, err := store.SelectMoneyByIDDB(balanceID)
	if err != nil {
		fmt.Println("Error retrieving balance:", err)
		return nil, ErrorResponse{
			Message: "Failed to retrieve balance",
			Code:    http.StatusInternalServerError,
		}
	}

	return balance, ErrorResponse{
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
