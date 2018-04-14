package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
)

var (
	// FileList - list of files that were parsed from the provided config
	FileList []string
	visited  map[string]bool
)

// TailFile - Accepts a websocket connection and a filename and tails the
// file and writes the changes into the connection. Recommended to run on
// a thread as this is blocking in nature
func TailFile(conn *websocket.Conn, fileName string) {
	t, err := tail.TailFile(fileName,
		tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Whence: os.SEEK_END,
			},
		})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error occurred in opening the file: ", err)
	}
	for line := range t.Lines {
		conn.WriteMessage(websocket.TextMessage, []byte(line.Text))
	}
}

// IndexFiles - takes argument as a list of files and directories and returns
// a list of absolute file strings to be tailed
func IndexFiles(fileList []string) {
	for _, file := range fileList {
		dfs(file)
	}
	fmt.Fprintln(os.Stderr, "Indexing complete !")
	for _, f := range FileList {
		fmt.Fprintln(os.Stderr, f)
	}
}

/* skip all files that are :
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
func dfs(file string) {
	file = filepath.Clean(file)
	absPath, err := filepath.Abs(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get absolute path of the file %s; err: %s\n", file, err)
	}
	if _, ok := visited[file]; ok {
		// if the absolute path has been visited, return without processing
		return
	}
	s, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to stat file %s; err: %s\n", file, err)
		return
	}
	// check if the file is a directory
	if s.IsDir() {
		basepath := filepath.Clean(file)
		filelist, _ := ioutil.ReadDir(absPath)
		for _, f := range filelist {
			dfs(basepath + string(os.PathSeparator) + f.Name())
		}
	} else if strings.ContainsAny(s.Mode().String(), "alTLDpSugct") {
		// skip these files
		// @TODO try including names PIPES
	} else {
		// only remaining file are ascii files that can be then differentiated
		// by the user as golang has only these many categorization
		// Note : this appends the absolute paths
		FileList = append(FileList, absPath)
	}
}
