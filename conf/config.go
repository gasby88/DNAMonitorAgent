package conf

import (
	"encoding/json"
	"fmt"
	log4 "github.com/alecthomas/log4go"
	"io/ioutil"
	"os"
)

var GCfg = NewConfig()

type Config struct {
	ProcName       string
	Port           string
	RequestPath    string
	StatInterval   int
	ReadTimeout    int
	WriteTimeout   int
	MaxHeaderBytes int
}

func NewConfig() *Config {
	return &Config{}
}

func (this *Config) Init(file string) {
	err := this.load(file)
	if err == nil {
		return
	}
	log4.Error("Config Init error:%s", err)
	this.ProcName = "node"
	this.Port = "8080"
	this.RequestPath = "/dna/monitor"
	this.StatInterval = 3
	this.ReadTimeout = 30
	this.WriteTimeout = 30
	this.MaxHeaderBytes = 104857600

	log4.Info("Config:%+v", this)
}

func (this *Config) load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("NewConfig open file:%s error:%s", file, err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("NewConfig read file:%s error:%s", file, err)
	}
	err = json.Unmarshal(data, this)
	if err != nil {
		return fmt.Errorf("json Unmarshal data:s error:%s", data, err)
	}
	return nil
}
