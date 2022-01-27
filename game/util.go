package game

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/wanghuiyt/ding"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//工具文件
// 1. 判定位置position是否合法(是否在九宫格内)	isPosValid(pos *gamedef.Position) bool
// 2. 计算两个position之间的距离				calDistance(pos1 *gamedef.Position, pos2 *gamedef.Position) int32
// 3. 计算卡牌到敌方军营的距离  	            calCampDistance(card *Card) int32
// 4. 计算两个position之间的距离（区别于2的position类型） calDistanceByPosition(pos1 Position, pos2 Position) int32
// 5. 取绝对值   							func Abs(x int32) int32
// 6. .......  								 func SeatRow(seatId int32) int
// 7. 判断是否为指定玩家的放置区  				isMyPlaceArea(seatId int32, x int32, y int32) bool
// 8. 随机生成一个[]int32，长度与ids一至     	 func shuffle(ids []int32) []int32
// 9. 将string转为int32						func toInt32(s string) (int32, error)
// 10.   从gamedef.position转换为position	func toPosition(position *gamedef.Position) Position
// 11. 查找id是否在切片内 						IsInSlice(id int32, ids []int32) bool
// 12. 判断卡牌是否在交战区						IsInWarZone(card *Card) bool
// 13. 判断卡牌是否在我方置牌区 				IsInMyZone（card *Card) bool
// 14. 判断卡牌是否在敌方置牌区 					IsInEnemyZone(card *Card) bool
// 15. 修改卡牌的hp和attack并同步给客户端 		func ModifyHPAttack(card *Card, hp int32, attack int32, srcCard *Card)
// 16. 判断卡牌是否为重装武将					  card.IsHeavyCard()
// 17.
// 18.
// 19.
// 20.
// 21.
// 22.
// 23.
// 24.
// 25.
// 26.
// 27.
// 28.

var uidToLog = make(map[uint32]string)

var AhoCorasick = gameutil.Constructor([]string{"zhubeijian", "panshi", "zhangjie"})

func isPosValid(pos *gamedef.Position) bool {
	if pos == nil {
		return false
	}

	if pos.Row < 0 || pos.Row > 2 {
		return false
	}

	if pos.Col < 0 || pos.Col > 2 {
		return false
	}

	return true
}

func calDistance(pos1 *gamedef.Position, pos2 *gamedef.Position) int32 {
	return Abs(pos1.Row-pos2.Row) + Abs(pos1.Col-pos2.Col)
}

//在我方区域 2，交战区 1，敌方区域 0
func calCampDistance(card *Card) int32 {
	if IsInEnemyZone(card) {
		return 0
	}

	if IsInWarZone(card) {
		return 1
	}

	return 2
}

func calDistanceByPosition(pos1 Position, pos2 Position) int32 {
	return Abs(pos1.Row-pos2.Row) + Abs(pos1.Col-pos2.Col)
}

func Abs(x int32) int32 {
	if x > 0 {
		return x
	} else {
		return -x
	}
}

func SeatRow(seatId int32) int {
	if seatId == 0 {
		return 0
	}

	if seatId == 1 {
		return 2
	}

	return -1
}

func shuffle(ids []int32) []int32 {
	count := len(ids)
	newIds := make([]int32, count)
	idxs := gameutil.Perm(count)
	for i, idx := range idxs {
		newIds[i] = ids[idx]
	}
	return newIds
}

func toInt32(s string) (int32, error) {
	i, err := strconv.Atoi(strings.TrimSpace(s))
	return int32(i), err
}

func isMyPlaceArea(seatId int32, x int32, y int32) bool {
	if seatId == 0 && x == 0 {
		return true
	}

	if seatId == 1 && x == 2 {
		return true
	}

	return false
}

func toPosition(position *gamedef.Position) Position {
	return Position{
		Row: position.Row,
		Col: position.Col,
	}
}

func IsInSlice(id int32, ids []int32) bool {
	for _, v := range ids {
		if id == v {
			return true
		}
	}
	return false
}

func IsInWarZone(card *Card) bool {
	return card.cell.Position.Row == 1
}

func IsInMyZone(card *Card) bool {
	return (card.owner.seatId == 0 && card.cell.Position.Row == 0) || (card.owner.seatId == 1 && card.cell.Position.Row == 2)
}

func IsInEnemyZone(card *Card) bool {
	return (card.owner.seatId == 0 && card.cell.Position.Row == 2) || (card.owner.seatId == 1 && card.cell.Position.Row == 0)
}
func ModifyHPAttack(card *Card, hp int32, attack int32, srcCard *Card, spellSkillId int32) {
	//如果变化为0，将不处理，也不会同步给客户端。否则会导致客户端出现 +0 -0 的无意义显示
	if hp > 0 {
		oldHp := card.GetHP()
		card.AddHP(hp)
		SyncChangeHP(card, oldHp, card.GetHP(), srcCard, spellSkillId)
	}
	if hp < 0 {
		hp = -hp
		oldHp := card.GetHP()

		card.SubHP(hp, false)

		SyncChangeHP(card, oldHp, card.GetHP(), srcCard, spellSkillId)
	}

	if attack != 0 {
		oldAt := card.attack
		newAt := card.attack + attack
		if newAt < 0 {
			newAt = 0
		}
		card.attack = newAt
		SyncChangeAttack(card, oldAt, newAt, srcCard)
	}

}

func setAttackDistanceAndNotify(card *Card, distance int32) {
	old := card.attackDistance
	card.attackDistance = distance
	new := card.attackDistance
	if old != new {
		SyncChangeAttackDistance(card, old, new, nil)
	}

}

func setAttackCountAndNotify(card *Card, cnt int32) {
	old := card.maxAtkCnt
	card.maxAtkCnt = cnt
	if old-card.attackCntInTurn != card.maxAtkCnt-card.attackCntInTurn {
		SyncChangeLeftAtkCnt(card, old-card.attackCntInTurn, card.maxAtkCnt-card.attackCntInTurn, nil)
	}

}

func setCanFaceUpAndNotify(card *Card, canFaceUp bool) {
	old := card.canTurnUp
	card.canTurnUp = canFaceUp

	SyncChangeCanTurnUp(card, old, canFaceUp, nil)
}

func FindAllMyOwnCards(card *Card) []*Card {
	p := card.owner
	var myCards []*Card
	for i := int32(0); i < 3; i++ {
		for j := int32(0); j < 3; j++ {
			if !p.game.board.cells[i][j].HasCard() {
				continue
			}
			nowCard := p.game.board.cells[i][j].Card
			if p.seatId == nowCard.owner.seatId { //同势力卡牌
				myCards = append(myCards, nowCard)
			}
		}
	}
	return myCards
}

func FindAllEnemyCards(card *Card) []*Card {
	p := card.owner
	var enemies []*Card
	for i := int32(0); i < 3; i++ {
		for j := int32(0); j < 3; j++ {
			if !p.game.board.cells[i][j].HasCard() {
				continue
			}
			nowCard := p.game.board.cells[i][j].Card
			if p.seatId != nowCard.owner.seatId { //敌对卡牌
				enemies = append(enemies, nowCard)
			}
		}
	}
	return enemies
}

//找出指定玩家的所有卡牌中 拥有 某个技能 的卡牌
func (g *GameBase) skillCards(p *Player, skillId int32) []*Card {
	var cards []*Card
	for row := int32(0); row < 3; row++ {
		for col := int32(0); col < 3; col++ {
			cell := g.board.cells[row][col]
			card, exist := cell.GetCard()
			if !exist {
				continue
			}
			if card.GetPlayer().seatId != p.seatId {
				continue
			}

			if card.HasSkill(skillId) {
				cards = append(cards, card)
			}
		}
	}
	return cards
}

func (g *GameBase) DrawHandCards(players ...*Player) {
	for _, p := range players {
		var cards []*gamedef.PoolCard

		//oldHandCard := p.HandCard.Clone()

		for p.HandCard.Count() < int(g.config.HandCardCount) && p.HandCardPool.Count() > 0 {
			card := p.HandCardPool.RemoveFront()
			cards = append(cards, card)
			p.HandCard.AddCard(card)
		}
		//logrus.Infof(" -!- %v  -!- %v  -!-%v ", p.user.userBrief.Nickname, p.HandCard, p.HandCardPool)
		SyncPlayerState(p, cmsg.SSyncPlayerState_CardPoolCount)

		msg := &cmsg.SNoticeDrawCard{
			OpSeatId: p.GetSeatID(),
			Cards:    cards,
		}

		g.Send2All(msg)

		g.SendSeatMsg(func(seatId int32) proto.Message {
			show := p.GetSeatID() == seatId
			return &cmsg.SSyncHandCard{
				ChangeTypes: cmsg.SSyncHandCard_Draw,
				SeatId:      p.GetSeatID(),
				GetCards:    cards,
				HandCards:   p.HandCard.Cards(show),
				SpellCard:   0,
			}
		})
	}
}

func (card *Card) IsHeavyCard() bool {
	return card.heroCfg.IsHeavy
}

func (g *GameBase) GetHeroIdsByNames(names ...string) []int32 {
	heroes := g.config.GetHeroes()
	var ret []int32
	for _, name := range names {

		for _, v := range *heroes {
			if v.Name == name {
				ret = append(ret, v.HeroID)
				break
			}
		}
	}
	return ret
}

func FindCardsByType(card *Card, selectType gamedef.SelectCardType) []*Card {
	var ret []*Card
	var cards []*Card

	for _, rows := range card.owner.game.board.cells {
		for _, cell := range rows {
			if cell.HasCard() {
				cards = append(cards, cell.Card)
			}
		}
	}

	switch selectType {
	case gamedef.SelectCardType_Any:
		ret = cards
	case gamedef.SelectCardType_other:
		for _, v := range cards {
			if v != card {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_Enemy:
		for _, v := range cards {
			if v.owner != card.owner {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_MyOwn:
		for _, v := range cards {
			if v.owner == card.owner {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_OtherMyOwnFaceUp:
		for _, v := range cards {
			if v.owner == card.owner && v != card && (!v.isBack) {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_OtherMyOwn:
		for _, v := range cards {
			if v.owner == card.owner && v != card {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_EnemyBack:
		for _, v := range cards {
			if v.owner != card.owner && (v.isBack) {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_EnemyFaceUp:
		for _, v := range cards {
			if v.owner != card.owner && (!v.isBack) {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_MyOwnBack:
		for _, v := range cards {
			if v.owner == card.owner && v.isBack {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_MyOwnFaceUp:
		for _, v := range cards {
			if v.owner == card.owner && !v.isBack {
				ret = append(ret, v)
			}
		}
	case gamedef.SelectCardType_OtherMyOwnBack:
		for _, v := range cards {
			if v.owner == card.owner && v != card && v.isBack {
				ret = append(ret, v)
			}
		}

	case gamedef.SelectCardType_NotHeavy: //非重装明将
		for _, v := range cards {
			if !v.isBack && !v.IsHeavyCard() {
				ret = append(ret, v)
			}
		}

	case gamedef.SelectCardType_OneOtherMyOwnAndOneEnemy: //选择一名己方其他武将 和一名 敌方武将
		for _, v := range cards { //找一名敌方武将 ，找到立即break
			if v.owner != card.owner {
				ret = append(ret, v)
				break
			}
		}
		for _, v := range cards { //找一名己方其他武将
			if v.owner == card.owner && v != card {
				ret = append(ret, v)
				break
			}
		}
		//合法的结果 一定是 len==2
	}
	return ret
}

func IsZbj() bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				if ipnet.IP.String() == "10.225.22.191" {
					return true
				}
			}
		}
	}
	return false

}

func IsTest() bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				if ipnet.IP.String() == "10.225.254.248" {
					return true
				}
			}
		}
	}
	return false

}

func DingSendMsgTestToTakeMyLord(info string) {

	(&ding.Webhook{
		AccessToken: "e40826d73faa464f418d6722b65ea5dac3e7b9fe40a90682c100bad409212716",
		Secret:      "SECe6676c2c07ad2463add54d4ff6b21d3c670a35bfc825cf3663133bc32e9a927e",
		//EnableAt:    true,
	}).SendMessage(info)

}

func (p *GameBase) DingSendMsg(info string) {
	at := AhoCorasick.Query(p.players[0].user.userBrief.Nickname + p.players[1].user.userBrief.Nickname)
	if IsZbj() {
		(&ding.Webhook{
			AccessToken: "e853c6d35234f07208ef02fd00e8555ba6046b2df50218ee8964e78fb9070435",
			Secret:      "SEC64d0d6b7ee4fccb9e69ff18432b91b5bf176adbe6cabc8c2492b345e50ac997e",
			EnableAt:    true,
		}).SendMessage(info, at...)
	} else if IsTest() {
		(&ding.Webhook{
			AccessToken: "86871fa45f066b38b6d0b779f5cb151b0221bd0f8311ad194006c095576b4b56",
			Secret:      "SECd97a494654ad4d8ab939ae496a28b491784ca2dc7253b1b769487456a31dd26b",
			EnableAt:    true,
		}).SendMessage(info, at...)
	} else {
		//https://oapi.dingtalk.com/robot/send?access_token=40ddd5547d060e773756a3748350f9002799ead669181319c1f559e72a9c9fba
		(&ding.Webhook{
			AccessToken: "40ddd5547d060e773756a3748350f9002799ead669181319c1f559e72a9c9fba",
			Secret:      "SEC8985298d5292f34289424e1184be6d0176f3f6ff45229b81782e25b25794b11c",
			EnableAt:    true,
		}).SendMessage(info, at...)
	}

}

func (p *Player) Log(info string) {
	logrus.Info(info)
	uidToLog[p.game.roomId] += "\n" + info
}
