package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	VerifyToken string `yaml:"verify_token`
	AccessToken string `yaml:"access_token"`
	Appsecret   string `yaml:"aap_secret"`
}

func parseContentFile() string {
	contentFile, err := ioutil.ReadFile("content,yml")
	if err != nil {
		log.Printf("Error opening content file: %s\n\n", err)
		panic(err)
	}
	er, _ := yaml.Marshal(contentFile)
	if er != nil {
		log.Printf("Couldn't marshal content fie: %s\n\n", er)
	}
	return string(er)
}

func (c *Config) readYml() *config {
	yamlFile, err := ioutil.ReadFile("bot.config.yml")
	if err != nil {
		log.Printf("Error opening config file: %s\n\n", err)

	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("Couldn't marshal config file:: %s\n", err)
	}
	fmt.Printf("Here is the parsed content.yml: %s\n\n", c)
	return c
}

func getToken() string {
	var c Config
	c.readYml()
	v, err := json.Marshal(c)
	if err != nil {
		log.Printf("Error marshalling our json file:: %s\n", err)
	}
	return string(v)
}
