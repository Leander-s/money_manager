package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/model"
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = bcrypt.DefaultCost

func (app *App) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	u.Password = hashPassword(u.Password)

	id, err := app.db.InsertUser(&u)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Println("Error inserting user:", err)
		return
	}

	u.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
	fmt.Println("Created user with ID:", id)
}

func (app *App) handleGetUsers(w http.ResponseWriter) {
	users, err := app.db.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		fmt.Println("Error retrieving users:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	fmt.Println("Retrieved all users")
}

func (app *App) handleGetUserByID(w http.ResponseWriter, id int64) {
	user, err := app.db.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Println("Error retrieving user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	fmt.Println("Retrieved user with ID:", id)
}

func (app *App) handleUpdateUser(w http.ResponseWriter, r *http.Request, id int64) {
	var u model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	u.ID = id

	if err := app.db.UpdateUser(&u); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		fmt.Println("Error updating user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
	fmt.Println("Updated user with ID:", id)
}

func (app *App) handleDeleteUser(w http.ResponseWriter, id int64) {
	if err := app.db.DeleteUser(id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		fmt.Println("Error deleting user:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Println("Deleted user with ID:", id)
}

func hashPassword(password string) string {
	// Placeholder for password hashing logic
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return ""
	}
	return string(hash)
}
