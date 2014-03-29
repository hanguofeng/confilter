package confilter

import (
	"bytes"
	"io/ioutil"

	"github.com/cloudflare/ahocorasick"
)

type ConfilterServerConfig struct {
	Dictionaries map[string]string
}

type ConfilterServer struct {
	config   *ConfilterServerConfig
	matchers map[string]*ahocorasick.Matcher
}

func CreateConfilterServer(config *ConfilterServerConfig) (*ConfilterServer, error) {
	server := new(ConfilterServer)
	server.config = config
	err := server.initMatchers()
	if nil != err {
		return nil, err
	}
	return server, nil
}

func (this *ConfilterServer) initMatchers() error {
	matchers := make(map[string]*ahocorasick.Matcher)
	for k, path := range this.config.Dictionaries {
		dictContent, err := ioutil.ReadFile(path)
		if nil != err {
			return err
		}
		dictContentSplited := bytes.Split(dictContent, []byte("\n"))
		matcher := ahocorasick.NewMatcher(dictContentSplited)
		matchers[k] = matcher
	}

	this.matchers = matchers
	return nil
}

func (this *ConfilterServer) Judge(s string) (matches []string) {
	sbytes := []byte(s)
	matches = this.JudgeBytes(sbytes)
	return matches
}

func (this *ConfilterServer) JudgeBytes(s []byte) (matches []string) {
	matches = []string{}
	for k, matcher := range this.matchers {
		hits := matcher.Match(s)
		if len(hits) > 0 {
			matches = append(matches, k)
		}
	}
	return matches
}
