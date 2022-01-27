package game

import (
	"fmt"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

// 陈登	斡旋	翻牌技：使一名武将后退一格（无空位无效）
type SkillWoXuan struct {
	HeroSkill
	targetCards []*Card
}

func (ss *SkillWoXuan) CanUse() ([]*Card, error) {
	selectCardType := gamedef.SelectCardType_Any

	//没有可选牌，直接视为 放弃发动机会，但仍会翻面
	if !hasSelectCard(ss.card, selectCardType) {
		return nil, nil
	}

	//检查选的牌是否符合
	//客户端传回的信息 为卡牌位置，根据位置检查类型并找出位置对应的卡牌
	var minSelect int32 = 1
	var maxSelect int32 = 1

	if ret, ok := IsSelectTargetUnique(ss.card, selectCardType, maxSelect); ok {
		ss.targetCards = ret
		return ret, nil
	}

	cards, ok := checkSelectCard(ss.card, minSelect, maxSelect, selectCardType, ss.actSelectParam)
	if !ok {
		return nil, ss
	}
	ss.targetCards = cards
	return cards, nil
}
func (s *SkillWoXuan) OnFaceUp(card *Card) {
	if s.targetCards == nil {
		return
	}

	t := s.targetCards[0]
	row := t.cell.Row
	col := t.cell.Col
	if t.owner.seatId == 0 {
		row--
	} else {
		row++
	}
	if row > 2 || row < 0 {
		// 0 1 2
		return
	}
	g := card.owner.game
	if g.board.cells[row][col].HasCard() {
		return
	}
	g.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动 斡旋：使一名武将（%v）后退一格. 从%v后退至 { %v,%v}", card.GetOwnInfo(), t.GetOwnInfo(), t.cell.Position, row, col))
	StartMoveToCell(t, g.board.cells[row][col])

}
