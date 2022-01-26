package serverstateservice

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/balance"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/smsg"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type NodeInfo struct {
	ID uint32

	online bool
	load   int64
	delay  int32
	limit  int32
}

func (node *NodeInfo) CheckLoadable(cmp int64) (bool, int64) {
	if !node.online {
		return false, 0
	}
	if node.limit > 0 {
		if node.load > int64(node.limit) {
			return false, node.load
		}
	}
	if node.load > cmp {
		return false, 0
	}
	return true, node.load
}

type serverStateService struct {
	appInstance *appframe.Application
	lock        sync.RWMutex

	svrWatched      map[appframe.ServerType]bool
	balances        []*balance.StateBalance
	nodes           map[uint32]*NodeInfo
	serverInfos     *gameshared.ServerInfos
	loadSettingAi   []int32
	loadSettingGame []int32
	unloadAIs       map[uint32]bool
	unloadGames     map[uint32]bool
	unloadSvrs      map[uint32]bool

	onServerInfoUpdate []func(serverInfos *gameshared.ServerInfos)
}

func newServerStateService(app *appframe.Application) *serverStateService {
	return &serverStateService{
		appInstance: app,
		serverInfos: gameshared.NewServerInfos(),
		nodes:       make(map[uint32]*NodeInfo),
		svrWatched:  make(map[appframe.ServerType]bool),
		unloadSvrs:  make(map[uint32]bool),
		unloadAIs:   make(map[uint32]bool),
		unloadGames: make(map[uint32]bool),
	}
}

func (sss *serverStateService) GetServerInfos() *gameshared.ServerInfos {
	return sss.serverInfos
}

func (sss *serverStateService) OnServerInfoUpdate(callback func(serverInfos *gameshared.ServerInfos)) {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	sss.onServerInfoUpdate = append(sss.onServerInfoUpdate, callback)
}

func (sss *serverStateService) NewStateBalance(svrType appframe.ServerType) *balance.StateBalance {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	b := balance.NewStateBalance(sss, sss.appInstance, svrType)
	sss.balances = append(sss.balances, b)
	if _, ok := sss.svrWatched[svrType]; !ok {
		sss.svrWatched[svrType] = true
		sss.appInstance.ListenServerEvent(svrType, sss.onServerEvent)
	}
	return b
}

func (sss *serverStateService) WatchServer(svrType appframe.ServerType) {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	if _, ok := sss.svrWatched[svrType]; !ok {
		sss.svrWatched[svrType] = true
		sss.appInstance.ListenServerEvent(svrType, sss.onServerEvent)
	}
}

func (sss *serverStateService) SetLoad(svrID uint32, load int64) {
	if load < 0 {
		load = 0
	}
	sss.lock.Lock()
	defer sss.lock.Unlock()

	node := sss._getNode(svrID)
	node.load = load
}
func (sss *serverStateService) ModLoad(svrID uint32, load int64) {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	node := sss._getNode(svrID)
	node.load += load
	if node.load < 0 {
		node.load = 0
	}
}

func (sss *serverStateService) SetServerDelay(svrID uint32, svrType appframe.ServerType, delay int64) {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	node := sss._getNode(svrID)
	if delay > 60000 {
		node.delay = 60000
	} else {
		node.delay = int32(delay)
	}
	switch svrType {
	case sgs_herox.SvrTypeGame:
		node.limit = sss._getGameServerLoadLimit(svrID, node.delay)
	case sgs_herox.SvrTypeAI:
		node.limit = sss._getAiServerLoadLimit(svrID, node.delay)
	}
}

func (sss *serverStateService) GetLoadInfo() string {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	sortKV := gameutil.SortUInt64KV{IsLess: true}
	for k, node := range sss.nodes {
		sortKV.List = append(sortKV.List, gameutil.UInt64KV{uint64(k), uint64(node.ID)})
	}
	sort.Sort(sortKV)

	loadInfo := []string{}
	for _, kv := range sortKV.List {
		k := uint32(kv.Value)
		node := sss.nodes[uint32(kv.Value)]
		if node.load == 0 {
			continue
		}
		switch appframe.ServerType(node.ID / 100) {
		case sgs_herox.SvrTypeAI:
			loadInfo = append(loadInfo, fmt.Sprintf("(AI_%d:%d:%d)", k, node.load, node.delay))
		case sgs_herox.SvrTypeGame:
			loadInfo = append(loadInfo, fmt.Sprintf("(GAME_%d:%d:%d)", k, node.load, node.delay))
		default:
			loadInfo = append(loadInfo, fmt.Sprintf("(SVR_%d:%d:%d)", k, node.load, node.delay))
		}
	}
	if len(loadInfo) != 0 {
		return strings.Join(loadInfo, " ")
	}
	//loadInfo := make([]string, 0, len(GameManagerInstance.aiNode)+len(GameManagerInstance.gameNode))
	//for k, node := range GameManagerInstance.aiNode {
	//	if node.Load != 0 {
	//		loadInfo = append(loadInfo, fmt.Sprintf("(AI_%d:%d:%d)", k, node.Load, node.Delay))
	//	}
	//}
	//for k, node := range GameManagerInstance.gameNode {
	//	if node.Load != 0 {
	//		loadInfo = append(loadInfo, fmt.Sprintf("(GAME_%d:%d:%d)", k, node.Load, node.Delay))
	//	}
	//}
	return ""
}

func (sss *serverStateService) ResetServerInfo(m map[string]string) {
	logrus.Info("ServerInfos Reset ", sss.appInstance.ID(), " ", m)
	sss.serverInfos.Reset(m)
	sss.ReloadServerInfo()
}

func (sss *serverStateService) onSyncServerInfo(sender appframe.Server, req *smsg.SyncServerInfo) {
	sss.serverInfos.SetAll(req.Keys, req.Values)
	sss.ReloadServerInfo()
}
func (sss *serverStateService) ReloadServerInfo() {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	sss.loadSettingAi = []int32{}
	sss.loadSettingGame = []int32{}
	sss.unloadAIs = make(map[uint32]bool)
	sss.unloadGames = make(map[uint32]bool)
	sss.unloadSvrs = make(map[uint32]bool)

	serverInfos := sss.serverInfos
	{
		value := serverInfos.Get("LoadAI")
		list := strings.Split(value, ",")
		if len(list) > 2 {
			for _, tmp := range list {
				n, e := strconv.Atoi(tmp)
				if e != nil {
					break
				}
				sss.loadSettingAi = append(sss.loadSettingAi, int32(n))
			}
		}
	}
	{
		value := serverInfos.Get("LoadGame")
		list := strings.Split(value, ",")
		if len(list) > 2 {
			for _, tmp := range list {
				n, e := strconv.Atoi(tmp)
				if e != nil {
					break
				}
				sss.loadSettingGame = append(sss.loadSettingGame, int32(n))
			}
		}
	}
	{
		value := serverInfos.Get("UnloadAI")
		ids, _ := gameutil.ParseUInt32s(value, ";")
		for _, id := range ids {
			sss.unloadAIs[id] = true
		}
	}
	{
		value := serverInfos.Get("UnloadGame")
		ids, _ := gameutil.ParseUInt32s(value, ";")
		for _, id := range ids {
			sss.unloadGames[id] = true
		}
	}
	{
		value := serverInfos.Get("UnloadSvr")
		ids, _ := gameutil.ParseUInt32s(value, ";")
		for _, id := range ids {
			sss.unloadSvrs[id] = true
		}
	}

	//("ServerInfos Update %d LoadAI %v LoadGame %v UnloadAI %v UnloadGame %v UnloadSvr %v", sss.appInstance.ID(), sss.loadSettingAi, sss.loadSettingGame, sss.unloadAIs, sss.unloadGames, sss.unloadSvrs)

	for _, callback := range sss.onServerInfoUpdate {
		c := callback
		sss.appInstance.Post(func() {
			c(sss.serverInfos)
		})
	}
}
func (sss *serverStateService) _getGameServerLoadLimit(svrID uint32, delay int32) int32 {
	load := sss.loadSettingGame
	limit := int32(math.MaxInt32)
	n := len(load) / 2
	for i := 0; i < n; i++ {
		if delay >= load[2*i] {
			limit = load[2*i+1]
		}
	}
	return limit
}
func (sss *serverStateService) _getAiServerLoadLimit(svrID uint32, delay int32) int32 {
	load := sss.loadSettingAi
	limit := int32(math.MaxInt32)
	n := len(load) / 2
	for i := 0; i < n; i++ {
		if delay >= load[2*i] {
			limit = load[2*i+1]
		}
	}
	return limit
}

func NoSvrError(svrType appframe.ServerType) error {
	switch svrType {
	case sgs_herox.SvrTypeGame:
		return balance.ErrNoAvailableGame
	case sgs_herox.SvrTypeAI:
		return balance.ErrNoAvailableAI
	default:
		return balance.ErrNoAvailableServer
	}
}

func (sss *serverStateService) GetLoadableServer(svrType appframe.ServerType) (appframe.Server, error) {
	svrids := sss.appInstance.GetAvailableServerIDs(svrType)
	if len(svrids) == 0 {
		return nil, NoSvrError(svrType)
	}

	sss.lock.RLock()
	defer sss.lock.RUnlock()

	var loadIds = []uint32{}
	var loadMin = int64(math.MaxInt64)
	var full = false
	for _, id := range svrids {
		if sss._isUnload(id, svrType) {
			continue
		}
		node := sss._getNode(id)
		ok, v := node.CheckLoadable(loadMin)
		if ok && v != int64(math.MaxInt64) {
			if loadMin == v {
				loadIds = append(loadIds, id)
			} else {
				loadMin = v
				loadIds = []uint32{id}
			}
		} else if v != 0 {
			full = true
		}
	}
	nOK := len(loadIds)
	if nOK == 0 {
		if full {
			return nil, balance.ErrServiceLoadFull
		}
		return nil, NoSvrError(svrType)
	}
	idx := 0
	if nOK != 1 {
		idx = gameutil.Rand() % nOK
	}
	return sss.appInstance.GetServer(loadIds[idx]), nil
}

func (sss *serverStateService) _isUnload(svrID uint32, svrType appframe.ServerType) bool {
	switch svrType {
	case sgs_herox.SvrTypeAI:
		if _, ok := sss.unloadAIs[svrID]; ok {
			return true
		}
	case sgs_herox.SvrTypeGame:
		if _, ok := sss.unloadGames[svrID]; ok {
			return true
		}
	default:
		if _, ok := sss.unloadSvrs[svrID]; ok {
			return true
		}
	}
	return false
}

func (sss *serverStateService) onServerEvent(svrid uint32, event netcluster.SvrEvent) {
	sss.lock.Lock()
	defer sss.lock.Unlock()

	switch event {
	case netcluster.SvrEventStart, netcluster.SvrEventReconnect:
		node := sss._getNode(svrid)
		node.online = true
		node.load = 0
		node.delay = 0
		node.limit = 0
	case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
		node := sss._getNode(svrid)
		node.online = false
	}
}
func (sss *serverStateService) _getNode(svrid uint32) *NodeInfo {
	sl, ok := sss.nodes[svrid]
	if !ok {
		sl = &NodeInfo{ID: svrid}
		sss.nodes[svrid] = sl
	}
	return sl
}
