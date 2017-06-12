package main

import (
	"DNAMonitorAgent/conf"
	"DNAMonitorAgent/monitor"
	"DNAMonitorAgent/httpServer"
	"flag"
	log4 "github.com/alecthomas/log4go"
	"os"
	"os/signal"
	"syscall"
)

var CfgFile string
var LogFile string

func init() {
	flag.StringVar(&CfgFile, "cf", "./etc/dnamonitor.json", "The path of config file")
	flag.StringVar(&LogFile, "lf", "./etc/log4go.xml", "The path of log config file")
}

func main() {
	log4.LoadConfiguration(LogFile)
	conf.GCfg.Init(CfgFile)

	monitor.MStat = monitor.NewMachineStat(conf.GCfg.StatInterval)
	monitor.MStat.Start()
	defer monitor.MStat.Close()

	server := httpServer.NewHttpServer(conf.GCfg.Port, conf.GCfg.RequestPath, conf.GCfg.ReadTimeout, conf.GCfg.WriteTimeout, conf.GCfg.MaxHeaderBytes)
	server.Start()

	waitToExit()
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-sc
		close(exit)

	}()
	<-exit
}
