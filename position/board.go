package position

type squares uint64

var ranks = map[int]squares{
	1: 0b0000000000000000000000000000000000000000000000000000000011111111,
	2: 0b0000000000000000000000000000000000000000000000001111111100000000,
	3: 0b0000000000000000000000000000000000000000111111110000000000000000,
	4: 0b0000000000000000000000000000000011111111000000000000000000000000,
	5: 0b0000000000000000000000001111111100000000000000000000000000000000,
	6: 0b0000000000000000111111110000000000000000000000000000000000000000,
	7: 0b0000000011111111000000000000000000000000000000000000000000000000,
	8: 0b1111111100000000000000000000000000000000000000000000000000000000,
}

var files = map[byte]squares{
	'a': 0b0000000100000001000000010000000100000001000000010000000100000001,
	'b': 0b0000001000000010000000100000001000000010000000100000001000000010,
	'c': 0b0000010000000100000001000000010000000100000001000000010000000100,
	'd': 0b0000100000001000000010000000100000001000000010000000100000001000,
	'e': 0b0001000000010000000100000001000000010000000100000001000000010000,
	'f': 0b0010000000100000001000000010000000100000001000000010000000100000,
	'g': 0b0100000001000000010000000100000001000000010000000100000001000000,
	'h': 0b1000000010000000100000001000000010000000100000001000000010000000,
}

type army struct {
	Pawns   squares
	Knights squares
	Bishops squares
	Queens  squares
	King    squares
}
type State struct {
	wPieces     army
	bPieces     army
	activeColor string
	castling    string
	ep          string
}

type Game struct {
	board  State
	hclock int
	moveNo int
	prev   map[string]int // maybe map[State]int? seems better to hash the states
}

func (s State) hashCode() string {
	return StateToFen(s).String()
}

func StateToFen(s State) *Fen {
	//TODO this
	return &Fen{}
}

func (g *Game) move(s State, pawnMoveOrCapture bool) {
	h := g.board.hashCode()
	g.prev[h] = g.prev[h] + 1
	g.board = s
	if s.activeColor == "w" {
		g.moveNo++
	}
	if pawnMoveOrCapture {
		g.hclock = 0
	} else {
		g.hclock++
	}
}

func (g Game) iMove(s State, pawnMoveOrCapture bool) Game {
	newGame := Game{}
	newGame.board = s
	newGame.prev = make(map[string]int)
	for k, v := range g.prev {
		newGame.prev[k] = v
	}
	h := g.board.hashCode()
	if s.activeColor == "w" {
		newGame.moveNo = g.moveNo + 1
	}
	if pawnMoveOrCapture {
		newGame.hclock = 0
	} else {
		newGame.hclock = g.hclock + 1
	}

	return newGame
}
