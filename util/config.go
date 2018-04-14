package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config - struct to hold the config
type Config struct {
	Dir       []string
	ForceAuth bool
	Whitelist map[string]bool
	Cron      string
}

// Conf global config
var Conf Config

// ParseConfig - function to manage config
func ParseConfig(dir []string, forceAuth bool, whitelist, cron string) error {
	IndexFiles(dir)
	Conf.Dir = FileList
	Conf.ForceAuth = forceAuth
	if err := getWhitelistFromFile(whitelist); err != nil {
		return err
	}
	Conf.Cron = cron
	return nil
}

func getWhitelistFromFile(file string) error {
	if len(file) == 0 {
		return nil
	}
	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error Opening file %s; err: %s", file, err)
		return err
	}
	reader := bufio.NewReader(f)
	s, err := reader.ReadString('\n')
	for err != nil {
		Conf.Whitelist[strings.Trim(s, "\n\r ")] = true
	}
	return err
}

// IsWhitelisted - returns if a username is in ACL or not
// If the ACL list is empty, it assumes no ACL
func IsWhitelisted(username string) bool {
	// if the Whitelist is empty, library assumes no ACL passed
	if len(Conf.Whitelist) == 0 {
		return true
	}
	_, ok := Conf.Whitelist[username]
	return ok
}
