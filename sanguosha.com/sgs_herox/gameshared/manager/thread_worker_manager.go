package manager

import (
	"fmt"
	"math/rand"
	"sanguosha.com/baselib/ioservice"
)

//多个工作线程
type ThreadWorkerMgr struct {
	workerCount uint32
	worker      []ioservice.IOService
}

func NewThreadWorkMgr(appID uint32, workerID int32, workerCnt uint32) *ThreadWorkerMgr {
	p := &ThreadWorkerMgr{}
	p.init(appID, workerID, workerCnt)
	return p
}

func (p *ThreadWorkerMgr) init(appID uint32, workerID int32, workerCnt uint32) error {
	p.workerCount = workerCnt
	p.worker = make([]ioservice.IOService, workerCnt)
	for i := 0; i < int(p.workerCount); i++ {
		p.worker[i] = ioservice.NewIOService(fmt.Sprintf("multi-worker-%d-%d-%d", appID, workerID, i), 10240)
		p.worker[i].Init()
	}
	return nil
}

func (p *ThreadWorkerMgr) Close() {
	for _, v := range p.worker {
		//Log("ThreadWorkerMgr.Close IOServiceIdx:%d is closing...", i)
		v.Fini()
		//Log("ThreadWorkerMgr.Close IOServiceIdex:%d is closed", i)
	}
}

func (p *ThreadWorkerMgr) Run() {
	for _, v := range p.worker {
		v.Run()
	}
}

func (p *ThreadWorkerMgr) WorkerLen() int32 {
	var sum int32
	for _, v := range p.worker {
		sum += int32(v.WorkerLen())
	}
	return sum
}

func (p *ThreadWorkerMgr) TryPost(f func()) bool {
	if len(p.worker) == 0 {
		return false
	}

	i := rand.Intn(len(p.worker))

	if p.worker[i].TryPost(f) {
		return true
	}
	return false
}

func (p *ThreadWorkerMgr) PostByID(id interface{}, f func()) bool {
	i := GetHashByID(id, p.workerCount)
	p.worker[i].Post(f)
	return true
}

func (p *ThreadWorkerMgr) Post(f func()) bool {
	i := rand.Intn(len(p.worker))
	p.worker[i].Post(f)
	return true
}
