package api

import (
	"fmt"
	"net/http"

	"github.com/Leander-s/money_manager/db"
	"github.com/Leander-s/money_manager/logic"
)

type Context struct {
	// Database connection
	Db             database.DatabaseInterface
	// CORS allowed origins
	AllowedOrigins string
	// Email configuration
	MailConfig     logic.EmailSender
	// Host address for links
	HostAddress    string
	// Frontend address for links
	FronendAddress string
	// Flag indicating if there are no users in the database
	NoUsers        bool
}

func (ctx *Context) RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received root", r.Method, "request from:", r.RemoteAddr)
	fmt.Fprintln(w, "Root Path Accessed with method:", r.Method)
	w.WriteHeader(http.StatusOK)
}
