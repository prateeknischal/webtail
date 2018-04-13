package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prateeknischal/webtail/controllers"
	"github.com/prateeknischal/webtail/util"
)

var (
	dir      = flag.String("dir", "./", "Log directory path")
	restrict = flag.Bool("restrict", false, "User authentication to allow access")
)

func main() {
	flag.Parse()
	_ = util.ParseConfig(*dir, *restrict)

	router := mux.NewRouter()

	router.HandleFunc("/ws/{file}", Use(controllers.WSHandler, controllers.AuthHandler, controllers.GetContext)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler, controllers.AuthHandler, controllers.GetContext)).Methods("GET")
	router.HandleFunc("/login", Use(controllers.LoginHandler, controllers.GetContext)).Methods("POST")
	router.HandleFunc("/login", Use(controllers.LoginPageHandler, controllers.GetContext)).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	csrfHandler := csrf.Protect([]byte(util.GenerateSecureKey()),
		csrf.FieldName("csrf_token"),
		csrf.Secure(false))
	csrfRouter := csrfHandler(router)
	server := &http.Server{Addr: ":8080", Handler: handlers.CombinedLoggingHandler(os.Stdout, csrfRouter)}
	panic(server.ListenAndServe())
}

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
