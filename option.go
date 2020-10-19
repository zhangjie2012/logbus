package logbus

import (
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type ServerOption struct {
	Redis struct {
		AppName  string `yaml:"appname"`
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
		ListKey  string `yaml:"listkey"`
	} `yaml:"redis"`
}

var once sync.Once
var option *ServerOption

func ParseOption(fn string) (err error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return
	}
	o := ServerOption{}
	if err = yaml.Unmarshal(bs, &o); err != nil {
		return
	}

	once.Do(func() {
		option = &o
	})
	return
}

func GetOption() *ServerOption {
	return option
}
