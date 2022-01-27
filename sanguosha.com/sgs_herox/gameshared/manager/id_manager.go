package manager

import (
	"github.com/pkg/errors"
	"math"
	"sanguosha.com/sgs_herox/gameutil"
	"sync"
	"time"
)

const sequenceMax = 1e5

type IDManager struct {
	serverID  uint32
	lastSeqID uint64
	timestamp int64
	mutex     sync.Mutex
}

func (p *IDManager) Init(serverID uint32) error {
	maxServerID := uint32(math.Pow(2, 12))
	if serverID > maxServerID {
		return errors.New("serverID too max")
	}
	p.serverID = serverID
	return nil
}

func (p *IDManager) GeneratePKID() uint64 {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	now := gameutil.GetCurrentTimestamp()
	if p.timestamp >= now {
		p.lastSeqID++
		if p.lastSeqID >= sequenceMax {
			p.timestamp++
			p.lastSeqID = 0
			time.Sleep(time.Second)
			//p.lastSeqID = 0
			//for now <= p.timestamp {
			//	time.Sleep(time.Second)
			//	now = gameutil.GetCurrentTimestamp()
			//}
		}
	} else {
		p.lastSeqID = 0
		p.timestamp = now
	}

	return uint64(p.timestamp)*1e9 + uint64(p.serverID)*sequenceMax + p.lastSeqID
}
