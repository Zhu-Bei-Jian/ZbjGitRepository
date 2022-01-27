package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sanguosha.com/baselib/util"
	"sync"
	"time"
)

var provides = make(map[string]Provider)

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session:Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session:Register called twice for provider " + name)
	}
	provides[name] = provider
}

type Manager struct {
	cookieName  string
	lock        sync.Mutex
	provider    Provider
	maxLifeTime int64
}

type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}

func NewManager(provideName, cookieName string, maxLiftTime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session:unknow provide %q", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLiftTime}, nil
}

func (p *Manager) RunGC() {
	util.SafeGo(func() { p.GC() })
}

func (p *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (p *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	p.lock.Lock()
	defer p.lock.Unlock()
	cookie, err := r.Cookie(p.cookieName)
	if err != nil || cookie.Value == "" {
		sid := p.sessionId()
		session, _ = p.provider.SessionInit(sid)
		cookie := http.Cookie{Name: p.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: 86400}
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = p.provider.SessionRead(sid)
	}
	return
}

func (p *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(p.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: p.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (p *Manager) GC() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.provider.SessionGC(p.maxLifeTime)
	time.AfterFunc(time.Duration(p.maxLifeTime)*time.Second, func() {
		p.GC()
	})
}
