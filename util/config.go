package util

// Config - struct to hold the config
type Config struct {
	Dir       string
	ForceAuth bool
}

var Conf Config

// ParseConfig - function to manage config
func ParseConfig(dir string, forceAuth bool) error {
	Conf.Dir = dir
	Conf.ForceAuth = forceAuth

	return nil
}
