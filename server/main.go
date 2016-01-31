package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hanguofeng/config"
	"github.com/hanguofeng/confilter"
	"github.com/howbazaar/loggo"
)

var (
	c          *confilter.Confilter
	configFile = flag.String("c", "confilter.conf", "the config file")
	logger     loggo.Logger

	cfgDictionaries map[string]string
	cfgLogLevel     string
	cfgListenAddr   string
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
	var err error

	flag.Parse()
	initLogger()

	err = initConfig()
	if nil != err {
		logger.Criticalf("fatal in init(initConfig),error:%s", err)
		os.Exit(1)
	}

	adjustLogLevel()

	err = initConfilter()
	if nil != err {
		logger.Criticalf("fatal in init(initConfilter),error:%s", err)
		os.Exit(1)
	}

	logger.Infof("starting server,listen at <%s>", cfgListenAddr)
	err = initHttp()
	if nil != err {
		logger.Criticalf("fatal in init(initHttp),error:%s", err)
		os.Exit(1)
	}

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
	//logLevel
	loglevel, err := conf.String("server", "logLevel")
	if nil != err {
		loglevel = "UNSPECIFIED"
	}
	//listenAddr
	listenAddr, err := conf.String("server", "listenAddr")
	if nil != err {
		listenAddr = ":80"
	}
	cfgDictionaries = dicts
	cfgLogLevel = loglevel
	cfgListenAddr = listenAddr
	return nil

}

func initLogger() {
	logger = loggo.GetLogger("confilter.server")
}

func adjustLogLevel() {
	var level loggo.Level

	switch strings.ToUpper(cfgLogLevel) {
	case "UNSPECIFIED":
		level = loggo.UNSPECIFIED
	case "TRACE":
		level = loggo.TRACE
	case "DEBUG":
		level = loggo.DEBUG
	case "INFO":
		level = loggo.INFO
	case "WARNING":
		level = loggo.WARNING
	case "ERROR":
		level = loggo.ERROR
	case "CRITICAL":
		level = loggo.CRITICAL
	default:
		level = loggo.UNSPECIFIED
	}
	logger.SetLogLevel(level)
}

func initConfilter() error {
	config := new(confilter.ConfilterConfig)
	config.Dictionaries = cfgDictionaries

	confilterObj, err := confilter.CreateConfilter(config)
	if nil != err {
		return err
	} else {
		c = confilterObj
	}

	return nil
}

func initHttp() error {
	http.HandleFunc("/api/confilter/judge", JudgeHandler)
	return http.ListenAndServe(cfgListenAddr, nil)
}
