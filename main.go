package main

import (
	"crypto/tls"
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
	dir      = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./").ExistingFilesOrDirs()
	port     = kingpin.Flag("port", "Port number to host the server").Short('p').Default("8080").Int()
	restrict = kingpin.Flag("restrict", "Enforce PAM authentication (single level)").Short('r').Bool()
	acl      = kingpin.Flag("acl", "enable Access Control List with users in the provided file").Short('a').ExistingFile()
	cron     = kingpin.Flag("cron", "configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]").Short('t').Default("0h").String()
	secure   = kingpin.Flag("secure", "Run Server with TLS").Short('s').Bool()
	cert     = kingpin.Flag("cert", "Server Certificate").Short('c').Default("server.crt").String()
	key      = kingpin.Flag("key", "Server Key File").Short('k').Default("server.key").String()
)

func main() {
	kingpin.Parse()
	err := util.ParseConfig(*dir, *restrict, *acl, *cron)

	if err != nil {
		panic(err)
	}

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

	if *secure == false {
		server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", *port), Handler: handlers.CombinedLoggingHandler(os.Stdout, csrfRouter)}
		panic(server.ListenAndServe())
	} else {
		serverCert, err := tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read the certificates %s", err)
			panic(err)
		}
		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS11,
			Certificates:             []tls.Certificate{serverCert},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		server := &http.Server{
			Addr:      fmt.Sprintf("0.0.0.0:%d", *port),
			Handler:   handlers.CombinedLoggingHandler(os.Stdout, csrfRouter),
			TLSConfig: tlsConfig,
		}

		panic(server.ListenAndServeTLS(*cert, *key))
	}
}

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
