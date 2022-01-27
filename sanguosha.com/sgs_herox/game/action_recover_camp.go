package game

//军营恢复HP
type ActionRecoverCamp struct {
	ActionDataBase

	pos       *Position
	player    *Player
	spellCard *Card
	skillId   int32

	recoverValue int32
}

func NewActionRecoverCamp(game *GameBase, player *Player, spellCard *Card, recover int32) *ActionRecoverCamp {
	ar := &ActionRecoverCamp{}
	ar.game = game
	ar.player = player
	ar.spellCard = spellCard
	ar.skillId = 0
	ar.recoverValue = recover
	return ar
}

func (ac *ActionRecoverCamp) DoRecover() {
	ac.game.PostActData(ac)
	ac.PostActStream(func() {

		old := ac.player.GetHP()
		ac.player.AddHP(ac.recoverValue)

		SyncCampChangeHP(ac.player, ac.spellCard, old, ac.recoverValue)
	})
}
