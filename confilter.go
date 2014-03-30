package confilter

import (
	"bytes"
	"io/ioutil"

	"github.com/cloudflare/ahocorasick"
)

type ConfilterConfig struct {
	Dictionaries map[string]string
}

type Confilter struct {
	config   *ConfilterConfig
	matchers map[string]*ahocorasick.Matcher
}

func CreateConfilter(config *ConfilterConfig) (*Confilter, error) {
	server := new(Confilter)
	server.config = config
	err := server.InitMatchers()
	if nil != err {
		return nil, err
	}
	return server, nil
}

func (this *Confilter) InitMatchers() error {
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

	needDoReplace := true

	if len(this.matchers) > 0 {
		if len(matchers) == 0 {
			//when there are running matchers,it will not be replace the well running matcher to empty 
			needDoReplace = false
		}
	}

	if needDoReplace {
		this.matchers = matchers
	}

	return nil
}

func (this *Confilter) Judge(s string) (matches []string) {
	sbytes := []byte(s)
	matches = this.JudgeBytes(sbytes)
	return matches
}

func (this *Confilter) JudgeBytes(s []byte) (matches []string) {
	matches = []string{}
	for k, matcher := range this.matchers {
		hits := matcher.Match(s)
		if len(hits) > 0 {
			matches = append(matches, k)
		}
	}
	return matches
}
