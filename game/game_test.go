package game

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"log"
//	"net/http"
//	"runtime/pprof"
//	"sanguosha.com/baselib/appframe"
//	"sanguosha.com/baselib/ioservice"
//	"sanguosha.com/sgs_herox/game/core"
//	"sanguosha.com/sgs_herox/gameshared/config"
//	"sanguosha.com/sgs_herox/gameshared/config/conf"
//	"sanguosha.com/sgs_herox/gameshared/manager"
//	"sanguosha.com/sgs_herox/proto/cmsg"
//	gamedef "sanguosha.com/sgs_herox/proto/def"
//	"strconv"
//	"sync"
//	"testing"
//)
//
//func Test_Game(t *testing.T) {
//	cfg, err := config.ParseConfigFile("../bin/app.yaml")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	netConfigFile := "../bin/netconfig.json"
//	app, err := appframe.NewApplication(netConfigFile, "game1")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	confManager := manager.NewConfManager(app, cfg.GameCfgPath, cfg.Develop, func(conf *conf.GameConfig) {
//		gCfg = conf
//	})
//	confManager.LoadConf()
//
//	roomSetting := &gamedef.RoomSetting{MaxPlayer: 2}
//	gameId := "abcd"
//
//	worker := NewLabelWorker()
//
//	game := newGame(worker, roomSetting, gameId, 100, 100, gCfg)
//	g := game.base
//
//	for i := 0; i < 2; i++ {
//		sessionId := appframe.SessionID{
//			SvrID: 0,
//			ID:    0,
//		}
//		session := app.GetSession(sessionId)
//
//		userId := uint64(i)
//		seatId := int32(i)
//		brief := &gamedef.UserBrief{
//			UserID: userId,
//		}
//		u := newUser(userId, seatId, brief, session)
//		g.OnUserJoin(u, seatId)
//		//playerMgr.add(player)
//	}
//
//	g.prepareStart()
//
//	http.HandleFunc("/action", ActionHandler(g))
//	http.HandleFunc("/ready", ReadyHandler(g))
//	err = http.ListenAndServe("0.0.0.0:80", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func ActionHandler(g *GameBase) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var wg sync.WaitGroup
//		wg.Add(1)
//		Command(g, w, r.FormValue("action"), &wg)
//		wg.Wait()
//	}
//}
//
//func ReadyHandler(g *GameBase) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var wg sync.WaitGroup
//		wg.Add(4)
//		Command(g, w, "place", &wg)
//		Command(g, w, "move", &wg)
//		Command(g, w, "cancel", &wg)
//		Command(g, w, "place", &wg)
//		wg.Wait()
//	}
//}
//
//func Command(g *GameBase, w http.ResponseWriter, action string, wg *sync.WaitGroup) {
//	g.DoNow(func() {
//		defer wg.Done()
//		//defer g.Print()
//		p := g.GetCurrentPlayer()
//		switch action {
//		case "place":
//			var cardId, row, col int32
//			if cardId <= 0 {
//				if v, ok := p.HandCard.Front(); !ok {
//					w.Write([]byte("没有可放置的牌"))
//					return
//				} else {
//					cardId = v
//				}
//
//				cells := g.board.SeatPlaceEmptyCell(p.GetSeatID())
//				if len(cells) == 0 {
//					w.Write([]byte("没有可放置的位置"))
//					return
//				}
//
//				row = cells[0].Row
//				col = cells[0].Col
//			}
//
//			msg := &cmsg.CReqAct{
//				ActType: gamedef.ActType_PlaceCard,
//				CardId:  cardId,
//				TargetPos: &gamedef.Position{
//					Row: row,
//					Col: col,
//				},
//			}
//			onReqAct(p, msg)
//
//			data, _ := json.Marshal(msg)
//			w.Write(data)
//		case "move":
//			cells := g.board.GetCellBySeat(p.GetSeatID())
//
//			var fromCell *Cell
//			var toPos Position
//			for _, v := range cells {
//				if to, ok := g.board.CellCanMovePos(v); ok {
//					fromCell = v
//					toPos = to
//					break
//				}
//			}
//
//			msg := &cmsg.CReqAct{
//				ActType: gamedef.ActType_MoveCard,
//				//CardId:    fromCell.cardCfg.CardID,
//				TargetPos: toPos.ToDef(),
//			}
//
//			onReqAct(p, msg)
//			data, _ := json.Marshal(msg)
//			w.Write(data)
//		case "cancel":
//			onReqCancel(p, &cmsg.CReqCancelCurOpt{})
//		case "attack":
//			cells := g.board.GetCellBySeat(p.GetSeatID())
//
//			var fromCell *Cell
//			var toPos Position
//			for _, v := range cells {
//				if to, ok := g.board.CellCanAttackPos(v); ok {
//					fromCell = v
//					toPos = to
//					break
//				}
//			}
//
//			msg := &cmsg.CReqAct{
//				ActType: gamedef.ActType_MoveCard,
//				//CardId:    fromCell.cardCfg.CardID,
//				TargetPos: toPos.ToDef(),
//			}
//			onReqAct(p, msg)
//			data, _ := json.Marshal(msg)
//			w.Write(data)
//		case "print":
//			//g.Print()
//		default:
//			w.Write([]byte("action not support"))
//		}
//	})
//	//g.Print()
//}
//
//type LabelWorker struct {
//	worker core.Worker
//}
//
//func NewLabelWorker() core.Worker {
//	worker := ioservice.NewIOService(fmt.Sprintf("app_%s_main", "game"), 102400)
//	worker.Init()
//	worker.Run()
//	return &LabelWorker{worker: worker}
//}
//
//func (p *LabelWorker) Post(f func()) {
//	p.worker.Post(func() {
//		pprof.Do(context.Background(), pprof.Labels("game", strconv.Itoa(0)), func(_ context.Context) {
//			f()
//		})
//	})
//}
