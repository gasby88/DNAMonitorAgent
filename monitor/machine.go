package monitor

import (
	"DNAMonitorAgent/conf"
	log4 "github.com/alecthomas/log4go"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"sync"
	"time"
)

type CpuStat struct {
	Idle        float64
	UsedPercent float64
}

type MemStat struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type DisStat struct {
	Total       uint64
	Free        uint64
	Used        uint64
	UsedPercent float64
}

type NetStat struct {
	Name      string
	recvBytes uint64
	RecvRate  float64
	sendBytes uint64
	SendRate  float64
}

type HostStat struct {
	Hostname string
	OS       string
	Platform string
}

type ProcStat struct {
	ProcName   string
	CreateTime int64
	RunTime    int64
}

type MachineStatMgr struct {
	cpu      *CpuStat
	mem      *MemStat
	dis      *DisStat
	net      map[string]*NetStat
	host     *HostStat
	proc     *ProcStat
	interval int
	exitCh   chan interface{}
	lock     sync.RWMutex
}

type MachineStat struct {
	Cpu  *CpuStat
	Mem  *MemStat
	Dis  *DisStat
	Net  []*NetStat
	Host *HostStat
	Proc *ProcStat
}

var MStat *MachineStatMgr

func NewMachineStat(interval int) *MachineStatMgr {
	if interval == 0 {
		interval = 1
	}
	return &MachineStatMgr{
		interval: interval,
		net:      make(map[string]*NetStat, 0),
		exitCh:   make(chan interface{}, 0),
	}
}

func (this *MachineStatMgr) Start() {
	log4.Info("MachineStat Start.")
	go func() {
		statTicker := time.NewTicker(time.Duration(this.interval) * time.Second)
		for {
			select {
			case <-this.exitCh:
				return
			case <-statTicker.C:
				go this.UpdMachineStat()
			}
		}
	}()
}

func (this *MachineStatMgr) Close() {
	close(this.exitCh)
}

func (this *MachineStatMgr) UpdMachineStat() {
	this.UpdCpuStat()
	this.UpdDisStat()
	this.UpdMemStat()
	this.UpdNetStat()
	this.UpdHostStat()
	this.UpdProcStat()
}

func (this *MachineStatMgr) UpdProcStat() {
	pids, err := process.Pids()
	if err != nil {
		log4.Error("UpdProcStat process.Pids error:%s", err)
		return
	}

	createTimeInit := int64(1<<62)
	olderCreateTime := int64(1<<62)
	for _, pid := range pids {
		proc, err := process.NewProcess(pid)
		if err != nil {
			log4.Error("UpdProcStat process.NewProcess Pid:%v error:%s", pid, err)
			continue
		}
		name, err := proc.Name()
		if err != nil {
			log4.Error("UpdProcStat proc.Name Pid:%v error:%s", pid, err)
			continue
		}
		if name != conf.GCfg.ProcName {
			continue
		}
		createTime, err := proc.CreateTime()
		if err != nil {
			log4.Error("UpdProcStat ProcName:%v proc.CreateTime error:%s", name, err)
			continue
		}
		if createTime < olderCreateTime {
			olderCreateTime = createTime
		}
	}

	if olderCreateTime == createTimeInit {
		olderCreateTime = 0
	}
	olderCreateTime = olderCreateTime / 1000
	runTime := int64(0)
	if olderCreateTime == 0 {
		runTime = 0
	}else{
		runTime = time.Now().Unix() - olderCreateTime
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.proc = &ProcStat{
		ProcName:   conf.GCfg.ProcName,
		CreateTime: olderCreateTime,
		RunTime:    runTime,
	}
}

func (this *MachineStatMgr) GetProcStat() *ProcStat {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.proc
}

func (this *MachineStatMgr) UpdHostStat() {
	stat, err := host.Info()
	if err != nil {
		log4.Error("UpdHostStat host.Info error:%s", err)
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.host = &HostStat{
		Hostname: stat.Hostname,
		OS:       stat.OS,
		Platform: stat.Platform,
	}
}

func (this *MachineStatMgr) GetHostStat() *HostStat {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.host
}

func (this *MachineStatMgr) UpdCpuStat() {
	stat, err := cpu.Percent(time.Duration(this.interval)*time.Second, false)
	if err != nil {
		log4.Error("SetCpuStat cpu.Percent error:%s", err)
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.cpu = &CpuStat{
		UsedPercent: stat[0],
		Idle:        100 - stat[0],
	}
}

func (this *MachineStatMgr) GetCpuStat() *CpuStat {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.cpu
}

func (this *MachineStatMgr) UpdMemStat() {
	stat, err := mem.VirtualMemory()
	if err != nil {
		log4.Error("SetMemStat mem.VirtualMemory error:%s", err)
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.mem = &MemStat{
		Total:       stat.Total,
		Free:        stat.Available,
		Used:        stat.Used,
		UsedPercent: stat.UsedPercent,
	}
}

func (this *MachineStatMgr) GetMemStat() *MemStat {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.mem
}

func (this *MachineStatMgr) UpdDisStat() {
	stat, err := disk.Usage("/")
	if err != nil {
		log4.Error("SetDisStat disk.Usage error:%s", err)
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.dis = &DisStat{
		Total:       stat.Total,
		Free:        stat.Free,
		Used:        stat.Used,
		UsedPercent: stat.UsedPercent,
	}
}

func (this *MachineStatMgr) GetDisStat() *DisStat {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.dis
}

func (this *MachineStatMgr) UpdNetStat() {
	stats, err := net.IOCounters(true)
	if err != nil {
		log4.Error("SetNetStat net.IOCounters(true) error:%s", err)
		return
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, stat := range stats {
		name := stat.Name
		oldStat, ok := this.net[name]
		recvRate := float64(0)
		sendRate := float64(0)
		if ok {
			recvRate = float64(stat.BytesRecv-oldStat.recvBytes) / float64(this.interval)
			sendRate = float64(stat.BytesSent-oldStat.sendBytes) / float64(this.interval)
		}
		this.net[stat.Name] = &NetStat{
			Name:      stat.Name,
			RecvRate:  recvRate,
			SendRate:  sendRate,
			recvBytes: stat.BytesRecv,
			sendBytes: stat.BytesSent,
		}
	}
}

func (this *MachineStatMgr) GetNetStat() []*NetStat {
	this.lock.RLock()
	defer this.lock.RUnlock()
	stats := make([]*NetStat, 0, len(this.net))
	for _, stat := range this.net {
		stats = append(stats, stat)
	}
	return stats
}

func (this *MachineStatMgr) GetMachineStat() *MachineStat {
	return &MachineStat{
		Cpu:  this.GetCpuStat(),
		Mem:  this.GetMemStat(),
		Dis:  this.GetDisStat(),
		Net:  this.GetNetStat(),
		Host: this.GetHostStat(),
		Proc: this.GetProcStat(),
	}
}