package position

import (
	"fmt"
	"strings"
)

type squares [8][8]byte

var startBoard = [8][8]byte{
	[8]byte{'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r'},
	[8]byte{'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p'},
	[8]byte{'.', '.', '.', '.', '.', '.', '.', '.'},
	[8]byte{'.', '.', '.', '.', '.', '.', '.', '.'},
	[8]byte{'.', '.', '.', '.', '.', '.', '.', '.'},
	[8]byte{'.', '.', '.', '.', '.', '.', '.', '.'},
	[8]byte{'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P'},
	[8]byte{'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R'},
}

type State struct {
	board          squares
	activeColor    string
	castling       string
	ep             string
	pMoveOrCapture bool
}

var startState = State{
	board:          startBoard,
	activeColor:    "w",
	castling:       "KQkq",
	ep:             "",
	pMoveOrCapture: false,
}

type Game struct {
	state  State
	hclock int
	moveNo int
	prev   map[string]int // maybe map[State]int? seems better to hash the states
}

func NewGame(startPos ...State) (g Game) {
	g = Game{}
	g.hclock = 0
	g.moveNo = 0
	g.prev = make(map[string]int)
	if len(startPos) > 0 {
		g.state = startPos[0]
	} else {
		g.state = startState
	}
	return
}

func (g Game) hashCode() string {
	return GameToFen(g).String()
}

func GameToFen(s Game) Fen {
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
	return Fen{
		pos:         str.String(),
		activeColor: s.state.activeColor,
		castling:    s.state.castling,
		ep:          s.state.ep,
		hClock:      s.hclock,
		move:        s.moveNo,
	}
}

func (g *Game) move(s State) {
	h := g.hashCode()
	g.prev[h] = g.prev[h] + 1
	g.state = s
	if s.activeColor == "w" {
		g.moveNo++
	}
	if s.pMoveOrCapture {
		g.hclock = 0
	} else {
		g.hclock++
	}
	//TODO put something here -- what was I supposed to put here?
}

func (g Game) IMove(s State) Game {
	//??TODO refactor by making a copy instead
	newGame := g
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
	if s.pMoveOrCapture {
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

func (s State) pawnMoves(pRank int, pFile int) []State {
	moves := make([]State, 0, 12)
	var allies string
	var enemies string
	var dir int         //up (-1) or down (+1) the board
	var sRank int       //starting rank of a pawn
	var promoteRank int //rank at which promotions must be handled
	if s.activeColor == "w" {
		allies = "KQRNBP"
		enemies = "kqrnbp"
		dir = -1
		sRank = 6
		promoteRank = 1

	} else {
		allies = "kqrnbp"
		enemies = "KQRNBP"
		dir = 1
		sRank = 1
		promoteRank = 6
	}
	//have to handle regular moves, captures, first move, en passant, and promotion

	//promotion
	if pRank == promoteRank {
		//TODO this nightmare.
		//return here, or call out to a different function.
		return moves
	}
	//captures
	for i := -1; i < 2; i += 2 {
		if pFile+i < 0 || pFile+i > 7 {
			continue
		}
		if strings.IndexByte(enemies, s.board[pRank+dir][pFile+i]) != -1 {
			t := s
			t.board = s.board.movePiece(pRank, pFile, pRank+dir, pFile+i)
			if s.activeColor == "w" {
				t.activeColor = "b"
			} else {
				t.activeColor = "w"
			}
			t.ep = ""
			t.pMoveOrCapture = true
			//TODO edit castling
			moves = append(moves, t)
		}
	}
	//regular move
	if strings.IndexByte(allies, s.board[pFile+dir][pRank]) == -1 &&
		strings.IndexByte(enemies, s.board[pFile+dir][pRank]) == -1 {
		t := s
		t.board = s.board.movePiece(pFile, pRank, pFile+dir, pRank)
		if s.activeColor == "w" {
			t.activeColor = "b"
		} else {
			t.activeColor = "w"
		}
		t.ep = ""
		t.pMoveOrCapture = true
		moves = append(moves, t)
		//starting move
		if pRank == sRank && strings.IndexByte(allies, s.board[pRank+dir*2][pFile]) == -1 &&
			strings.IndexByte(enemies, s.board[pRank+dir*2][pFile]) == -1 {
			t := s
			t.board = s.board.movePiece(pRank, pFile, pRank+dir*2, pFile)
			if s.activeColor == "w" {
				t.activeColor = "b"
			} else {
				t.activeColor = "w"
			}
			t.pMoveOrCapture = true
			moves = append(moves, t)
			//TODO add ep for this case
		}
	}

	return moves
}
func (s State) knightMoves(pRank int, pFile int) []State {
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
				strings.IndexByte(allies, s.board[pRank+i][pFile+j]) != -1 {
				continue
			}
			t := s
			t.board = t.board.movePiece(pRank, pFile, pRank+i, pFile+j)
			t.ep = ""
			if s.activeColor == "w" {
				t.activeColor = "b"
			} else {
				t.activeColor = "w"
			}
			t.pMoveOrCapture = false
			moves = append(moves, t)
		}
	}
	return moves
}

func (s State) bishopMoves(pRank int, pFile int) []State {
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
					strings.IndexByte(allies, s.board[pRank+rOff][pFile+fOff]) != -1 {
					break
				}
				t := s
				t.board = t.board.movePiece(pRank, pFile, pRank+rOff, pFile+fOff)
				t.ep = ""
				if s.activeColor == "w" {
					t.activeColor = "b"
				} else {
					t.activeColor = "w"
				}
				//TODO check for capture
				t.pMoveOrCapture = false
				moves = append(moves, t)
			}
		}
	}
	return moves
}

func (s State) rookMoves(pRank int, pFile int) []State {
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
				rOff := (1 - card) * dir * i
				fOff := card * dir * i
				//exactly one is non-zero
				if pFile+fOff < 0 || pFile+fOff > 7 || pRank+rOff < 0 || pRank+rOff > 7 ||
					strings.IndexByte(allies, s.board[pRank+rOff][pFile+fOff]) != -1 {
					break
				}
				t := s
				t.board = t.board.movePiece(pRank, pFile, pRank+rOff, pFile+fOff)
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
				//TODO check for capture
				t.pMoveOrCapture = false
				moves = append(moves, t)
			}
		}
	}
	return moves
}

func (s State) queenMoves(pRank int, pFile int) []State {
	return append(s.bishopMoves(pRank, pFile), s.rookMoves(pRank, pFile)...)
}

func (s State) kingMoves(pRank int, pFile int) []State {
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
			t.board = t.board.movePiece(pRank, pFile, pRank+i, pFile+j)
			t.ep = ""
			if s.activeColor == "w" {
				t.activeColor = "b"
				t.castling = strings.Trim(s.castling, "KQ")
			} else {
				t.activeColor = "w"
				t.castling = strings.Trim(s.castling, "kq")
			}
			t.pMoveOrCapture = false
			moves = append(moves, t)
		}
	}
	return moves
}

func (b squares) movePiece(srcRank int, srcFile int, destRank int, destFile int) squares {
	p := b[srcRank][srcFile]
	b[srcRank][srcFile] = '.'
	b[destRank][destFile] = p
	return b
}

func (g *Game) Result() (bool, int) { //return true if game over, and -1/0/1 for result
	if g.hclock > 50 { //50 move rule
		return true, 0
	}
	if g.prev[g.hashCode()] == 2 { //3-fold rep
		return true, 0
	}
	var wKing bool
	var bKing bool
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if g.state.board[i][j] == 'K' {
				wKing = true
			} else if g.state.board[i][j] == 'k' {
				bKing = true
			}
		}
	}
	if !wKing {
		return true, -1
	} else if !bKing {
		return true, 1
	}
	if len(g.state.Moves()) == 0 {
		return true, 0
	}
	return false, 0
}
