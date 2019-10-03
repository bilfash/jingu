package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type (
	Config interface {
		Host() string
		Port() string
		Username() string
		Password() string
		Mailbox() string
		Subjects() []string
		FilePattern() string
		SinkFolder() string
	}
	config map[string]string
)

func New(path string) Config {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read yaml file #%v ", err)
	}

	var configs config
	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		log.Fatalf("cannot unmarshal data: #%v", err)
	}
	return configs
}

func (c config) Host() string {
	return c["host"]
}

func (c config) Port() string {
	return c["port"]
}

func (c config) Username() string {
	return c["username"]
}

func (c config) Password() string {
	return c["password"]
}

func (c config) Mailbox() string {
	return c["mailbox"]
}

func (c config) Subjects() []string {
	return strings.Split(c["subjects"], ",")
}

func (c config) FilePattern() string {
	return c["file_pattern"]
}

func (c config) SinkFolder() string {
	return c["sink_folder"]
}
