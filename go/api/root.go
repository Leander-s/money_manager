package api

import (
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
)

type Context struct {
	Db *database.Database
}

func (ctx *Context) RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received root", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Root Path Accessed with method:", r.Method)
	w.WriteHeader(http.StatusOK)
}

