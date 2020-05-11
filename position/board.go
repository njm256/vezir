package position

import (
	"fmt"
	"strings"
)

type squares [8][8]byte

type State struct {
	board       squares
	activeColor string
	castling    string
	ep          string
}

type Game struct {
	state  State
	hclock int
	moveNo int
	prev   map[string]int // maybe map[State]int? seems better to hash the states
}

func (g *Game) hashCode() string {
	return GameToFen(g).String()
}

func GameToFen(s *Game) *Fen {
	//TODO this
	counter := 0
	var str strings.Builder
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if s.state.board[i][j] == '.' {
				counter++
			} else if counter > 0 {
				fmt.Fprintf(&str, "%d", counter)
				counter = 0
			}
			fmt.Fprintf(&str, "%c", s.state.board[i][j])
		}
		if counter > 0 {
			fmt.Fprintf(&str, "%d", counter)
			counter = 0
		}
		if i != 7 {
			str.WriteString("/")
		}
	}
	return &Fen{
		pos:         str.String(),
		activeColor: s.state.activeColor,
		castling:    s.state.castling,
		ep:          s.state.ep,
		hClock:      s.hclock,
		move:        s.moveNo,
	}
}

func (g *Game) move(s State, pawnMoveOrCapture bool) {
	h := g.hashCode()
	g.prev[h] = g.prev[h] + 1
	g.state = s
	if s.activeColor == "w" {
		g.moveNo++
	}
	if pawnMoveOrCapture {
		g.hclock = 0
	} else {
		g.hclock++
	}
	//TODO put something here
}

func (g Game) iMove(s State, pawnMoveOrCapture bool) Game {
	newGame := Game{}
	newGame.state = s
	newGame.prev = make(map[string]int)
	for k, v := range g.prev {
		newGame.prev[k] = v
	}
	h := g.hashCode()
	newGame.prev[h] = newGame.prev[h] + 1
	if s.activeColor == "w" {
		newGame.moveNo = g.moveNo + 1
	}
	if pawnMoveOrCapture {
		newGame.hclock = 0
	} else {
		newGame.hclock = g.hclock + 1
	}
	//TODO put something here
	return newGame
}

func (s State) Moves() (states []State) {
	/*
		var activePlayer army
		var passivePlayer army
		if g.board.activeColor == "w" {
			activePlayer = g.board.wPieces
			passivePlayer = g.board.bPieces
		} else {
			activePlayer = g.board.bPieces
			passivePlayer = g.board.wPieces
		}
	*/

	//this seems super dumb. Better way?
	states = make([]State, 0)
	states = append(states, s.pawnMoves()...)
	states = append(states, s.bishopMoves()...)
	states = append(states, s.knightMoves()...)
	states = append(states, s.rookMoves()...)
	states = append(states, s.queenMoves()...)
	states = append(states, s.kingMoves()...)
	return
}
func (s *State) pawnMoves() []State {
	moves := []State{}
	var marker byte
	var dir int
	if s.activeColor == "w" {
		marker = 'P'
		dir = -1
	} else {
		marker = 'p'
		dir = 1
	}
	for i := 0; i < 8; i++ {
		for j := 1; j < 7; j++ {
			//TODO en passant
			//TODO promotion
			if s.board[i][j] == marker {
				//TODO everything ARGH
			}
		}

	}
	return moves
}
func (s State) knightMoves() []State {
	moves := []State{}
	var marker byte
	var friendlies string
	if s.activeColor =="w" {
		marker = 'N'
		friendlies = 'KQRBNP'
	} else {
		marker = 'n'
		friendlies = 'kqrbnp'
	}
	return nil
}
func (s State) queenMoves() []State {
	return nil
}
func (s State) rookMoves() []State {
	return nil
}
func (s State) bishopMoves() []State {
	return nil
}

func (b squares) movePiece(srcFile int, srcRank int, destFile int, destRank int) squares {
	p := b[srcFile][srcRank]
	b[srcFile][srcRank] = '.'
	b[destFile][destRank] = p
	return b
}
func (s State) kingMoves() []State {
	moves := make([]State, 8)

	return moves
}
