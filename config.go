package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type config struct {
	Port                uint          `json:"port"`
	TimeoutMs           time.Duration `json:"timeoutMs"`
	ReadHeaderTimeoutMs time.Duration `json:"readHeaderTimeoutMs"`
	ReadTimeoutMs       time.Duration `json:"readTimeoutMs"`
	Root                string        `json:"root"`
	PathPrefix          string        `json:"pathPrefix"`
	LogRequests         *bool         `json:"logRequests,omitempty"`
	HeaderServerName    string        `json:"headerServerName"`
}

var ptrTrue bool = true

var defaultValues config = config{
	Port:                3000,
	TimeoutMs:           1000,
	ReadHeaderTimeoutMs: 1000,
	ReadTimeoutMs:       1000,
	Root:                "/var/www",
	PathPrefix:          "",
	LogRequests:         &ptrTrue,
	HeaderServerName:    "static-lite-server:" + version,
}

func fillDefaultValues(c *config) {
	if c.Port <= 0 {
		c.Port = defaultValues.Port
	}
	if c.TimeoutMs <= 0 {
		c.TimeoutMs = defaultValues.TimeoutMs
	}
	if c.ReadHeaderTimeoutMs <= 0 {
		c.ReadHeaderTimeoutMs = defaultValues.ReadHeaderTimeoutMs
	}
	if c.ReadTimeoutMs <= 0 {
		c.ReadTimeoutMs = defaultValues.ReadTimeoutMs
	}
	if c.Root == "" {
		c.Root = defaultValues.Root
	}
	if c.LogRequests == nil {
		c.LogRequests = defaultValues.LogRequests
	}
	if c.HeaderServerName == "" {
		c.HeaderServerName = defaultValues.HeaderServerName
	}
}

func parseConfig(configPath string) (config, error) {
	configFile, err := os.Open(configPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println("Cannot open config at", configPath, err)
		return config{}, err
	}

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Println("Cannot close config: ", err)
		}
	}(configFile)

	byteValue, err := ioutil.ReadAll(configFile)

	if err != nil {
		log.Println("Cannot read config: ", err)
		return config{}, err
	}

	var c config
	err = json.Unmarshal([]byte(byteValue), &c)

	if err != nil {
		log.Println("Cannot parse config: ", err)
		return config{}, err
	}

	fillDefaultValues(&c)

	return c, nil
}
