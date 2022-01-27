package memory

import (
	"container/list"
	"sanguosha.com/sgs_herox/gameshared/manager/session"
	"sync"
	"time"
)

type SessionStore struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
}

func (p *SessionStore) Set(key, value interface{}) error {
	pder.SessionUpdate(p.sid)
	p.value[key] = value
	return nil
}

func (p *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(p.sid)
	if v, ok := p.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (p *SessionStore) Delete(key interface{}) error {
	pder.SessionUpdate(p.sid)
	delete(p.value, key)
	return nil
}

func (p *SessionStore) SessionID() string {
	return p.sid
}

type Provider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element
	list     *list.List
}

func (p *Provider) SessionInit(sid string) (session.Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	sess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := p.list.PushBack(sess)
	p.sessions[sid] = element
	return sess, nil
}

func (p *Provider) SessionRead(sid string) (session.Session, error) {
	if e, ok := p.sessions[sid]; ok {
		return e.Value.(*SessionStore), nil
	} else {
		sess, err := p.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (p *Provider) SessionDestroy(sid string) error {
	if element, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		p.list.Remove(element)
		return nil
	}
	return nil
}

func (p *Provider) SessionGC(maxLifeTime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		e := p.list.Back()
		if e == nil {
			break
		}

		ss := e.Value.(*SessionStore)
		if (ss.timeAccessed.Unix() + maxLifeTime) < time.Now().Unix() {
			p.list.Remove(e)
			delete(p.sessions, ss.sid)
		} else {
			break
		}
	}
}

func (p *Provider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if e, ok := p.sessions[sid]; ok {
		e.Value.(*SessionStore).timeAccessed = time.Now()
		p.list.MoveToFront(e)
		return nil
	}
	return nil
}

var pder = &Provider{list: list.New()}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}
