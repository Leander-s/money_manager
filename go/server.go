package main

import (
	"context"
	"fmt"
	"github.com/Leander-s/money_manager/model"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type App struct {
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
	mux.Handle("/budget", app.withAuth(http.HandlerFunc(app.budgetHandler)))
	mux.Handle("/balance", app.withAuth(http.HandlerFunc(app.balanceHandler)))
	mux.Handle("/user", app.withAuth(http.HandlerFunc(app.userHandler)))
	mux.Handle("/user/", app.withAuth(http.HandlerFunc(app.userHandlerByID)))
	mux.HandleFunc("/login", app.handleLogin)
	mux.HandleFunc("/createAccount", app.handleCreateAccount)

	muxWithCORS := withCORS(mux)

	err := http.ListenAndServe("0.0.0.0:8080", muxWithCORS)
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}
}

func (app *App) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		auth = strings.TrimPrefix(auth, "Bearer ")
		userID, err := app.validateToken(auth)
		if(err != nil) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If authentication succeeds, proceed to the next handler
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow your Angular dev server
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// handle preflight (OPTIONS) requests quickly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
