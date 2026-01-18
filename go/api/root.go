package api

import (
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
	"github.com/Leander-s/money_manager/logic"
)

type Context struct {
	Db             *database.Database
	AllowedOrigins string
	MailConfig     *logic.BrevoConfig
	HostAddress    string
}

func (ctx *Context) RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received root", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Root Path Accessed with method:", r.Method)
	w.WriteHeader(http.StatusOK)
}
