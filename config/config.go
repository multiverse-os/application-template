package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Environment int

const (
	Development Environment = iota
	Testing
	Production
)

func (self Environment) String() string {
	switch self {
	case Testing:
		return "testing"
	case Production:
		return "production"
	default:
		return "development"
	}
}

type Config struct {
	Environment string `yaml:"environment"`
}

func LoadConfig(path string) (config *Config, err error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (self *Config) Save(path string) error {
	configPath, _ := filepath.Split(path)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	} else {
		yamlData, err := yaml.Marshal(&self)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(path, yamlData, 0600)
	}
}

// Initialize - First run config folder structure and file
// initialization using default config, unless otherwise
// specified using flags.
func (self *Config) InitializeConfig(path string) error {
	configPath, _ := filepath.Split(path)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, 0700)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := self.Save(path)
		if err != nil {
			return nil
		}
	}
	return nil
}
