package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type Card struct {
	id       int32 //游戏内卡牌唯一ID
	heroCfg  *gameconf.HeroDefine
	skillCfg *gameconf.SkillConfDefine

	owner *Player
	cell  *Cell //卡牌所在格子

	isBack    bool
	canTurnUp bool //是否可以翻到正面（本回合刚放置的卡牌不能翻面，以及被孙权制衡翻回背面的卡牌也不能立即翻至正面

	buffHP int32 //来自buff的HP，扣血结算时优先扣除buffHP，被沉默后这部分hp清零。 加血时优先级在本身hp之后
	hp     int32 //本身hp，扣血结算时当buffHP为零时，再开始扣除本身hp。  加血时，优先加到本身HP，如果加满本身hp，再将溢出的hp加到buffHP上，最多可以加到到血量总上限为止
	hurtHP int32 //扣除在本身HP上的血量 //伤害继承
	hpMax  int32 //实际血量（buffhp+hp）的总上限，

	attack int32

	attackDistance  int32 //攻击距离（默认为1)
	attackCntInTurn int32 //本回合卡牌已经攻击的次数
	maxAtkCnt       int32 //本回合卡牌 最多可以攻击几次

	skillId int32 // -1 代表被沉默 ，0代表未翻面 ，>0代表正面拥有对应ID的技能

	skillData map[int32]interface{}

	BuffManager
}

func newBackCard(cardId int32, heroCfg *gameconf.HeroDefine, skillCfg *gameconf.SkillConfDefine, owner *Player) *Card {

	card := &Card{
		id:     cardId,
		owner:  owner,
		isBack: true,
		buffHP: 0,

		heroCfg:        heroCfg,
		skillCfg:       skillCfg,
		hp:             heroCfg.HP,
		attack:         heroCfg.Attack,
		canTurnUp:      false,
		attackDistance: 1,
	}

	card.hp = heroCfg.HP
	card.hpMax = heroCfg.HP
	card.attack = heroCfg.Attack
	card.attackCntInTurn = 0
	card.maxAtkCnt = 1
	card.skillData = map[int32]interface{}{}
	return card
}

//沉默 伤害继承机制 下 的 血量变化和攻击力变化

func (p *Card) AddHP(v int32) {
	for p.hurtHP > 0 && v > 0 {
		p.hurtHP--
		p.hp++
		v--
	}
	for v > 0 && p.hp+p.buffHP < p.hpMax {
		p.buffHP++
		v--
	}
}
func (p *Card) SubHP(v int32, isDamage bool) {
	for p.buffHP > 0 && v > 0 {
		v--
		p.buffHP--
	}
	for v > 0 {
		v--
		p.hp--
		if isDamage {
			p.hurtHP++
		}
	}
}
func (p *Card) GetHP() int32 {
	return p.buffHP + p.hp
}

func (p *Card) AddHpMax(v int32) {
	p.hpMax += v
}
func (p *Card) SubHpMax(v int32) {
	p.hpMax -= v
	for v > 0 && p.buffHP > 0 && p.buffHP+p.hp > p.hpMax {
		p.buffHP--
		v--
	}
	for v > 0 && p.hp > 0 {
		v--
		p.hp--
	}
}
func (p *Card) IsDead() bool {
	return p.GetHP() <= 0
}

func (p *Card) AddAttack(v int32) {
	p.attack += v
}
func (p *Card) SubAttack(v int32) {
	p.attack -= v
	if p.attack < 0 {
		p.attack = 0
	}
}
func (p *Card) GetAttack() int32 {
	return p.attack
}

func (p *Card) BeSilenced() { //被沉默后，重置
	p.hpMax = p.heroCfg.HP
	p.attack = p.heroCfg.Attack
	p.buffHP = 0
	p.hp = p.hpMax - p.hurtHP
	//伤害继承 （只继承扣除在hp上 ，扣除在buffHP上的不算在内）
}

func (p *Card) GetPlayer() *Player {
	return p.owner
}
func (p *Card) GetUserName() string {
	return p.owner.user.userBrief.Nickname
}
func (p *Card) GetName() string {
	return p.heroCfg.Name
}

func (p *Card) GetOwnInfo() string { //用于方便打印归属信息
	return "[" + p.GetUserName()[3:] + "]" + p.GetName()
}

func (p *Card) ID() int32 {
	if p == nil {
		return 0
	}
	return p.id
}

func (p *Card) ToDef(viewSeatId int32) *gamedef.Card {
	if p == nil {
		return nil
	}

	pos := &gamedef.Position{}
	if p.cell != nil {
		pos = p.cell.Position.ToDef()
	}
	c := &gamedef.Card{
		Id:             p.id,
		IsBack:         p.isBack,
		BuffHp:         p.buffHP,
		HeroId:         0,
		Hp:             0,
		Attack:         0,
		Buffs:          p.BuffManager.ToDef(),
		Position:       pos,
		SkillId:        p.skillId,
		CanTurnUp:      p.canTurnUp,
		AttackDistance: p.attackDistance,
		LeftAttackCnt:  p.maxAtkCnt - p.attackCntInTurn,
	}

	if viewSeatId == p.owner.GetSeatID() || viewSeatId == -1 || p.isBack == false {
		c.HeroId = p.heroCfg.HeroID
		c.Hp = p.GetHP()
		c.HpMax = p.hpMax
		c.Attack = p.attack
	}

	return c
}

func (p *Card) TurnUpFace() {
	p.isBack = false
	p.skillId = p.heroCfg.SkillID
}

func (p *Card) GetSkillId() int32 {
	return p.skillId
}

func (p *Card) SilenceSkillID() { //卡牌被沉默时，skillID 以 -1 代表
	p.skillId = -1
}
func (p *Card) FaceDownSkillID() { //翻回背面（制衡）且没有被沉默时 归零
	if p.skillId != -1 {
		p.skillId = 0
	}
}

func (p *Card) HasSkill(skillId int32) bool {
	return p.skillId == skillId
}
func (p *Card) onEnterPhase(phase gamedef.GamePhase) {
	switch phase {
	case gamedef.GamePhase_PhaseBegin:
		p.attackCntInTurn = 0
		SyncChangeLeftAtkCnt(p, 0, p.maxAtkCnt, nil)
	case gamedef.GamePhase_PhaseEnd:
		p.canTurnUp = true
		SyncChangeCanTurnUp(p, false, p.canTurnUp, nil)
	default:

	}
}

func cardIds(cards ...*Card) (ret []int32) {
	for _, v := range cards {
		ret = append(ret, v.ID())
	}
	return
}
