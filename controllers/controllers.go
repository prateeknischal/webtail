package controllers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	ctx "github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/prateeknischal/webtail/util"
)

// RootHandler - http handler for handling / path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	absPath, _ := filepath.Abs(util.Conf.Dir)
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		panic(err)
	}
	t := template.New("index").Delims("<<", ">>")
	t, err = t.ParseFiles("templates/index.tmpl")
	t = template.Must(t, err)
	if err != nil {
		panic(err)
	}
	var fileList struct {
		FileList []string
	}

	for _, f := range files {
		/* skip all files that are :
		   d: is a directory
		   a: append-only
		   l: exclusive use
		   T: temporary file; Plan 9 only
		   L: symbolic link
		   D: device file
		   p: named pipe (FIFO)
		   S: Unix domain socket
		   u: setuid
		   g: setgid
		   c: Unix character device, when ModeDevice is set
		   t: sticky
		*/
		if !strings.ContainsAny(f.Mode().String(), "dalTLDpSugct") {
			fileList.FileList = append(fileList.FileList, f.Name())
		}
	}

	t.Execute(w, fileList)
}

// WSHandler - Websocket handler
func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	go tailFile(conn, mux.Vars(r)["file"])
}

// LoginHandler - handles the POST reques to /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	session := ctx.Get(r, "session").(*sessions.Session)
	var isValid = false
	var username = "Anon"
	var err error
	fmt.Println(util.Conf.Dir, util.Conf.ForceAuth)
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
		// flash(w, r, "danger", "Invalid Username/Password")
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
	var csrf = struct {
		Token string
	}{
		Token: csrf.Token(r),
	}
	t.Execute(w, csrf)
}

func tailFile(conn *websocket.Conn, fileName string) {
	currDir, _ := filepath.Abs(util.Conf.Dir)
	// Get the last element of the path
	fileLeafPath := filepath.Base(fileName)
	t, err := tail.TailFile(currDir+string(os.PathSeparator)+fileLeafPath,
		tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Whence: os.SEEK_END,
			},
		})
	if err != nil {
		fmt.Println("Error occurred in opening the file: ", err)
	}
	for line := range t.Lines {
		conn.WriteMessage(websocket.TextMessage, []byte(line.Text))
	}
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
		} else {
			ctx.Set(r, "user", nil)
		}

		// If running on No-Login enforced mode then will set an anon context
		if !util.Conf.ForceAuth {
			ctx.Set(r, "user", "Anon")
		}

		handler.ServeHTTP(w, r)
		// Remove context contents
		ctx.Clear(r)
	}
}
