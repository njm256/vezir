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
	//TODO refactor by making a copy instead
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
	moves := make([]State, 0, 12)
	var allies string
	var enemies string
	var dir int
	if s.activeColor == "w" {
		allies = "KQRNBP"
		enemies = "kqrnbp"
		dir = -1
	} else {
		allies = "kqrnbp"
		enemies = "KQRNBP"
		dir = 1
	}
	//TODO everything argh
	//have to handle first move, en passant, and promotion
	return moves
}
func (s State) knightMoves(pFile int, pRank int) []State {
	moves := make([]State, 0, 8)
	var allies string
	if s.activeColor == "w" {
		allies = "KQRNBP"
	} else {
		allies = "kqrnbp"
	}
	offs := [4]int{-2, -1, 1, 2}
	for i := range offs {
		for j := range offs {
			//skip if target square is out of bounds or holds a friendly piece.
			if i == j || pFile+i < 0 || pFile+i > 7 || pRank+j < 0 || pRank+j > 7 ||
				strings.IndexByte(allies, s.board[pFile+i][pRank+j]) != -1 {
				continue
			}
			t := s
			t.board = t.board.movePiece(pFile, pRank, pFile+i, pRank+j)
			t.ep = ""
			if s.activeColor == "w" {
				t.activeColor = "b"
			} else {
				t.activeColor = "w"
			}
			moves = append(moves, t)
		}
	}
	return moves
}

func (s State) bishopMoves(pFile int, pRank int) []State {
	moves := make([]State, 0, 13)
	var allies string
	if s.activeColor == "w" {
		allies = "KQRNBP"
	} else {
		allies = "kqrnbp"
	}
	//TODO get rid of these dumb ranges and loop like a normal human being.
	for i := range [2]int{-1, 1} {
		for j := range [2]int{-1, 1} {
			for k := 1; k < 8; k++ {
				fOff := i * k
				rOff := j * k
				if pFile+fOff < 0 || pFile+fOff > 7 || pRank+rOff < 0 || pRank+rOff > 7 ||
					strings.IndexByte(allies, s.board[pFile+fOff][pRank+rOff]) != -1 {
					break
				}
				t := s
				t.board = t.board.movePiece(pFile, pRank, pFile+fOff, pRank+rOff)
				t.ep = ""
				if s.activeColor == "w" {
					t.activeColor = "b"
				} else {
					t.activeColor = "w"
				}
				moves = append(moves, t)
			}
		}
	}
	return moves
}

func (s State) rookMoves(pFile int, pRank int) []State {
	moves := make([]State, 0, 14)
	var allies string
	if s.activeColor == "w" {
		allies = "KQRNBP"
	} else {
		allies = "kqrnbp"
	}
	//this is probably too clever for my own good.
	for card := range [2]int{0, 1} {
		for dir := range [2]int{-1, 1} {
			for i := 1; i < 8; i++ {
				fOff := (1 - card) * dir * i
				rOff := card * dir * i
				//exactly one is non-zero
				if pFile+fOff < 0 || pFile+fOff > 7 || pRank+rOff < 0 || pRank+rOff > 7 ||
					strings.IndexByte(allies, s.board[pFile+fOff][pRank+rOff]) != -1 {
					break
				}
				t := s
				t.board = t.board.movePiece(pFile, pRank, pFile+fOff, pRank+rOff)
				t.ep = ""
				if s.activeColor == "w" {
					t.activeColor = "b"
					if pRank == 7 && (pFile == 0 || pFile == 7) {
						t.castling = strings.Trim(s.castling, "KQ")
					}
				} else {
					t.activeColor = "w"
					if pRank == 0 && (pFile == 0 || pFile == 7) {
						t.castling = strings.Trim(s.castling, "kq")
					}
				}
				moves = append(moves, t)
			}
		}
	}
	return moves
}

func (s State) queenMoves(pFile int, pRank int) []State {
	return append(s.bishopMoves(pFile, pRank), s.rookMoves(pFile, pRank)...)
}

func (s State) kingMoves(pFile int, pRank int) []State {
	moves := make([]State, 0, 14)
	var allies string
	if s.activeColor == "w" {
		allies = "KQRNBP"
	} else {
		allies = "kqrnbp"
	}
	//TODO add castling
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == j || pFile+i < 0 || pFile+i > 7 || pRank+j < 0 || pRank+j > 7 ||
				strings.IndexByte(allies, s.board[pFile+i][pRank+j]) != -1 {
				break
			}

			t := s
			t.board = t.board.movePiece(pFile, pRank, pFile+i, pRank+j)
			t.ep = ""
			if s.activeColor == "w" {
				t.activeColor = "b"
				t.castling = strings.Trim(s.castling, "KQ")
			} else {
				t.activeColor = "w"
				t.castling = strings.Trim(s.castling, "kq")
			}
			moves = append(moves, t)
		}
	}
	return moves
}

func (b squares) movePiece(srcFile int, srcRank int, destFile int, destRank int) squares {
	p := b[srcFile][srcRank]
	b[srcFile][srcRank] = '.'
	b[destFile][destRank] = p
	return b
}
