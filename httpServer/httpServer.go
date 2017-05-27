package httpServer

import (
	. "DNAMonitorAgent/monitor"
	"encoding/json"
	"fmt"
	log4 "github.com/alecthomas/log4go"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type HttpServer struct {
	httpServer *http.Server
}

func NewHttpServer(port string, readTimeout, writeTimeout int, maxHeaderBytes int) *HttpServer {
	server := &HttpServer{}
	httpRouter := httprouter.New()
	httpRouter.GET("/dna/monitor/cpu", server.GetCpuStat)
	httpRouter.GET("/dna/monitor/mem", server.GetMemStat)
	httpRouter.GET("/dna/monitor/dis", server.GetDisStat)
	httpRouter.GET("/dna/monitor/net", server.GetNetStat)
	httpRouter.GET("/dna/monitor/mechine", server.GetMachineStat)
	httpRouter.NotFound = &NotFound{}

	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        httpRouter,
		ReadTimeout:    time.Second * time.Duration(readTimeout),
		WriteTimeout:   time.Second * time.Duration(writeTimeout),
		MaxHeaderBytes: maxHeaderBytes,
	}
	server.httpServer = httpServer

	return server
}

func (this *HttpServer) Start() {
	doStart := func() {
		err := this.httpServer.ListenAndServe()
		if err != nil {
			panic(fmt.Errorf("JsonRpcServer ListenAndServe error:%s", err))
		}
	}
	go doStart()
}

func (this *HttpServer) GetCpuStat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := MStat.GetCpuStat()
	this.returnResponse(data, w)
}

func (this *HttpServer) GetMemStat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := MStat.GetMemStat()
	this.returnResponse(data, w)
}

func (this *HttpServer) GetDisStat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := MStat.GetDisStat()
	this.returnResponse(data, w)
}

func (this *HttpServer) GetNetStat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := MStat.GetNetStat()
	this.returnResponse(data, w)
}

func (this *HttpServer)GetMachineStat(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	data := MStat.GetMachineStat()
	this.returnResponse(data, w)
}

func (this *HttpServer) returnResponse(rsp interface{}, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(rsp)
	w.Write(data)
}

type NotFound struct{}

func (this *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log4.Debug("Cannot handle:%s\n", r.URL.String())
	w.WriteHeader(http.StatusNotFound)
}
