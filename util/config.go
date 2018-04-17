package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
	// Parse cron
	// Rules for cron :
	// the string should be of type [^0](\d*)(h|d) and the integer should be positive
	// If this exact format is not presented, it will fail.

	timeUnit := cron[len(cron)-1]
	if timeUnit != 'h' && timeUnit != 'd' {
		return fmt.Errorf("Invalid time unit in cron arg: %s", cron)
	}

	timeValue, err := strconv.ParseInt(cron[:len(cron)-1], 10, 32)
	if err != nil {
		return fmt.Errorf("Invalid time value in cron arg: %s", cron)
	}
	if timeValue < 0 {
		return fmt.Errorf("Invalid time value in cron arg: %s", cron)
	}

	if (timeUnit == 'h' && timeValue >= 10000) || (timeUnit == 'd' && timeValue >= 365) {
		fmt.Fprintf(os.Stderr, "Whoah Dude !, That's a long time you put there...")
	}

	// First Index
	IndexFiles(dir)
	tmp := make([]interface{}, len(dir))
	for idx, x := range dir {
		tmp[idx] = x
	}

	// Setting up cron job to keep indexing the files
	if timeValue > 0 {
		repeat := time.Duration(timeValue) * time.Hour

		if timeUnit == 'd' {
			repeat = repeat * 24
		}
		fmt.Println(repeat)
		go MakeAndStartCron(repeat, func(v ...interface{}) error {
			tmp := make([]string, len(v))
			for idx, val := range v {
				tmp[idx] = val.(string)
			}
			IndexFiles(tmp)
			return nil
		}, tmp...)
	}

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
