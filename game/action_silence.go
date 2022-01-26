package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type ActionSilence struct {
	ActionDataBase
	card      *Card
	spellCard *Card
}

func NewActionSilence(game *GameBase, card *Card, spellCard *Card) *ActionSilence {
	ac := &ActionSilence{}
	ac.game = game
	ac.card = card
	ac.spellCard = spellCard
	return ac
}

func (ad *ActionSilence) DoSilence() {
	g := ad.game
	card := ad.card
	g.PostActData(ad)

	ad.PostActStream(func() {
		buffs := card.BuffManager.Buffs()
		for _, v := range buffs { // 触发卡牌上面所有buff的onDisable ，驱散buff效果
			buff, ok := NewBuff(v.buffCfg)
			if !ok {
				continue
			}

			ad.PostActStream(func() {
				buff.OnDisable(card)
			})
		}
	})

	ad.PostActStream(func() {
		card.BuffManager.Clear() //清空卡牌上buffId
		//由于沉默本身算buff，所有清空所有其他buff后，再为卡牌加上一个 沉默buff（断肠buff ，buffID=1）
		card.buffs = append(card.buffs, &buff{buffCfg: &gameconf.BuffConfDefine{BuffID: 1}})

		card.SilenceSkillID() //消去 卡牌本身的技能 ，将skillID置为-1 ，代表已被沉默
		//攻击能力（攻击次数和攻击距离 ）变为白板默认属性
		setAttackDistanceAndNotify(card, 1)
		setAttackCountAndNotify(card, 1)

		oldHp := card.GetHP()
		oldAttack := card.attack
		oldHpMax := card.hpMax

		//进行卡牌被沉默后身材变化的结算
		card.BeSilenced()

		ad.game.GetCurrentPlayer().Log(fmt.Sprintf("%v已被沉默，失去所有buff和技能，身材重置。", card.GetOwnInfo()))
		var syncCard []*cmsg.SSyncCard_Change //存储卡牌所有的变化，统一同步给客户端

		if oldHp != card.GetHP() { //确认血量发生了变化才会同步变化
			syncCard = append(syncCard, &cmsg.SSyncCard_Change{
				ChangeType: cmsg.SSyncCard_HP,
				Old:        oldHp,
				New:        card.GetHP(),
			})
		}
		if card.attack != oldAttack {
			syncCard = append(syncCard, &cmsg.SSyncCard_Change{
				ChangeType: cmsg.SSyncCard_Attack,
				Old:        oldAttack,
				New:        card.attack,
			})
		}
		syncCard = append(syncCard, &cmsg.SSyncCard_Change{
			ChangeType: cmsg.SSyncCard_Buff,
			NewBuffId:  1,
		})
		if oldHpMax != card.hpMax {
			syncCard = append(syncCard, &cmsg.SSyncCard_Change{
				ChangeType: cmsg.SSyncCard_HPMax,
				Old:        oldHpMax,
				New:        card.hpMax,
			})
		}

		card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
			return &cmsg.SSyncCard{
				SeatId:    card.GetPlayer().GetSeatID(),
				Card:      card.ToDef(seatId),
				SpellCard: ad.spellCard.ID(),
				Changes:   syncCard,
			}
		})
	})
	ad.PostActStream(func() {
		g.GetTriggerMgr().Trigger(TriggerType_LoseSkill, ad, card)
	})
}
