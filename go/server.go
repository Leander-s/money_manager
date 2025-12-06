package main

import (
	"strings"
	"strconv"
	"fmt"
	"net/http"
	"os"
	"github.com/Leander-s/money_manager/model"
)

type App struct{
	db model.Database
}

func initServer() (app *App) {

	dsn := os.Getenv("POSTGRES_DSN")
	fmt.Println("Connecting to database with DSN:", dsn)
	db, err := model.OpenDB(dsn)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		panic(err)
	}
	fmt.Println("Successfully connected to the database")
	app = &App{
		db: db,
	}

	return
}

func (app *App) runServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.rootHandler)
	mux.HandleFunc("/budget", app.budgetHandler)
	mux.HandleFunc("/balance", app.balanceHandler)
	mux.HandleFunc("/user", app.userHandler)
	mux.HandleFunc("/user/", app.userHandlerByID)

	err := http.ListenAndServe("0.0.0.0:8080", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}
}

func (app *App) deInitServer() {
	app.db.Close()
}

func (app *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received root", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Root Path Accessed with method:", r.Method)
}

func (app *App) balanceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received budget", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Balance Path Accessed with method:", r.Method)
}

func (app *App) budgetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received balance", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Budget Path Accessed with method:", r.Method)
}

func (app *App) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.handleCreateUser(w, r)
	case http.MethodGet:
		app.handleGetUsers(w)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (app *App) userHandlerByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/user/")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var id int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		app.handleGetUserByID(w, id)
	case http.MethodPut:
		app.handleUpdateUser(w, r, id)
	case http.MethodDelete:
		app.handleDeleteUser(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
