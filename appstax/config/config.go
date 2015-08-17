package config

import (
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	AppKey     string `json:"appKey"`
	PublicDir  string `json:"publicDir"`
	ServerDir  string `json:"serverDir"`
	ApiBaseUrl string `json:"apiBaseUrl,omitempty"`
}

const fileName = "appstax.conf"

func Exists() bool {
	_, err := ioutil.ReadFile(fileName)
	return err == nil
}

func Write(values map[string]string) {
	config := Read()
	config.AppKey = values["AppKey"]
	config.PublicDir = values["PublicDir"]
	config.ServerDir = values["ServerDir"]
	encoded, err := json.MarshalIndent(config, "", "    ")
	fail.Handle(err)
	ioutil.WriteFile(fileName, encoded, 0644)
	log.Debugf("Wrote config file: %s", encoded)
}

func Read() Config {
	var config Config
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Debugf("Could not find appstax.conf")
	} else {
		err = json.Unmarshal(dat, &config)
		fail.Handle(err)
	}
	insertDefaults(&config)
	return config
}

func insertDefaults(config *Config) {
	if config.PublicDir == "" {
		config.PublicDir = "./public"
	}
	if config.ServerDir == "" {
		config.ServerDir = "./server"
	}
}
