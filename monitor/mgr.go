package monitor

import (
	"DNAMonitorAgent/common"
	log4 "github.com/alecthomas/log4go"
)

var MonitorMgr = NewDNAMonitorMgr()

type MonitorHandler interface {
	GetName() string
	Handle(req *common.DNAMonitorRequest) (result interface{}, errorCode int)
}

type DNAMonitorMgr struct {
	handlers map[string]MonitorHandler
}

func NewDNAMonitorMgr() *DNAMonitorMgr {
	return &DNAMonitorMgr{
		handlers: make(map[string]MonitorHandler, 0),
	}
}

func (this *DNAMonitorMgr) GetHandler(name string) MonitorHandler {
	handler, ok := this.handlers[name]
	if ok {
		return handler
	}
	return nil
}

func (this *DNAMonitorMgr) setHandler(handler MonitorHandler) {
	this.handlers[handler.GetName()] = handler
}

func (this *DNAMonitorMgr) RegHandler(handler MonitorHandler) {
	h := this.GetHandler(handler.GetName())
	if h != nil {
		log4.Error("DNAMonitorMgr RegHandler:%s error, Handler has already existed", handler.GetName())
		return
	}
	this.setHandler(handler)
}

func (this *DNAMonitorMgr) Handle(req *common.DNAMonitorRequest) (interface{}, int) {
	handler := this.GetHandler(req.Method)
	if handler == nil {
		return nil, common.Err_Method_Not_EXIST
	}
	return handler.Handle(req)
}
