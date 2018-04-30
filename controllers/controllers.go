package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ctx "github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/prateeknischal/webtail/util"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// RootHandler - http handler for handling / path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	t := template.New("index").Delims("<<", ">>")
	t, err := t.ParseFiles("templates/index.tmpl")
	t = template.Must(t, err)
	if err != nil {
		panic(err)
	}

	var fileList = make(map[string]interface{})

	fileList["FileList"] = util.Conf.Dir
	fileList[csrf.TemplateTag] = csrf.Token(r)
	fileList["token"] = csrf.Token(r)
	t.Execute(w, fileList)
}

// WSHandler - Websocket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	filenameB, _ := base64.StdEncoding.DecodeString(mux.Vars(r)["b64file"])
	filename := string(filenameB)
	// sanitize the file if it is present in the index or not.
	filename = filepath.Clean(filename)
	ok := false
	for _, wFile := range util.Conf.Dir {
		if filename == wFile {
			ok = true
			break
		}
	}

	// If the file is found, only then start tailing the file.
	// This is to prevent arbitrary file access. Otherwise send a 403 status
	// This should take care of stacking of filenames as it would first
	// be searched as a string in the index, if not found then rejected.
	if ok {
		go util.TailFile(conn, filename)
	}
	w.WriteHeader(http.StatusUnauthorized)
}

// LoginHandler - handles the POST reques to /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session := ctx.Get(r, "session").(*sessions.Session)
	var isValid = false
	var username = "Anon"
	var err error

	if util.Conf.ForceAuth {
		isValid, username, err = util.Login(r)
		fmt.Println(isValid, username)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Login Failure for %s: %s", username, err)
	}
	if isValid {
		session.Values["id"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/", 302)
	} else {
		session.Save(r, w)
		http.Redirect(w, r, "/login?err=invalid", 302)
	}
}

// LoginPageHandler - GET response to login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if util.Conf.ForceAuth == false {
		http.Redirect(w, r, "/", 302)
	}
	t := template.New("login").Delims("<<", ">>")
	t, err := t.ParseFiles("templates/login.tmpl")
	t = template.Must(t, err)
	if err != nil {
		panic(err)
	}

	t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

// LogoutHandler - handles logout requests
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if util.Conf.ForceAuth == false {
		http.Redirect(w, r, "/", 302)
	}

	session := ctx.Get(r, "session").(*sessions.Session)
	delete(session.Values, "id")
	session.Save(r, w)
	http.Redirect(w, r, "/login?logout=success", 302)
}

// AuthHandler - checks if user is logged in
func AuthHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if u := ctx.Get(r, "user"); u != nil {
			handler.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	}
}

// GetContext wraps each request in a function which fills in the context for a given request.
// This includes setting the User and Session keys and values as necessary for use in later functions.
func GetContext(handler http.Handler) http.HandlerFunc {
	// Set the context here
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing request", http.StatusInternalServerError)
		}
		// Set the context appropriately here.
		session, _ := util.Store.Get(r, "webtail")
		// Put the session in the context so that we can
		// reuse the values in different handlers
		ctx.Set(r, "session", session)
		if id, ok := session.Values["id"]; ok {
			ctx.Set(r, "user", id)
			ctx.Set(r, "isLoggedIn", true)
		} else {
			ctx.Set(r, "user", nil)
			ctx.Set(r, "isLoggedIn", false)
		}

		// If running on No-Login enforced mode then will set an anon context
		if !util.Conf.ForceAuth {
			ctx.Set(r, "user", "Anon")
			ctx.Set(r, "isLoggedIn", false)
		}
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		handler.ServeHTTP(w, r)
		// Remove context contents
		ctx.Clear(r)
	}
}

// UserDetails - returns user name who is logged in
func UserDetails(w http.ResponseWriter, r *http.Request) {
	username := ctx.Get(r, "user").(string)
	isLoggedIn := ctx.Get(r, "isLoggedIn").(bool)
	var resp = struct {
		Username   string `json:"username"`
		IsLoggedIn bool   `json:"isLoggedIn"`
	}{
		Username:   username,
		IsLoggedIn: isLoggedIn,
	}
	json.NewEncoder(w).Encode(resp)
}

// CSRFExemptPrefixes - list of endpoints that does not require csrf protction
var CSRFExemptPrefixes = []string{
	// "/user",
}

// CSRFExceptions - exempts ajax calls from csrf tokens
func CSRFExceptions(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, prefix := range CSRFExemptPrefixes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				r = csrf.UnsafeSkipCheck(r)
				break
			}
		}
		handler.ServeHTTP(w, r)
	}
}
