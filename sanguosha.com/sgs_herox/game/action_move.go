package game

type ActionMove struct {
	ActionDataBase

	card        *Card
	to          Position
	spellPlayer *Player
	skillId     int32
}

func NewActionMove(game *GameBase, card *Card, to Position, spellPlayer *Player) *ActionMove {
	ac := &ActionMove{}
	ac.game = game
	ac.spellPlayer = spellPlayer
	ac.card = card
	ac.to = to
	return ac
}

func (ac *ActionMove) Do() {
	ac.game.PostActData(ac)

	ac.PostActStream(func() {
		fromPos := ac.card.cell.Position
		card := ac.card.cell.RemoveCard()
		ac.game.board.GetCell(int(ac.to.Row), int(ac.to.Col)).SetCard(card)

		SyncChangePos(&fromPos, ac.game, card)
	})
	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_MoveCard, ac)
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_AfterPosChange, ac)
	})
}
