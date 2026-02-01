package logic

import (
	"testing"

	"github.com/Leander-s/money_manager/db"
	"github.com/google/uuid"
)

func TestRecalculateBudget_UpdatedLastEntry(t *testing.T) {
	ID0 := uuid.New()
	ID1 := uuid.New()
	ID2 := uuid.New()
	ID3 := uuid.New()

	balanceEntries := []*database.MoneyEntry{
		{
			ID:      ID3,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  -200.0,
		},
		{
			ID:      ID2,
			Balance: 1200.0,
			Ratio:   0.5,
			Budget:  200.0,
		},
		{
			ID:      ID1,
			Balance: 1000.0,
			Ratio:   0.5,
			Budget:  100.0,
		},
		{
			ID:      ID0,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
	}

	updatedEntry := &database.MoneyEntry{
		ID:      ID3,
		Balance: 1300.0,
		Ratio:   0.5,
	}

	expectedBudgets := []float64{250.0, 200.0, 100.0, 0.0}

	newBalances, _ := updateBalanceEntry(balanceEntries, updatedEntry)

	for i, entry := range newBalances {
		if entry.Budget != expectedBudgets[i] {
			t.Errorf("Expected budget %.2f, but got %.2f", expectedBudgets[i], entry.Budget)
		}
	}
}

func TestRecalculateBudget_UpdatedMidEntry(t *testing.T) {
	ID0 := uuid.New()
	ID1 := uuid.New()
	ID2 := uuid.New()
	ID3 := uuid.New()

	balanceEntries := []*database.MoneyEntry{
		{
			ID:      ID3,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  -200.0,
		},
		{
			ID:      ID2,
			Balance: 1200.0,
			Ratio:   0.5,
			Budget:  200.0,
		},
		{
			ID:      ID1,
			Balance: 1000.0,
			Ratio:   0.5,
			Budget:  100.0,
		},
		{
			ID:      ID0,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
	}

	updatedEntry := &database.MoneyEntry{
		ID:      ID1,
		Balance: 1300.0,
		Ratio:   0.5,
	}

	expectedBudgets := []float64{-250.0, 150.0, 250.0, 0.0}

	newBalances, _ := updateBalanceEntry(balanceEntries, updatedEntry)

	for i, entry := range newBalances {
		if entry.Budget != expectedBudgets[i] {
			t.Errorf("Expected budget %.2f, but got %.2f", expectedBudgets[i], entry.Budget)
		}
	}
}

func TestRecalculateBudget_UpdatedFirstEntry(t *testing.T) {
	ID0 := uuid.New()
	ID1 := uuid.New()
	ID2 := uuid.New()
	ID3 := uuid.New()

	balanceEntries := []*database.MoneyEntry{
		{
			ID:      ID3,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  -200.0,
		},
		{
			ID:      ID2,
			Balance: 1200.0,
			Ratio:   0.5,
			Budget:  200.0,
		},
		{
			ID:      ID1,
			Balance: 1000.0,
			Ratio:   0.5,
			Budget:  100.0,
		},
		{
			ID:      ID0,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
	}

	updatedEntry := &database.MoneyEntry{
		ID:      ID0,
		Balance: 600.0,
		Ratio:   0.5,
	}

	expectedBudgets := []float64{-100.0, 300.0, 200.0, 0.0}

	newBalances, _ := updateBalanceEntry(balanceEntries, updatedEntry)

	for i, entry := range newBalances {
		if entry.Budget != expectedBudgets[i] {
			t.Errorf("Expected budget %.2f, but got %.2f", expectedBudgets[i], entry.Budget)
		}
	}
}

func TestRecalulateBudget_DeletedLastEntry(t *testing.T) {
	ID0 := uuid.New()
	ID1 := uuid.New()
	ID2 := uuid.New()
	ID3 := uuid.New()

	balanceEntries := []*database.MoneyEntry{
		{
			ID:      ID3,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
		{
			ID:      ID2,
			Balance: 1200.0,
			Ratio:   0.5,
			Budget:  200.0,
		},
		{
			ID:      ID1,
			Balance: 1000.0,
			Ratio:   0.5,
			Budget:  100.0,
		},
		{
			ID:      ID0,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
	}

	expectedBudgets := []float64{200.0, 100.0, 0.0}

	newBalances, _ := deleteBalanceEntry(balanceEntries, &ID3)

	for i, entry := range newBalances {
		if entry.Budget != expectedBudgets[i] {
			t.Errorf("Expected budget %.2f, but got %.2f", expectedBudgets[i], entry.Budget)
		}
	}
}

func TestRecalulateBudget_DeletedMidEntry(t *testing.T) {
	ID0 := uuid.New()
	ID1 := uuid.New()
	ID2 := uuid.New()
	ID3 := uuid.New()

	balanceEntries := []*database.MoneyEntry{
		{
			ID:      ID3,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
		{
			ID:      ID2,
			Balance: 1200.0,
			Ratio:   0.5,
			Budget:  200.0,
		},
		{
			ID:      ID1,
			Balance: 1000.0,
			Ratio:   0.5,
			Budget:  100.0,
		},
		{
			ID:      ID0,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
	}

	expectedBudgets := []float64{-100.0, 100.0, 0.0}

	newBalances, _ := deleteBalanceEntry(balanceEntries, &ID2)

	for i, entry := range newBalances {
		if entry.Budget != expectedBudgets[i] {
			t.Errorf("Expected budget %.2f, but got %.2f", expectedBudgets[i], entry.Budget)
		}
	}
}

func TestRecalulateBudget_DeletedFirstEntry(t *testing.T) {
	ID0 := uuid.New()
	ID1 := uuid.New()
	ID2 := uuid.New()
	ID3 := uuid.New()

	balanceEntries := []*database.MoneyEntry{
		{
			ID:      ID3,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
		{
			ID:      ID2,
			Balance: 1200.0,
			Ratio:   0.5,
			Budget:  200.0,
		},
		{
			ID:      ID1,
			Balance: 1000.0,
			Ratio:   0.5,
			Budget:  100.0,
		},
		{
			ID:      ID0,
			Balance: 800.0,
			Ratio:   0.5,
			Budget:  0.0,
		},
	}

	expectedBudgets := []float64{-300.0, 100.0, 0.0}

	newBalances, _ := deleteBalanceEntry(balanceEntries, &ID0)

	for i, entry := range newBalances {
		if entry.Budget != expectedBudgets[i] {
			t.Errorf("Expected budget %.2f, but got %.2f", expectedBudgets[i], entry.Budget)
		}
	}
}

func TestCalculateBudgetIncrease(t *testing.T) {
	testBalance := &database.MoneyEntry{
		Balance: 1000.0,
		Ratio: 0.2,
	}
	
	testLastBalance := &database.MoneyEntry{
		Balance: 800.0,
		Budget: 400.0,
		Ratio: 0.5,
	}

	expectedBudget := 440.0

	calculatedBudget := calculateBudget(testBalance, testLastBalance)

	if calculatedBudget != expectedBudget {
		t.Errorf("Expected budget %.2f, but got %.2f", expectedBudget, calculatedBudget)
	}
}

func TestCalculateBudgetDecrease(t *testing.T) {
	testBalance := &database.MoneyEntry{
		Balance: 600.0,
		Ratio: 0.9,
	}
	
	testLastBalance := &database.MoneyEntry{
		Balance: 800.0,
		Budget: 400.0,
		Ratio: 0.5,
	}

	expectedBudget := 200.0

	calculatedBudget := calculateBudget(testBalance, testLastBalance)

	if calculatedBudget != expectedBudget {
		t.Errorf("Expected budget %.2f, but got %.2f", expectedBudget, calculatedBudget)
	}
}
