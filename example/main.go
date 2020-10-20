package main

import (
	"flag"
	"fmt"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/zhangjie2012/cbl-go/cache"
	"github.com/zhangjie2012/logbus"
)

var config string

func init() {
	flag.StringVar(&config, "config", "/etc/logbus.yml", "the server config")
}

func main() {
	flag.Parse()

	if err := logbus.ParseOption(config); err != nil {
		log.Fatalf("parse option failure, config=%s, error=%s", config, err)
	}

	option := logbus.GetOption()

	if err := cache.InitCache(option.Redis.AppName, option.Redis.Addr, option.Redis.Password, option.Redis.Db); err != nil {
		log.Fatalf("init cache failure, error=%s", err)
	}

	metadata := logrus.Fields{}
	stateid, _ := metadata["stateid"].(string)

	fmt.Println(len(stateid))
}
