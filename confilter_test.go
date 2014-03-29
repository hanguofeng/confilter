package confilter

import (
	"testing"
)

func TestCreate(t *testing.T) {
	_, err := createServer()
	if nil != err {
		t.Fatalf("create server error:%s", err)
		t.Fail()
	}

}

func TestJudge(t *testing.T) {
	server, err := createServer()
	if nil != err {
		t.Fatalf("create server error:%s", err)
	}
	result := server.Judge("inaaa")
	t.Logf("result:%s", result)
	if len(result) < 1 {
		t.Fatalf("result length <1")
		t.Fail()
	}
}

func BenchmarkJudge(b *testing.B) {
	server, err := createServer()
	if nil != err {
		b.Fatalf("create server error:%s", err)
	}
	for i := 0; i < b.N; i++ {
		server.Judge("inainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaainaaaaa")
	}

}

func createServer() (*ConfilterServer, error) {
	config := new(ConfilterServerConfig)
	dicts := make(map[string]string)
	dicts["aaa"] = "./data/aaa.txt"
	dicts["bbb"] = "./data/bbb.txt"
	config.Dictionaries = dicts

	return CreateConfilterServer(config)

}
