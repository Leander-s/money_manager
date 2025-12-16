package logic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = bcrypt.DefaultCost

func HandleCreateUser(db *database.Database, w http.ResponseWriter, r *http.Request) {
	var u database.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := CreateUser(db, &u)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		fmt.Println("Error inserting new user:", err)
		return
	}

	u.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
	fmt.Println("Created user with ID:", id)
}

func CreateUser(db *database.Database, user *database.User) (int64, error) {
	user.Password = hashPassword(user.Password)

	id, err := db.InsertUser(user)
	return id, err
}

func HandleGetUsers(db *database.Database, w http.ResponseWriter) {
	users, err := db.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		fmt.Println("Error retrieving users:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	fmt.Println("Retrieved all users")
}

func HandleGetUserByID(db *database.Database, w http.ResponseWriter, id int64) {
	user, err := db.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Println("Error retrieving user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
	fmt.Println("Retrieved user with ID:", id)
}

func HandleUpdateUser(db *database.Database, w http.ResponseWriter, r *http.Request, id int64) {
	var u database.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	u.ID = id

	if err := db.UpdateUser(&u); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		fmt.Println("Error updating user:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
	fmt.Println("Updated user with ID:", id)
}

func HandleDeleteUser(db *database.Database, w http.ResponseWriter, id int64) {
	if err := db.DeleteUser(id); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		fmt.Println("Error deleting user:", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Println("Deleted user with ID:", id)
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return ""
	}
	return string(hash)
}
