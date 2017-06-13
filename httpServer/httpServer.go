package httpServer

import (
	"DNAMonitorAgent/common"
	. "DNAMonitorAgent/monitor"
	"fmt"
	log4 "github.com/alecthomas/log4go"
	"io/ioutil"
	"net/http"
)

type HttpServer struct {
	port string
	path string
}

func NewHttpServer(port, path string) *HttpServer {
	return &HttpServer{
		port: port,
		path: path,
	}
}

func (this *HttpServer) Start() {
	http.HandleFunc(this.path, this.OnHandle)
	doStart := func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", this.port), nil)
		if err != nil {
			panic(fmt.Errorf("JsonRpcServer ListenAndServe error:%s", err))
		}
	}
	go doStart()
}

func (this *HttpServer) OnHandle(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log4.Error("DNAMonitorAgent Handle read body error:%s", err)
		this.writeResponse(common.NewDNAMonitorResponse("", "", nil, common.Err_Unknow), w)
		return
	}
	req, err := common.NewDNAMonitorRequest(data)
	if err != nil {
		log4.Error("HttpServer NewDNAMonitorRequest from:%s error:%s", data, err)
		this.writeResponse(common.NewDNAMonitorResponse("", "", nil, common.Err_Unknow), w)
		return
	}
	res, errCode := MonitorMgr.Handle(req)
	this.writeResponse(common.NewDNAMonitorResponse(req.Qid, req.Method, res, errCode), w)
}

func (this *HttpServer) writeResponse(rsp *common.DNAMonitorResponse, w http.ResponseWriter) {
	rspData, _ := rsp.Marshal()
	w.WriteHeader(http.StatusOK)
	w.Write(rspData)
}

type NotFound struct{}

func (this *NotFound) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log4.Debug("Cannot handle:%s\n", r.URL.String())
	w.WriteHeader(http.StatusNotFound)
}
