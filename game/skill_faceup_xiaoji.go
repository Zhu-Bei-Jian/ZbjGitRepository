package game

import "sanguosha.com/sgs_herox/proto/gameconf"

type SkillXiaoJi struct {
	HeroSkill
}

func (ss *SkillXiaoJi) OnFaceUp(card *Card) {
	M := card.owner.game.board.cells
	row := card.cell.Position.Row
	col := card.cell.Position.Col
	direction := [5][2]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {0, 0}}
	g := card.owner.game
	g.PostActData(ss)
	for _, dir := range direction {
		if !IsInsideBoard(row+dir[0], col+dir[1]) {

			continue
		}

		if M[row+dir[0]][col+dir[1]].HasCard() && M[row+dir[0]][col+dir[1]].Card.owner.seatId == card.owner.seatId {
			tDir := dir
			ss.PostActStream(func() {
				StartGetBuff(M[row+tDir[0]][col+tDir[1]].Card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
			})

		}
	}

}
func IsInsideBoard(row int32, col int32) bool {
	return row <= 2 && row >= 0 && col <= 2 && col >= 0
}
