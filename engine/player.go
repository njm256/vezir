package engine

import (
	"github.com/notnil/chess"
)

type Game struct {
	*chess.Game
}

func NewGame() *Game {
	return &Game{chess.NewGame()}
}

func (g *Game) ActivePlayer() int {
	s := g.Position().Turn()
	if s == chess.White {
		return 1
	}
	return -1
}

func (g *Game) Result() (bool, int) {
	switch o := g.Outcome(); o {
	case chess.NoOutcome:
		return false, 0
	case chess.WhiteWon:
		return true, 1
	case chess.BlackWon:
		return true, -1
	case chess.Draw:
		return true, 0
	default:
		return false, -1
	}
}

func (g *Game) Moves() []GameState {
	moves := g.ValidMoves()
	ret := make([]GameState, 0, len(moves))
	for _, m := range moves {
		ng := Game{g.Clone()}
		ng.Move(m)
		ret = append(ret, &ng)
	}

	return ret
}
