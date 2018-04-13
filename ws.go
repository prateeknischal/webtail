package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
)

var (
	dir      = flag.String("dir", "./", "Log directory path")
	restrict = flag.Bool("restrict", false, "User authentication to allow access")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/ws/{file}", wsHandler).Methods("GET")
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	server := &http.Server{Addr: ":8080", Handler: handlers.CombinedLoggingHandler(os.Stdout, router)}
	panic(server.ListenAndServe())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	absPath, _ := filepath.Abs(*dir)
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	go echo(conn, mux.Vars(r)["file"])
}

func echo(conn *websocket.Conn, fileName string) {
	currDir, _ := filepath.Abs(*dir)
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
