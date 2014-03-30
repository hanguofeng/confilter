package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/hanguofeng/config"
	"github.com/hanguofeng/confilter"
	"github.com/howbazaar/loggo"
)

var (
	c          *confilter.Confilter
	configFile = flag.String("c", "confilter.conf", "the config file")
	logger     loggo.Logger
)

func JudgeHandler(res http.ResponseWriter, req *http.Request) {
	stime := time.Now()
	req.ParseForm()
	content := req.FormValue("content")
	matches := c.Judge(content)
	jsonBytes, _ := json.Marshal(matches)
	res.Write(jsonBytes)

	timeCost := (time.Now().UnixNano() - stime.UnixNano()) / 1000000
	logger.Infof("cmd:<%s>,remoteAddr:<%s>,timeCost:<%d>,input:<%s>,output:<%s>", "JudgeHandler", req.RemoteAddr, timeCost, req.Form, jsonBytes)
}

func main() {
	initLogger()

	err := initConfig()
	if nil != err {
		logger.Errorf("fatal in init(initConfig),error:%s", err)
		os.Exit(1)
	}

	initConfilter()
	initHttp()

	logger.Infof("Server start...")
	http.ListenAndServe(":8080", nil)
}

func initConfig() error {

	conf, err := config.ReadDefault(*configFile)
	if nil != err {
		return err
	}

	//dictionaries
	dicts := make(map[string]string)
	dictOpts, cfgerr := conf.Options("dictionary")
	if nil != cfgerr {
		return cfgerr
	}
	for _, dictOpt := range dictOpts {
		dictFile, cfgerr := conf.String("dictionary", dictOpt)
		if nil != cfgerr {
			return cfgerr
		}
		dicts[dictOpt] = dictFile
	}
	logger.Infof("dictionaries:%s", dicts)

	return nil

}

func initLogger() {
	logger = loggo.GetLogger("confilter.server")
	logger.SetLogLevel(loggo.DEBUG)
}

func initConfilter() {
	config := new(confilter.ConfilterConfig)
	dicts := make(map[string]string)
	dicts["aaa"] = "../data/aaa.txt"
	dicts["bbb"] = "../data/bbb.txt"
	config.Dictionaries = dicts

	confilterObj, err := confilter.CreateConfilter(config)
	if nil != err {
		logger.Errorf("error when create confilter:%s", err)
		os.Exit(1)
	}
	c = confilterObj
}

func initHttp() {
	http.HandleFunc("/api/confilter/judge", JudgeHandler)
}
