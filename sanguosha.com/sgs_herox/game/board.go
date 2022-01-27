package game

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

type Board struct {
	cells [3][3]*Cell
}

func newBoard() *Board {
	b := &Board{}
	b.init()
	return b
}

func (b *Board) init() {
	for row := int32(0); row < 3; row++ {
		for col := int32(0); col < 3; col++ {
			c := &Cell{
				Position: Position{
					Row: row,
					Col: col,
				},
				Card: nil,
			}
			b.cells[row][col] = c
		}
	}
}

func (b *Board) GetCell(row, col int) *Cell {
	return b.cells[row][col]
}

func (b *Board) GetCellByPos(pos *gamedef.Position) *Cell {
	return b.GetCell(int(pos.Row), int(pos.Col))
}

func (b *Board) SeatCardCount(seatId int32) int32 {
	var count int32 = 0
	for i := int32(0); i < 3; i++ {
		for j := int32(0); j < 3; j++ {
			card, exist := b.cells[i][j].GetCard()
			if !exist {
				continue
			}
			if card.GetPlayer().GetSeatID() == seatId {
				count++
			}
		}
	}
	return count
}

func (b *Board) SeatEmptyCellCount(seatId int32) (count int32) {
	row := SeatRow(seatId)

	for _, v := range b.cells[row] {
		if !v.HasCard() {
			count++
		}
	}
	return
}

func (b *Board) GetCellBySeat(seatId int32) (ret []*Cell) {
	for i := int32(0); i < 3; i++ {
		for j := int32(0); j < 3; j++ {
			cell := b.cells[i][j]
			card, exist := b.cells[i][j].GetCard()
			if !exist {
				continue
			}
			if card.GetPlayer().GetSeatID() == seatId {
				ret = append(ret, cell)
			}
		}
	}
	return
}

func (b *Board) GetCardByPos(pos *gamedef.Position) (*Card, bool) {
	if !isPosValid(pos) {
		return nil, false
	}
	cell := b.GetCellByPos(pos)
	return cell.GetCard()
}

func (b *Board) GetCardsBySeatId(seatId int32) (ret []*Card) {
	for i := int32(0); i < 3; i++ {
		for j := int32(0); j < 3; j++ {
			//cell := b.cells[i][j]
			card, exist := b.cells[i][j].GetCard()
			if !exist {
				continue
			}
			if card.GetPlayer().GetSeatID() == seatId {
				ret = append(ret, card)
			}
		}
	}
	return ret
}

func (b *Board) CellCanAttackPos(cell *Cell) (Position, bool) {
	if !cell.HasCard() {
		return Position{}, false
	}

	for i := int32(0); i < 3; i++ {
		for j := int32(0); j < 3; j++ {
			v := b.cells[i][j]
			card, exist := v.GetCard()
			if !exist {
				continue
			}

			if card.GetPlayer() == cell.GetPlayer() {
				continue
			}

			if calDistanceByPosition(cell.Position, v.Position) == 1 {
				return v.Position, true
			}
		}
	}
	return Position{}, false
}

func (b *Board) SeatPlaceEmptyCell(seatId int32) (ret []*Cell) {
	row := SeatRow(seatId)

	for _, v := range b.cells[row] {
		if !v.HasCard() {
			ret = append(ret, v)
		}
	}
	return
}

//-1 means show all
func (p *Board) ToDef(seatIdView int32) *gamedef.GameBoard {
	cells := make([]*gamedef.Cell, 0, 9)

	b := &gamedef.GameBoard{}
	for row := int32(0); row < 3; row++ {
		for col := int32(0); col < 3; col++ {
			cell := p.cells[row][col]
			cells = append(b.Cells, cell.ToDef(seatIdView))
		}
	}

	return &gamedef.GameBoard{
		Cells: cells,
	}
}
func (b *Board) GetCardById(seatId int32, cardId int32) (*Card, bool) {
	for i := int32(0); i < 3; i++ {
		for j := int32(0); j < 3; j++ {
			//cell := b.cells[i][j]
			card, exist := b.cells[i][j].GetCard()
			if !exist {
				continue
			}
			if card.GetPlayer().GetSeatID() == seatId && card.id == cardId {
				return card, true
			}
		}
	}
	return nil, false
}

func (p *Board) GetCardByUniqueCardId(cardId int32) *Card {
	for _, rows := range p.cells {
		for _, v := range rows {
			if !v.HasCard() {
				continue
			}
			if v.Card.id == cardId {
				return v.Card
			}
		}
	}
	return nil
}

func (b *Board) onEnterPhase(phase gamedef.GamePhase) {
	for row := int32(0); row < 3; row++ {
		for col := int32(0); col < 3; col++ {
			cell := b.cells[row][col]
			card, exist := cell.GetCard()
			if !exist {
				continue
			}
			card.onEnterPhase(phase)
		}
	}

}

type Cell struct {
	Position
	*Card
}

func (p *Cell) GetCard() (*Card, bool) {
	if p.Card == nil {
		return nil, false
	}

	return p.Card, true
}

func (p *Cell) HasCard() bool {
	return p.Card != nil
}

func (p *Cell) SetCard(card *Card) {
	p.Card = card
	card.cell = p
}

func (p *Cell) RemoveCard() *Card {
	ret, ok := p.GetCard()
	if !ok {
		return nil
	}
	ret.cell = nil
	p.Card = nil
	return ret
}

func (p *Cell) ToDef(viewSeatId int32) *gamedef.Cell {
	return &gamedef.Cell{
		Pos:  p.Position.ToDef(),
		Card: p.Card.ToDef(viewSeatId),
	}
}

type Position struct {
	Row int32
	Col int32
}

func (p *Position) ToDef() *gamedef.Position {
	return &gamedef.Position{
		Row: p.Row,
		Col: p.Col,
	}
}

func (b *Board) pierceThrough(src *Card,
	tar *Card) []*Card {

	var cards []*Card
	//src 攻击 tar  从 src 指向tar 的射线 方向上的其他武将

	if src.cell.Row == tar.cell.Row {
		//src 和 tar 处于同一行
		step := src.cell.Position.Col - tar.cell.Col
		pos := Position{tar.cell.Row, tar.cell.Col}
		for { //向射线方向扫描
			pos.Col -= step
			if !IsInsideBoard(pos.Row, pos.Col) {
				break
			}
			if !b.cells[pos.Row][pos.Col].HasCard() {
				continue
			}
			if b.cells[pos.Row][pos.Col].Card.owner.seatId == src.owner.seatId {
				continue
			}
			cards = append(cards, b.cells[pos.Row][pos.Col].Card)
		}
	}

	if src.cell.Col == tar.cell.Col {
		//同一列
		step := src.cell.Position.Row - tar.cell.Row
		pos := Position{tar.cell.Row, tar.cell.Col}
		for {
			pos.Row -= step
			if !IsInsideBoard(pos.Row, pos.Col) {
				break
			}
			if !b.cells[pos.Row][pos.Col].HasCard() {
				continue
			}
			if b.cells[pos.Row][pos.Col].Card.owner.seatId == src.owner.seatId {
				continue
			}
			cards = append(cards, b.cells[pos.Row][pos.Col].Card)
		}
	}

	return cards
}
