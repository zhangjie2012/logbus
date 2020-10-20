package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zhangjie2012/logbus"
	"gopkg.in/yaml.v2"
)

var config string

func init() {
	flag.StringVar(&config, "config", "/etc/logbus.yml", "the server config")
}

// redis:
//   appname:
//   addr:
//   password:
//   db:
//   listkey:
// mongo:
//   host:
//   port:
//   username:
//   password:
//   dbname:
type ServerOption struct {
	Redis struct {
		AppName  string `yaml:"appname"`
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
		ListKey  string `yaml:"listkey"`
	} `yaml:"redis"`
	Mongo struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbname"`
	} `yaml:"mongo"`
}

func ParseOption(fn string) (*ServerOption, error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	o := ServerOption{}
	if err = yaml.Unmarshal(bs, &o); err != nil {
		return nil, err
	}

	return &o, nil
}

func main() {
	flag.Parse()

	option, err := ParseOption(config)
	if err != nil {
		log.Printf("parse option failure, config=%s, error=%s", config, err)
		return
	}

	in, err := logbus.NewRedisListInput(
		option.Redis.AppName,
		option.Redis.Addr,
		option.Redis.Password,
		option.Redis.Db,
		option.Redis.ListKey,
	)
	if err != nil {
		log.Printf("init cache failure, error=%s", err)
		return
	}
	defer in.Close()

	stdoutOut, err := logbus.NewStdoutOutput(logbus.DefaultTransformer)
	if err != nil {
		log.Printf("new stdout output failure, error=%s", err)
		return
	}
	defer stdoutOut.Close()

	mongoOut, err := logbus.NewMongoOutput(
		option.Mongo.Host,
		option.Mongo.Port,
		option.Mongo.Username,
		option.Mongo.Password,
		option.Mongo.DbName,
		logbus.StatLogTransformer,
	)
	if err != nil {
		log.Printf("new mongo output failure, error=%s", err)
		return
	}
	defer mongoOut.Close()

	wg := sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		go logbus.Serve(ctx, in, []logbus.Output{stdoutOut, mongoOut})
	}()

	log.Println("logbus start")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("receive stop singal, stop ...")
}
