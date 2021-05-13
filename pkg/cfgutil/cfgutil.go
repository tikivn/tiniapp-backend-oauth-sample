package cfgutil

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func LoadConfig(path string, config interface{}) error {
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if err := ValidateConfigPath(path); err != nil {
		return err
	}

	tmpl, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	tmpl = []byte(os.ExpandEnv(string(tmpl)))

	err = yaml.Unmarshal(tmpl, config)
	if err != nil {
		return err
	}

	return nil
}
