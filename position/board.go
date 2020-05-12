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
	possibleMoves := make([][]State, 0, 16)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			switch {
			case s.board[i][j] == '.':
			case s.board[i][j] == 'P' || s.board[i][j] == 'p':
				possibleMoves = append(possibleMoves, s.pawnMoves(i, j))
			case s.board[i][j] == 'N' || s.board[i][j] == 'n':
				possibleMoves = append(possibleMoves, s.knightMoves(i, j))
			case s.board[i][j] == 'B' || s.board[i][j] == 'b':
				possibleMoves = append(possibleMoves, s.bishopMoves(i, j))
			case s.board[i][j] == 'R' || s.board[i][j] == 'r':
				possibleMoves = append(possibleMoves, s.rookMoves(i, j))
			case s.board[i][j] == 'Q' || s.board[i][j] == 'q':
				possibleMoves = append(possibleMoves, s.queenMoves(i, j))
			case s.board[i][j] == 'K' || s.board[i][j] == 'k':
				possibleMoves = append(possibleMoves, s.kingMoves(i, j))
			default:
				fmt.Println("Wtf.")

			}
		}
	}
	sum := 0
	for _, ml := range possibleMoves {
		sum += len(ml)
	}
	moves := make([]State, 0, sum)
	for _, ml := range possibleMoves {
		moves = append(moves, ml...)
	}
	return moves
}

func (s State) pawnMoves(pFile int, pRank int) []State {
	return nil
}
func (s State) knightMoves(pFile int, pRank int) []State {
	return nil
}
func (s State) bishopMoves(pFile int, pRank int) []State {
	return nil
}
func (s State) rookMoves(pFile int, pRank int) []State {
	return nil
}
func (s State) queenMoves(pFile int, pRank int) []State {
	return nil
}
func (s State) kingMoves(pFile int, pRank int) []State {
	return nil
}
func (b squares) movePiece(srcFile int, srcRank int, destFile int, destRank int) squares {
	p := b[srcFile][srcRank]
	b[srcFile][srcRank] = '.'
	b[destFile][destRank] = p
	return b
}
