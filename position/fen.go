package position

import (
	"strconv"
	"strings"
)

type Fen struct {
	pos         string
	activeColor string
	castling    string
	ep          string
	hClock      int
	move        int
}

func NewFen(str string) *Fen {
	s := strings.Split(str, " ")
	newPos := new(Fen)
	newPos.pos = s[0]
	newPos.activeColor = s[1]
	newPos.castling = s[2]
	newPos.ep = s[3]
	newPos.hClock, _ = strconv.Atoi(s[4])
	newPos.move, _ = strconv.Atoi(s[5])
	return newPos
}

func (f Fen) String() string {
	ranks := strings.Split(f.pos, "/")
	r := strings.NewReplacer(
		"1", " ",
		"2", "  ",
		"3", "   ",
		"4", "    ",
		"5", "     ",
		"6", "      ",
		"7", "       ",
		"8", "        ",
	)
	var newRanks [8]string
	for i := 0; i < 8; i++ {
		newRanks[i] = r.Replace(ranks[i])
	}
	return strings.Join(newRanks[:], "\n")

}

func Wtf() {
	println("why")
}
