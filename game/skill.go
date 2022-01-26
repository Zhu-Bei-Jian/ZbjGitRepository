package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type Skill interface {
	CanUse() ([]*Card, error)
	PreUseSkill()
	OnFaceUp(card *Card)

	TriggerHandler() []TriggerHandler
	SetSkillCfg(cfg *gameconf.SkillConfDefine)
	GetSkillId() int32
	SetActSelect(actSelect *gamedef.ActSelectParam)
	SetCard(card *Card)
	SetTargets(cards []*Card)
	AddTarget(card *Card)
	SetDataToClient(data []int32)
}

type HeroSkill struct {
	core.ActionDataCore
	skillCfg         gameconf.SkillConfDefine
	actSelectParam   *gamedef.ActSelectParam
	card             *Card
	dataToClients    []int32
	targetCards      []*Card
	dataToClientCard []*Card
}

func (p *HeroSkill) GetBuffId0() int32 {
	if len(p.skillCfg.Buffs) < 1 {
		logrus.Fatal("没有对应的buffID")
		return -1
	}

	return p.skillCfg.Buffs[0]
}

func (p *HeroSkill) CanUse() ([]*Card, error) {
	return nil, nil
}

func (p *HeroSkill) Error() string {
	return "选牌类型或个数异常"
}

//使用技能前触发
func (p *HeroSkill) PreUseSkill() {
	msg := &cmsg.SSyncUseSkill{
		Seat:         p.card.owner.seatId,
		SkillId:      p.skillCfg.SkillID,
		Card:         p.card.ID(),
		TargetCards:  cardIds(p.targetCards...),
		Data:         p.dataToClients,
		TargetSeatId: 0,
	}

	p.card.owner.game.Send2All(msg)
}

//TODO 废除入参card，使用HeroSkill里的对象card
func (p *HeroSkill) OnFaceUp(card *Card) {
	return
}

func (p *HeroSkill) TriggerHandler() []TriggerHandler {
	return []TriggerHandler{}
}

func (p *HeroSkill) SetActSelect(actSelect *gamedef.ActSelectParam) {
	p.actSelectParam = actSelect
}

func (p *HeroSkill) SetSkillCfg(cfg *gameconf.SkillConfDefine) {
	p.skillCfg = *cfg
}

func (p *HeroSkill) SetCard(card *Card) {
	p.card = card
}

func (p *HeroSkill) GetSkillId() int32 {
	return p.skillCfg.SkillID
}

func (p *HeroSkill) SetTargets(cards []*Card) {
	p.targetCards = cards
}

func (p *HeroSkill) AddTarget(card *Card) {
	p.targetCards = append(p.targetCards, card)
}

func (p *HeroSkill) SetDataToClient(data []int32) {
	p.dataToClients = data
}

func NewSkill(skillCfg *gameconf.SkillConfDefine) (Skill, bool) {
	skill, ok := newSkill(skillCfg.SkillID)
	if !ok {
		return nil, false
	}
	skill.SetSkillCfg(skillCfg)
	return skill, true
}

func newSkill(skillId int32) (Skill, bool) {
	switch skillId {
	case 1:
		return &SkillDuanChang{}, true
	case 2:
		return &SkillDianHou{}, true
	case 3:
		return &SkillJianXiong{}, true
	case 4:
		return &SkillZhiChi{}, true
	case 5:
		return &SkillGuoSe{}, true
	case 6:
		return &SkillJiXi{}, true
	case 7:
		return &SkillQiangXi{}, true
	case 8:
		return &SkillLiJian{}, true
	case 9:
		return &SkillBaoNue{}, true
	case 10:
		return &SkillXuanHuo{}, true
	case 11:
		return &SkillPoXi{}, true
	case 12:
		return &SkillXianZhen{}, true
	case 13:
		return &SkillWuSheng{}, true
	case 14:
		return &SkillTianDu{}, true
	case 15:
		return &SkillJieFan{}, true
	case 16:
		return &SkillZhenDu{}, true
	case 17:
		return &SkillQingNang{}, true
	case 18:
		return &SkillKuRou{}, true
	case 19:
		return &SkillJiQiao{}, true
	case 20:
		return &SkillLieGong{}, true
	case 21:
		return &SkillWeiMu{}, true
	case 22:
		return &SkillShangWu{}, true
	case 23:
		return &SkillFengCheng{}, true
	case 24:
		return &SkillJiJiang{}, true
	case 25:
		return &SkillDiMeng{}, true
	case 26:
		return &SkillQianXun{}, true
	case 27:
		return &SkillWuShuang{}, true
	case 28:
		return &SkillTieJi{}, true
	case 29:
		return &SkillShiShou{}, true
	case 30:
		return &SkillHuoShou{}, true
	case 31:
		return &SkillZiYuan{}, true
	case 32:
		return &SkillFengChu{}, true
	case 33:
		return &SkillJianJie{}, true
	case 34:
		return &SkillTianMing{}, true
	case 35:
		return &SkillJiAng{}, true
	case 36:
		return &SkillYingHun{}, true
	case 37:
		return &SkillZhiHeng{}, true
	case 38:
		return &SkillXiaoJi{}, true
	case 39:
		return &SkillTianYi{}, true
	case 40:
		return &SkillRaoShe{}, true
	case 41:
		return &SkillShiXue{}, true
	case 42:
		return &SkillGangLie{}, true
	case 43:
		return &SkillTianXiang{}, true
	case 44:
		return &SkillDuanLiang{}, true
	case 45:
		return &SkillXiongHuo{}, true
	case 46:
		return &SkillZhuHai{}, true
	case 47:
		return &SkillLuoYi{}, true
	case 48:
		return &SkillJieMing{}, true
	case 49:
		return &SkillGuHuo{}, true
	case 50:
		return &SkillYiZhong{}, true
	case 51:
		return &SkillLuanJi{}, true
	case 52:
		return &SkillWeiDi{}, true
	case 53:
		return &SkillPaoXiao{}, true
	case 54:
		return &SkillQiaoBian{}, true
	case 55:
		return &SkillYinLei{}, true
	case 56:
		return &SkillTuXi{}, true
	case 57:
		return &SkillXianTu{}, true
	case 58:
		return &SkillLongDang{}, true
	case 59:
		return &SkillLuoShen{}, true
	case 60:
		return &SkillFanJian{}, true
	case 61:
		return &SkillWeiZhong{}, true
	case 62:
		return &SkillKanPo{}, true
	case 63:
		return &SkillManYi{}, true
	case 64:
		return &SkillTongMing{}, true
	case 65:
		return &SkillLiRang{}, true
	case 66:
		return &SkillJuShou{}, true
	case 67:
		return &SkillShuShen{}, true
	case 68:
		return &SkillPoJun{}, true
	case 69:
		return &SkillFangQuan{}, true
	case 70:
		return &SkillXuShen{}, true
	case 71:
		return &SkillSongCi{}, true
	case 72:
		return &SkillHanYong{}, true
	case 73:
		return &SkillBuQu{}, true
	case 74:
		return &SkillXingShang{}, true
	case 75:
		return &SkillZiShou{}, true
	case 76:
		return &SkillWoXuan{}, true
	//case 75:

	//zzz
	default:
		return nil, false
	}
}

//根据选牌类型 、个数 判断是否有可选的牌 //如果无牌可选，直接放弃发动技能，但会翻至正面。此时不依赖客户端的选牌信息
func hasSelectCard(card *Card, selectType gamedef.SelectCardType) bool {
	var flag = true
	switch selectType {
	case gamedef.SelectCardType_OneOtherMyOwnAndOneEnemy:
		cards1 := FindCardsByType(card, gamedef.SelectCardType_OtherMyOwn)
		cards2 := FindCardsByType(card, gamedef.SelectCardType_Enemy)
		if len(cards1) < 1 || len(cards2) < 1 {
			flag = false
		}
	default:
		cards := FindCardsByType(card, selectType)
		if len(cards) < 1 {
			flag = false
		}
	}
	if !flag {
		card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf(" %v翻面，检测场上没有可选目标，不发动该武将的翻牌技能(%v)，但仍翻至正面.", card.GetOwnInfo(), card.skillCfg.Name))
	}

	return flag
}

//根据选牌类型 、个数 判断选择 是否唯一 //如果唯一，那么直接选择唯一目标 作为技能目标，此时不依赖客户端的选牌信息
func IsSelectTargetUnique(card *Card, selectType gamedef.SelectCardType, selectCount int32) ([]*Card, bool) {

	switch selectType {
	case gamedef.SelectCardType_OneOtherMyOwnAndOneEnemy:
		cards1 := FindCardsByType(card, gamedef.SelectCardType_OtherMyOwn)
		cards2 := FindCardsByType(card, gamedef.SelectCardType_Enemy)
		if len(cards1) == 1 && len(cards2) == 1 { //离间 的选择方式 唯一
			card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 离间可选目标组合唯一（%v和%v），该两名武将进行决斗.", card.GetOwnInfo(), cards1[0].GetOwnInfo(), cards2[0].GetOwnInfo()))
			return append(cards1, cards2...), true
		}
	default:
		cards := FindCardsByType(card, selectType)
		if len(cards) == int(selectCount) {
			var logInfo []string
			for _, v := range cards {
				logInfo = append(logInfo, v.GetOwnInfo())
			}
			card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动%v,此次选择目标个数为%v,检测到场上可选目标组合唯一（%v），将直接对唯一目标组合发动该技能.", card.GetOwnInfo(), card.skillCfg.Name, selectCount, logInfo))
			return cards, true
		}
	}
	return nil, false
}

func checkSelectCard(srcCard *Card, minSelect int32, maxSelect int32, cardType gamedef.SelectCardType, selectParam *gamedef.ActSelectParam) ([]*Card, bool) {

	//positon 去重 ，保留唯一position
	var tPos []*gamedef.Position
	for _, pos1 := range selectParam.Positions {
		hasSamePos := false
		for _, pos2 := range tPos {
			if pos1.Row == pos2.Row && pos1.Col == pos2.Col {
				hasSamePos = true
				break
			}
		}
		if !hasSamePos {
			tPos = append(tPos, pos1)
		}
	}

	// 检查 选择个数区间
	if int32(len(tPos)) < minSelect || int32(len(tPos)) > maxSelect {
		return nil, false
	}
	// 将位置 转换为 对应卡牌
	var cards []*Card
	for _, pos := range tPos {
		card, exist := srcCard.owner.game.board.GetCardByPos(pos)
		if !exist { //传了空的位置
			return nil, false
		}
		cards = append(cards, card)
	}

	// 根据selectCardType检查牌合法性
	switch cardType {
	case gamedef.SelectCardType_Enemy:
		for _, card := range cards {
			if card.owner == srcCard.owner {
				return nil, false
			}
		}
	case gamedef.SelectCardType_MyOwn:
		for _, card := range cards {
			if card.owner != srcCard.owner {

				return nil, false
			}
		}
	case gamedef.SelectCardType_OtherMyOwnFaceUp:
		for _, card := range cards {
			if card.owner != srcCard.owner || (card.cell.Position.Row == srcCard.cell.Row && card.cell.Position.Col == srcCard.cell.Col) || card.isBack {

				return nil, false
			}
		}
	case gamedef.SelectCardType_OtherMyOwn:
		for _, card := range cards {
			if card.owner != srcCard.owner || (card.cell.Position.Row == srcCard.cell.Row && card.cell.Position.Col == srcCard.cell.Col) {

				return nil, false
			}
		}
	case gamedef.SelectCardType_EnemyBack:
		for _, card := range cards {
			if card.owner == srcCard.owner || (!card.isBack) {

				return nil, false
			}
		}
	case gamedef.SelectCardType_EnemyFaceUp:
		for _, card := range cards {
			if card.owner == srcCard.owner || (card.isBack) {

				return nil, false
			}
		}
	case gamedef.SelectCardType_MyOwnBack:
		for _, card := range cards {
			if card.owner != srcCard.owner || (!card.isBack) {

				return nil, false
			}
		}
	case gamedef.SelectCardType_MyOwnFaceUp:
		for _, card := range cards {
			if card.owner != srcCard.owner || (card.isBack) {

				return nil, false
			}
		}
	case gamedef.SelectCardType_OtherMyOwnBack:
		for _, card := range cards {
			if card.owner != srcCard.owner || (card.cell.Position.Row == srcCard.cell.Row && card.cell.Position.Col == srcCard.cell.Col) || (!card.isBack) {

				return nil, false
			}
		}

	case gamedef.SelectCardType_NotHeavy: //非重装明将
		for _, card := range cards {
			if card.isBack || (card.IsHeavyCard()) {

				return nil, false
			}
		}
	case gamedef.SelectCardType_OneOtherMyOwnAndOneEnemy: //选择一名己方其他武将 和一名 敌方武将
		if len(cards) != 2 {

			return nil, false
		}
		if cards[0].owner == cards[1].owner {

			return nil, false
		}
	case gamedef.SelectCardType_other:
		for _, card := range cards {
			if srcCard == card {
				return nil, false
			}
		}
	}
	var cardInfo []string
	for _, v := range cards {
		cardInfo = append(cardInfo, v.GetOwnInfo())
	}
	srcCard.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v发动翻牌技%v,玩家选择目标为：%v ", srcCard.GetOwnInfo(), srcCard.skillCfg.Name, cardInfo))
	return cards, true
}
