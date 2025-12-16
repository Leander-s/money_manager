package main

import (
	"fmt"
	"github.com/Leander-s/money_manager/db"
	"github.com/Leander-s/money_manager/api"
	"net/http"
	"os"
)

func initContext() (ctx *api.Context) {

	dsn := os.Getenv("POSTGRES_DSN")
	fmt.Println("Connecting to database with DSN:", dsn)
	db, err := database.OpenDB(dsn)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		panic(err)
	}
	fmt.Println("Successfully connected to the database")
	ctx = &api.Context{
		Db: &db,
	}

	return
}

func runServer(ctx *api.Context) {
	mux := http.NewServeMux()

	mux.Handle("/", ctx.WithAuth(http.HandlerFunc(ctx.RootHandler)))
	mux.Handle("/budget", ctx.WithAuth(http.HandlerFunc(ctx.BudgetHandler)))
	mux.Handle("/balance", ctx.WithAuth(http.HandlerFunc(ctx.BalanceHandler)))
	mux.Handle("/balance/", ctx.WithAuth(http.HandlerFunc(ctx.BalanceHandlerByCount)))
	mux.Handle("/user", ctx.WithAuth(http.HandlerFunc(ctx.UserHandler)))
	mux.Handle("/user/", ctx.WithAuth(http.HandlerFunc(ctx.UserHandlerByID)))
	mux.HandleFunc("/login", ctx.LoginHandler)
	mux.HandleFunc("/register", ctx.RegisterHandler)

	muxWithCORS := withCORS(mux)

	err := http.ListenAndServe("0.0.0.0:8080", muxWithCORS)
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow your Angular dev server
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// handle preflight (OPTIONS) requests quickly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func deinitContext(ctx *api.Context) {
	ctx.Db.Close()
}

