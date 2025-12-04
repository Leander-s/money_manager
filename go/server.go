package main

import (
	"fmt"
	"net/http"
	"os"
	"database/sql"
)

type App struct{
	db *sql.DB
}

func initServer() (app *App) {

	dsn := os.Getenv("POSTGRES_DSN")
	fmt.Println("Connecting to database with DSN:", dsn)
	db, err := openDB(dsn)
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
	http.HandleFunc("/", app.rootHandler)
	http.HandleFunc("/budget", app.budgetHandler)
	http.HandleFunc("/balance", app.balanceHandler)
	http.HandleFunc("/user", app.userHandler)

	err := http.ListenAndServe("0.0.0.0:8080", nil)
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
	fmt.Println("Received user", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "User Path Accessed with method:", r.Method)
}
