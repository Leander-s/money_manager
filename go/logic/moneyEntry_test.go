package logic

import (
	"testing"

	"github.com/Leander-s/money_manager/db"
)

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
