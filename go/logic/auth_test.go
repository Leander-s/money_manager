package logic

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordGeneration(t *testing.T) {
	password := "password"
	hashedPassword := hashPassword(password)

	bcryptCompareErr := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if bcryptCompareErr != nil {
		t.Errorf("Hashed password did not match original password: %v", bcryptCompareErr)
	}
}
