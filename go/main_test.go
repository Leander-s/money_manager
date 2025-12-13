package main

import(
	"testing"
)

func TestPasswordGeneration(t *testing.T) {
	password := "password"
	hashedPassword := hashPassword(password)
	otherHash := hashPassword(password)
	if (hashedPassword != otherHash){
		t.Errorf("Hashed passwords were not equal")
	}
}
