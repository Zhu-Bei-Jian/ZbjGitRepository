package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"

	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"time"
)

func (g *GameBase) onReqAct(p *Player, req *cmsg.CReqAct) {
	cardId := req.CardId
	actType := req.ActType
	targetPos := req.TargetPos
	seatId := p.GetSeatID()

	resp := &cmsg.SRespAct{
		CardId:    cardId,
		ActType:   actType,
		TargetPos: targetPos,
	}

	//严格控制 reqAct  和 respAct 一对一
	if p != g.GetCurrentPlayer() {
		resp.ErrCode = cmsg.SRespAct_ErrNotMyTurn
		p.SendMsg(resp)
		return
	}

	if !isPosValid(targetPos) {
		resp.ErrCode = cmsg.SRespAct_ErrInputParam
		p.SendMsg(resp)
		return
	}

	if actType == gamedef.ActType_PlaceCard {
		g.placeCard(p, cardId, targetPos, resp) //内部 回复resp
		return
	}

	card, exist := g.board.GetCardById(seatId, cardId)
	if !exist {
		resp.ErrCode = cmsg.SRespAct_ErrCardNotExist
		p.SendMsg(resp)
		return
	}

	//srcPos := card.cell.Position

	switch actType {
	case gamedef.ActType_TurnUpCard:
		if !card.isBack {
			resp.ErrCode = cmsg.SRespAct_ErrCardAlreadyFaceUp
			p.SendMsg(resp)
			return
		}
		if !card.canTurnUp {
			resp.ErrCode = cmsg.SRespAct_ErrCardRoundNotEnough
			p.SendMsg(resp)
			return
		}

		a := NewActionFaceUpCard(g, card, false, nil, req.ActSelect)
		if req.NeedSelect && req.ActSelect == nil {
			resp.ErrCode = cmsg.SRespAct_ErrNoSelectTarget
			p.SendMsg(resp)
			return
		}
		a.DoFaceUp(func(err error) {
			defer p.SendMsg(resp)
			if err != nil {
				resp.ErrCode = cmsg.SRespAct_ErrActSelectNotRight
			} else {
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v翻至正面", card.GetOwnInfo()))
			}

		})

	case gamedef.ActType_MoveCard:
		if calDistance(card.cell.Position.ToDef(), targetPos) != 1 {
			resp.ErrCode = cmsg.SRespAct_ErrDistanceLimit
			p.SendMsg(resp)
			return
		}

		ac := NewActionCheck(g, card, actType, nil)
		ac.DoCheck(func() {
			defer p.SendMsg(resp)
			if p.GetAP() < ac.ap[MOVE_CARD] {
				resp.ErrCode = cmsg.SRespAct_ErrAPNotEnough
				return
			}
			SubAPAndSync(p, ac.ap[MOVE_CARD])
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v从%v移动到%v", card.GetOwnInfo(), card.cell.Position, toPosition(targetPos)))
			a := NewActionMove(g, card, toPosition(targetPos), p)
			a.Do()

		})
	case gamedef.ActType_AttackCard:
		if card.isBack {
			resp.ErrCode = cmsg.SRespAct_ErrAttackCardNotFaceUp
			p.SendMsg(resp)
			return
		}

		targetCard, exist := g.board.GetCardByPos(targetPos)
		if !exist {
			resp.ErrCode = cmsg.SRespAct_ErrTargetPosCardNotExist
			p.SendMsg(resp)
			return
		}

		distance := calDistance(card.cell.Position.ToDef(), targetPos)

		ac := NewActionCheck(g, card, actType, targetCard)
		ac.DoCheck(func() {
			defer p.SendMsg(resp)
			if distance > ac.card.attackDistance {
				resp.ErrCode = cmsg.SRespAct_ErrDistanceLimit
				return
			}

			if card.attackCntInTurn >= card.maxAtkCnt {
				resp.ErrCode = cmsg.SRespAct_ErrAttackCountNotEnough
				return
			}

			if p.GetAP() < ac.ap[ATTACK_CARD] {
				resp.ErrCode = cmsg.SRespAct_ErrAPNotEnough
				return
			}

			SubAPAndSync(p, ac.ap[ATTACK_CARD])

			ac := NewActionActAttack(g, card, []*Card{targetCard}, card, ac.effectSkills)
			ac.DoAttack()

		})
	case gamedef.ActType_AttackCamp:
		if card.isBack {
			resp.ErrCode = cmsg.SRespAct_ErrAttackCardNotFaceUp
			p.SendMsg(resp)
			return
		}

		distance := calCampDistance(card)

		ac := NewActionCheck(g, card, actType, nil)
		ac.DoCheck(func() {
			defer p.SendMsg(resp)
			if distance > ac.maxCampDistance {
				resp.ErrCode = cmsg.SRespAct_ErrDistanceLimit
				return
			}

			if card.attackCntInTurn >= card.maxAtkCnt {
				resp.ErrCode = cmsg.SRespAct_ErrAttackCountNotEnough
				return
			}
			if !ac.canAttackCamp {
				resp.ErrCode = cmsg.SRespAct_ErrCampCannotAttack
				return
			}
			logrus.Infof("game_msg.go Line:169,玩家AP %v ，卡牌攻击需要消耗的AP%v ,roomId: %v", p.GetAP(), ac.ap[ATTACK_CAMP], g.roomId)
			if p.GetAP() < ac.ap[ATTACK_CAMP] {
				resp.ErrCode = cmsg.SRespAct_ErrAPNotEnough
				return
			}
			SubAPAndSync(p, ac.ap[ATTACK_CAMP])

			targetPlayer := g.GetNextAlivePlayer(seatId)
			a := NewActionAttackCamp(g, targetPlayer, card)
			a.DoAttack()

		})
	default:
		resp.ErrCode = cmsg.SRespAct_ErrActionNotSupport
		p.SendMsg(resp)
		return
	}

	//g.Send2All(&cmsg.SSyncAct{
	//	ActType:   actType,
	//	OpSeatId:  seatId,
	//	SrcPos:    srcPos.ToDef(),
	//	TargetPos: targetPos,
	//})

	g.StopWaiting()
}

func (g *GameBase) placeCard(p *Player, cardId int32, targetPos *gamedef.Position, resp *cmsg.SRespAct) {

	card, exist := p.HandCard.GetCard(cardId)
	if !exist {
		resp.ErrCode = cmsg.SRespAct_ErrCardNotExist
		p.SendMsg(resp)
		return
	}

	heroCfg, exist := g.config.Hero.GetHero(card.HeroId)
	if !exist {
		resp.ErrCode = cmsg.SRespAct_ErrHeroCfgNotExist
		p.SendMsg(resp)
		return
	}

	skillCfg, exist := g.config.Skill.GetSkill(heroCfg.SkillID)
	if !exist {
		resp.ErrCode = cmsg.SRespAct_ErrSkillCfgNotExist
		p.SendMsg(resp)
		return
	}

	if !p.HandCard.Has(cardId) {
		resp.ErrCode = cmsg.SRespAct_ErrCardNotExist
		p.SendMsg(resp)
		return
	}

	if !isMyPlaceArea(p.GetSeatID(), targetPos.Row, targetPos.Col) {
		resp.ErrCode = cmsg.SRespAct_ErrTargetPosNotMyPlaceArea
		p.SendMsg(resp)
		return
	}

	cell := g.board.GetCellByPos(targetPos)
	if cell.HasCard() {
		resp.ErrCode = cmsg.SRespAct_ErrTargetPosNotEmpty
		p.SendMsg(resp)
		return
	}

	ad := NewActionCheckPlaceCard(g, p, heroCfg, skillCfg)
	ad.DoCheck(func() {
		defer p.SendMsg(resp)
		if p.GetAP() < ad.ap {
			resp.ErrCode = cmsg.SRespAct_ErrAPNotEnough
			return
		}

		SubAPAndSync(p, ad.ap)

		//oldHandCard := p.HandCard.Clone()

		p.HandCard.Remove(cardId)

		card := newBackCard(cardId, heroCfg, skillCfg, p)
		cell.SetCard(card)

		if card.heroCfg.IsHeavy {
			card.isBack = false
			card.skillId = card.skillCfg.SkillID
			for _, bf := range ad.buffs { //重装武将的光环buff,自放置之初即开始生效
				switch bf.effectType {
				case EffectType_camp:
					AddCampBuff(card, bf.buffCfg.GetBuffID(), bf.ExpireType, bf.ExpireV)
				case EffectType_Card:
					for _, t := range bf.targetCards {
						StartGetBuff(t, bf.buffCfg.GetBuffID(), bf.ExpireType, bf.ExpireV, card)
					}
				default:
					continue
				}

			}
		}

		g.SendSeatMsg(func(seatId int32) proto.Message {
			targetP := []*gamedef.Position{targetPos}
			return &cmsg.SSyncAct{
				ActType:   gamedef.ActType_PlaceCard,
				OpSeatId:  p.GetSeatID(),
				OpCard:    card.id,
				SpellCard: 0,
				TargetPos: targetP,
			}
		})

		SyncChangePlace(card)
		g.GetCurrentPlayer().Log(fmt.Sprintf("%v放置于%v", card.GetOwnInfo(), card.cell.Position))

		g.SendSeatMsg(func(seatId int32) proto.Message {
			show := p.GetSeatID() == seatId
			return &cmsg.SSyncHandCard{
				ChangeTypes: cmsg.SSyncHandCard_Place,
				SeatId:      p.GetSeatID(),
				LoseCards:   []*gamedef.PoolCard{&gamedef.PoolCard{CardId: cardId}},
				HandCards:   p.HandCard.Cards(show),
				SpellCard:   0,
			}
		})
	})

	g.StopWaiting()
}

func (g *GameBase) onReqCancelCurOpt(p *Player, req *cmsg.CReqCancelCurOpt) {
	resp := &cmsg.SRespCancelCurOpt{}
	defer p.SendMsg(resp)

	currPlayer := g.GetCurrentPlayer()
	if p != currPlayer {
		resp.ErrCode = cmsg.SRespCancelCurOpt_ErrNotMyTurn
		return
	}

	g.ForceTimeout()
}

//请求游戏场景,客户端发出此消息代表客户端场景已准备好
func (g *GameBase) onReqGameScene(p *Player, req *cmsg.CReqGameScene) {
	p.setClientReady(true)

	seats := make([]*cmsg.SRespGameScene_Seat, 0)
	for _, v := range g.players {
		if v == nil {
			continue
		}
		seat := &cmsg.SRespGameScene_Seat{
			SeatId:       v.seatId,
			UserId:       v.GetUser().userId,
			UserBrief:    v.GetUser().userBrief,
			ConnectState: v.connectState,
			GameSeat:     v.toSeatInfo(false),
		}
		seats = append(seats, seat)
	}

	//waitPlayer, waitEndTime := g.waitMgr.waitInfo()

	resp := &cmsg.SRespGameScene{
		ErrCode:      0,
		GameMode:     g.gameMode,
		RoomSetting:  g.setting,
		RoomNO:       g.roomNO,
		Phase:        g.Phase,
		PhaseEndTime: g.PhaseEndTime,
		//OpSeatIds:    getPlayerSeatIds(waitPlayer...),
		//OpEndTime:    waitEndTime,
		Seats:       seats,
		MySeatId:    p.seatId,
		LookerType:  p.lookerType,
		LookerCount: int32(len(g.lookers)),
		Board:       g.board.ToDef(p.GetSeatID()),
		ServerTime:  time.Now().Unix(),
	}

	p.SendMsg(resp)

	g.resendLastOpMsg(p)
}
