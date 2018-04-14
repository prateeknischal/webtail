package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prateeknischal/webtail/controllers"
	"github.com/prateeknischal/webtail/util"
)

var (
	dir       = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./").ExistingFilesOrDirs()
	port      = kingpin.Flag("port", "Port number to host the server").Short('p').Default("8080").Int()
	restrict  = kingpin.Flag("restrict", "Enforce PAM authentication (single level)").Short('r').Bool()
	whitelist = kingpin.Flag("whitelist", "enable whitelisting with users in the provided file").Short('w').ExistingFile()
	cron      = kingpin.Flag("cron", "configure cron for re-indexing files (Not supported right now)").Short('c').Default("1h").String()
)

func main() {
	kingpin.Parse()
	_ = util.ParseConfig(*dir, *restrict, *whitelist, *cron)

	router := mux.NewRouter()

	router.HandleFunc("/ws/{b64file}", Use(controllers.WSHandler, controllers.AuthHandler, controllers.GetContext)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler, controllers.AuthHandler, controllers.GetContext)).Methods("GET")
	router.HandleFunc("/login", Use(controllers.LoginHandler, controllers.GetContext)).Methods("POST")
	router.HandleFunc("/login", Use(controllers.LoginPageHandler, controllers.GetContext)).Methods("GET")
	router.HandleFunc("/logout", Use(controllers.LogoutHandler, controllers.AuthHandler, controllers.GetContext)).Methods("POST")
	router.HandleFunc("/user", Use(controllers.UserDetails, controllers.AuthHandler, controllers.GetContext))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	csrfHandler := csrf.Protect([]byte(util.GenerateSecureKey()),
		csrf.FieldName("csrf_token"),
		csrf.Secure(false))
	csrfRouter := Use(csrfHandler(router).ServeHTTP, controllers.CSRFExceptions)
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", *port), Handler: handlers.CombinedLoggingHandler(os.Stdout, csrfRouter)}
	panic(server.ListenAndServe())
}

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
