package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type Player struct {
	user *User
	game *GameBase

	roleType   gameconf.RoleTyp
	seatId     int32
	lookerType gameconf.LookerTyp //旁观类型 上帝视角 玩家视角 双盲视角

	clientReady bool

	dead bool //是否已阵亡

	connectState gameconf.UserConnectState

	HandCard     CardPool //手牌
	HandCardPool CardPool

	actionPoint int32 //行动点
	attackCount int32 //本回合已经攻击的次数 包括攻击卡牌 和攻击基地
	Camp
}

func newPlayer(u *User, game *GameBase, seatId int32) *Player {
	p := &Player{
		user:        u,
		game:        game,
		seatId:      seatId,
		clientReady: false,
		Camp: Camp{
			hp: game.config.CampHP,
		},
	}
	if seatId == -1 {
		p.lookerType = gameconf.LookerTyp_LTBlind
	} else {
		p.lookerType = gameconf.LookerTyp_LTInvalid
	}

	return p
}

func (p *Player) GetGame() *GameBase {
	return p.game
}

func (p *Player) GetUser() *User {
	return p.user
}

func (p *Player) GetSeatID() int32 {
	return p.seatId
}

func (p *Player) SendMsg(msg proto.Message) {
	//data, _ := json.Marshal(msg)
	//logrus.WithFields(logrus.Fields{
	//	"seatId": p.GetSeatID(),
	//}).Debug(string(data))

	if !p.isClientReady() {
		return
	}
	p.user.SendMsg(msg)
}

func (p *Player) IsDead() bool {
	return p.dead
}

func (p *Player) IsAlive() bool {
	return !p.IsDead()
}

func (p *Player) SetDead() {
	p.dead = true
}

func (p *Player) Print() {
	fmt.Printf("seatId:%d hp:%d card:%v \n", p.GetSeatID(), p.GetHP(), p.HandCard.cards)
}

func (p *Player) disconnect() {
	p.connectState = gameconf.UserConnectState_USDisconnect
	p.setClientReady(false)
	p.user.OnDisconnect()
	//p.onPlayerVoiceGameDisconnect()
}

func (p *Player) reconnect(session appframe.Session) {
	p.connectState = gameconf.UserConnectState_USConnect
	p.setClientReady(false)
	p.user.OnReConnect(session)
}

func (p *Player) quit() {
	p.connectState = gameconf.UserConnectState_USQuit
	p.setClientReady(false)
	p.user.OnDisconnect()
}

func (p *Player) toSeatInfo(showRole bool) *gamedef.GameSeat {
	var roleType gameconf.RoleTyp
	if showRole || p.dead {
		roleType = p.roleType
	}
	return &gamedef.GameSeat{
		SeatId:        p.seatId,
		RoleType:      roleType,
		Dead:          p.dead,
		Hp:            p.GetHP(),
		Ap:            p.GetAP(),
		CardPoolCount: int32(p.HandCardPool.Count()),
	}
}

func (p *Player) setClientReady(v bool) {
	p.clientReady = v
}

func (p *Player) isClientReady() bool {
	return p.clientReady
}

func (p *Player) IsLooker() bool {
	return p.seatId == -1
}

func (p *Player) GetAP() int32 {
	return p.actionPoint
}

func (p *Player) SubAP(v int32) {
	p.actionPoint -= v
}

func (p *Player) SetAP(v int32) {
	p.actionPoint = v
}
func shuffleFisherYates(cards []int32) {
	n := len(cards)
	for i, _ := range cards {
		j := gameutil.Intn(n-i) + i
		cards[i], cards[j] = cards[j], cards[i]
	}
}

//func (p *Player) LoadCardPool() []int32 {
//	chInfo := entity.GetUserBySessionID(p.user.SessionID()).GetCharInfo()
//	cardIds := p.game.config.CardPool.GetPoolCardIds(1)
//	shuffleFisherYates(cardIds)
//	if len(cardIds) < int(p.game.config.CardPoolCount) {
//		logrus.WithFields(logrus.Fields{
//			"RealCardCount": len(cardIds),
//			"NeedCardCount": p.game.config.CardPoolCount,
//		}).Fatalf("总卡牌数量小于抽牌所需的卡牌数")
//		return nil
//	}
//
//	if chInfo.CardGroups == nil {
//		logrus.Warnf("ID为 %v 的玩家牌库数据为空，将为其生成一副随机卡组", p.user.userId)
//		cardIds = cardIds[:p.game.config.CardPoolCount]
//		logrus.Info(cardIds)
//		return cardIds
//	}
//	nowId := -1
//	for id, val := range chInfo.CardGroups {
//		if val.GroupId == chInfo.NowUseId {
//			nowId = id
//			break
//		}
//	}
//	if nowId == -1 {
//		logrus.Warnf("当前使用卡组ID为空，将生成一副随机卡组")
//		cardIds = cardIds[:p.game.config.CardPoolCount]
//		return cardIds
//	}
//	shuffleFisherYates(chInfo.CardGroups[nowId].Cards)
//	return chInfo.CardGroups[nowId].Cards
//}

func (p *Player) onEnterPhase(phase gamedef.GamePhase) {
	switch phase {
	case gamedef.GamePhase_PhaseBegin:
		p.attackCount = 0                               //该玩家本回合攻击次数 重置为0
		p.actionPoint = p.game.config.CommonActionPoint //重置 该玩家本回合的可用行动点
		if p.game.roundCount == 1 {
			p.actionPoint = p.game.config.FirtstRoundActionPoint
		}
		SyncPlayerState(p, cmsg.SSyncPlayerState_AP)
	}
}

type CardPool struct {
	cards []*gamedef.PoolCard
}

func (c *CardPool) Has(cardId int32) bool {
	for _, v := range c.cards {
		if v.CardId == cardId {
			return true
		}
	}
	return false
}

func (c *CardPool) IsEmpty() bool {
	return len(c.cards) == 0
}

func (c *CardPool) AddCard(card ...*gamedef.PoolCard) {
	c.cards = append(c.cards, card...)
}

func (c *CardPool) Count() int {
	return len(c.cards)
}

func (c *CardPool) RemoveFront() *gamedef.PoolCard {
	front := c.cards[0]
	c.cards = c.cards[1:]
	return front
}

func (c *CardPool) Front() (*gamedef.PoolCard, bool) {
	if len(c.cards) <= 0 {
		return nil, false
	}
	return c.cards[0], true
}

func (p *CardPool) Remove(cardId int32) bool {
	for i := 0; i < len(p.cards); {
		if p.cards[i].CardId == cardId {
			p.cards = append(p.cards[:i], p.cards[i+1:]...)
			return true
		} else {
			i++
		}
	}
	return false
}

func (p *CardPool) RemoveRandom() (*gamedef.PoolCard, bool) {
	count := p.Count()
	if count == 0 {
		return nil, false
	}

	randIndex := gameutil.Intn(count)
	card := p.cards[randIndex]
	p.cards = append(p.cards[:randIndex], p.cards[randIndex+1:]...)
	return card, true
}

func (p *CardPool) GetCard(cardId int32) (*gamedef.PoolCard, bool) {
	for _, v := range p.cards {
		if v.CardId == cardId {
			return v, true
		}
	}
	return nil, false
}

func (p *CardPool) Cards(show bool) []*gamedef.PoolCard {
	if show {
		return p.cards
	}

	var ret []*gamedef.PoolCard
	for _, v := range p.cards {
		ret = append(ret, &gamedef.PoolCard{
			CardId: v.CardId,
			HeroId: 0,
		})
	}

	return ret
}

func (p *CardPool) Clone() *CardPool {
	ret := new(CardPool)
	for _, v := range p.cards {
		ret.cards = append(ret.cards, &gamedef.PoolCard{
			CardId: v.CardId,
			HeroId: v.HeroId,
		})
	}
	return ret
}

type Camp struct {
	hp int32
	BuffManager
}

type BuffManager struct {
	buffs []*buff
}

type buff struct {
	buffCfg     *gameconf.BuffConfDefine
	useCount    int32 //已使用次数
	createRound int32 //buff创建回合
	ExpireType  gameconf.ExpireTyp
	ExpireV     int32
	buffCount   int32
}

func (p *Camp) SubHP(v int32) {
	p.hp -= v
}
func (p *Camp) AddHP(v int32) {
	p.hp += v
}

func (p *Camp) GetHP() int32 {
	return p.hp
}

func (p *Camp) SetHP(v int32) {
	p.hp = v
}

func (bm *BuffManager) Add(buffCfg *gameconf.BuffConfDefine, createRound int32, eType gameconf.ExpireTyp, eV int32) {

	tempBuff := &buff{
		buffCfg:     buffCfg,
		useCount:    0,
		createRound: createRound,
		ExpireType:  eType,
		ExpireV:     eV,
		buffCount:   1,
	}

	b, exist := bm.GetBuff(buffCfg.BuffID)
	if exist {
		if b.buffCfg.IsOverlap {
			b.buffCount++
		} else {
			bm.buffs = append(bm.buffs, tempBuff)
		}
	} else {
		bm.buffs = append(bm.buffs, tempBuff)
	}

}

func (bm *BuffManager) GetBuff(id int32) (*buff, bool) {
	for _, v := range bm.buffs {
		if v.buffCfg.BuffID == id {
			return v, true
		}
	}

	return nil, false
}

func (bm *BuffManager) AddBuffUseCount(buffId int32, v int32) {
	buff, exist := bm.GetBuff(buffId)
	if !exist {
		return
	}
	buff.useCount += v
}

func (bm *BuffManager) HasBuff(id int32) bool {
	_, ok := bm.GetBuff(id)
	return ok
}

func (bm *BuffManager) Remove(id int32) {
	for i := 0; i < len(bm.buffs); {
		if bm.buffs[i].buffCfg.BuffID == id {
			bm.buffs = append(bm.buffs[:i], bm.buffs[i+1:]...)
		} else {
			i++
		}
	}
}

func (bm *BuffManager) Clear() {
	bm.buffs = nil
}

func (bm *BuffManager) RemoveIfRoundExpire(roundCnt int32) []int32 {

	var remove []int32
	for i := 0; i < len(bm.buffs); {
		v := bm.buffs[i]
		if v.ExpireType == gameconf.ExpireTyp_ETRound {
			fmt.Printf("roundcnt = %v  endrounde= %v \n", roundCnt, v.ExpireV+v.createRound-1)
		}
		if v.ExpireType == gameconf.ExpireTyp_ETRound && roundCnt == v.ExpireV+v.createRound-1 {
			remove = append(remove, bm.buffs[i].buffCfg.BuffID)
			bm.buffs = append(bm.buffs[:i], bm.buffs[i+1:]...)
		} else {
			i++
		}
	}
	return remove
}

func (bm *BuffManager) ExpireBuffAtRound(roundCnt int32) ([]int32, bool) {
	var ret []int32
	hasRoundBuff := false
	for _, v := range bm.buffs {
		if v.ExpireType != gameconf.ExpireTyp_ETRound {
			continue
		}
		hasRoundBuff = true
		if roundCnt >= v.ExpireV-1 {
			ret = append(ret, v.buffCfg.BuffID)
		}
	}
	return ret, hasRoundBuff
}

func (bm *BuffManager) ExpireBuffInTimes() []int32 {
	var ret []int32
	for _, v := range bm.buffs {
		if v.ExpireType == gameconf.ExpireTyp_ETTimes && v.useCount >= v.ExpireV {
			ret = append(ret, v.buffCfg.BuffID)
		}
	}
	return ret
}

func (bm *BuffManager) ToDef() (ret []*gamedef.Buff) {

	for _, v := range bm.buffs {

		ret = append(ret, &gamedef.Buff{
			BuffId:      v.buffCfg.BuffID,
			ExpireType:  v.ExpireType,
			ExpireV:     v.ExpireV,
			CreateRound: v.createRound,
			UseCount:    v.useCount,
			BuffCount:   v.buffCount,
			IsOverlap:   v.buffCfg.IsOverlap,
		})
	}
	return
}
func (v *buff) ToDef() *gamedef.Buff {
	return &gamedef.Buff{
		BuffId:      v.buffCfg.BuffID,
		ExpireType:  v.ExpireType,
		ExpireV:     v.ExpireV,
		CreateRound: v.createRound,
		UseCount:    v.useCount,
		BuffCount:   v.buffCount,
		IsOverlap:   v.buffCfg.IsOverlap,
	}
}

func (bm *BuffManager) Buffs() []*buff {
	return bm.buffs
}

func (bm *BuffManager) IsBuffsSame(bm2 *BuffManager) bool {
	if len(bm.buffs) != len(bm2.buffs) {
		return false
	}
	for id, _ := range bm.buffs {
		if bm.buffs[id] != bm2.buffs[id] {
			return false
		}
	}
	return true
}
