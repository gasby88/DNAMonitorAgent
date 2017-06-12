package monitor

import "DNAMonitorAgent/common"

func init()  {
	MonitorMgr.RegHandler(&CpuStatHandler{})
	MonitorMgr.RegHandler(&MemStatHandler{})
	MonitorMgr.RegHandler(&DisStatHandler{})
	MonitorMgr.RegHandler(&NetStatHandler{})
	MonitorMgr.RegHandler(&HostStatHandler{})
	MonitorMgr.RegHandler(&ProcStatHandler{})
	MonitorMgr.RegHandler(&MachineStatHandler{})
}

type CpuStatHandler struct{}

func (this *CpuStatHandler) GetName() string {
	return "cpu"
}

func (this *CpuStatHandler) Handle(req *common.DNAMonitorRequest) (interface{}, int) {
	return MStat.GetCpuStat(), common.Err_OK
}

type MemStatHandler struct{}

func (this *MemStatHandler) GetName() string {
	return "mem"
}

func (this *MemStatHandler) Handle(req *common.DNAMonitorRequest) (interface{}, int) {
	return MStat.GetMemStat(), common.Err_OK
}

type DisStatHandler struct{}

func (this *DisStatHandler) GetName() string {
	return "dis"
}

func (this *DisStatHandler) Handle(req *common.DNAMonitorRequest) (interface{}, int) {
	return MStat.GetDisStat(), common.Err_OK
}

type NetStatHandler struct{}

func (this *NetStatHandler) GetName() string {
	return "net"
}

func (this *NetStatHandler) Handle(req *common.DNAMonitorRequest) (interface{}, int) {
	return MStat.GetNetStat(), common.Err_OK
}

type HostStatHandler struct {}

func (this *HostStatHandler) GetName() string {
	return "host"
}

func (this *HostStatHandler) Handle(req *common.DNAMonitorRequest)(interface{}, int){
	return MStat.GetHostStat(), common.Err_OK
}

type ProcStatHandler struct {}

func (this *ProcStatHandler) GetName()string {
	return "proc"
}

func (this *ProcStatHandler) Handle(req *common.DNAMonitorRequest)(interface{}, int){
	return MStat.GetProcStat(), common.Err_OK
}

type MachineStatHandler struct{}

func (this *MachineStatHandler) GetName() string {
	return "machine"
}

func (this *MachineStatHandler) Handle(req *common.DNAMonitorRequest) (interface{}, int) {
	return MStat.GetMachineStat(), common.Err_OK
}
