package gameutil

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/pprof"
	"sanguosha.com/baselib/log"
	"strings"
	"sync"
	"time"

	_ "net/http/pprof"
)

type BathAuth struct {
	User     string
	Password string
	Key      string
}

func Pprof(addr string, auth *BathAuth) error {
	if auth == nil {
		auth = &BathAuth{User: "a", Password: "@qwe4rrt6fah5", Key: "@qwe4rrt6fah6"}
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	router.HandleFunc("/debug/pprof/", pprof.Index).Methods("GET", "POST", "PUT")
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline).Methods("GET", "POST", "PUT")
	router.HandleFunc("/debug/pprof/profile", pprof.Profile).Methods("GET", "POST", "PUT")
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol).Methods("GET", "POST", "PUT")
	router.HandleFunc("/debug/pprof/trace", pprof.Trace).Methods("GET", "POST", "PUT")

	router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	serverMux := http.NewServeMux()
	serverMux.Handle("/debug/", NewHttpBasicAuth(auth.User, auth.Password, auth.Key, router))
	return http.Serve(listener, serverMux)
	//fmt.Println(http.ListenAndServe(addr, nil))
}

var pprofAuth map[string]int64
var pprofLock sync.Mutex

func init() {
	if pprofAuth == nil {
		pprofAuth = make(map[string]int64)
	}
	rand.Seed(time.Now().UnixNano())
}

func AuthInit(w http.ResponseWriter) {
	pprofLock.Lock()
	defer pprofLock.Unlock()
	var s string
	for {
		s = RandString(8)
		if s == "" {
			continue
		}
		if _, ok := pprofAuth[s]; ok {
			continue
		}
		pprofAuth[s] = time.Now().Unix()
		http.SetCookie(w, &http.Cookie{Name: "Auth", Value: s, Path: "/", MaxAge: 86400})
		//log.Debug("AuthCheck init ", s)
		return
	}
}

func AuthCheck(r *http.Request) bool {
	cookie, err := r.Cookie("Auth")
	if err != nil {
		return false
	}
	pprofLock.Lock()
	defer pprofLock.Unlock()
	timeTag, ok := pprofAuth[cookie.Value]
	if !ok {
		return false
	}

	if time.Now().Unix() > (timeTag + 86400) {
		return false
	}
	pprofAuth[cookie.Value] = time.Now().UnixNano()

	//log.Debug("AuthCheck pass ", cookie.Value)
	return true
}

type httpBasicAuth struct {
	user     string
	password string
	key      string
	handler  http.Handler
}

func NewHttpBasicAuth(user string, password string, key string, handler http.Handler) *httpBasicAuth {
	if user != "" && password != "" {
		//log.Debug("require authentication")
	}
	return &httpBasicAuth{user: user, password: password, key: key, handler: handler}
}

func (h *httpBasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.user == "" || h.password == "" {
		log.Debug("no auth required")
		h.handler.ServeHTTP(w, r)
		return
	}
	//log.Debug("auth required:", r.RequestURI)

	if AuthCheck(r) {
		h.handler.ServeHTTP(w, r)
		return
	}

	if len(h.key) != 0 {
		r.ParseForm()
		keyAuth := r.FormValue("KeyAuth")
		if keyAuth == h.key {
			AuthInit(w)
			h.handler.ServeHTTP(w, r)
			return
		}
	}

	username, password, ok := r.BasicAuth()
	if ok && username == h.user && len(h.user) != 0 {
		if strings.HasPrefix(h.password, "{SHA}") {
			log.Debug("auth with SHA")
			hash := sha1.New()
			io.WriteString(hash, password)
			if hex.EncodeToString(hash.Sum(nil)) == h.password[5:] {
				AuthInit(w)
				h.handler.ServeHTTP(w, r)
				return
			}
		} else if password == h.password {
			//log.Debug("Auth with normal password")
			AuthInit(w)
			h.handler.ServeHTTP(w, r)
			return
		}
	}
	w.Header().Set("WWW-Authenticate", "Basic realm=\"rnode\"")
	w.WriteHeader(401)
}
