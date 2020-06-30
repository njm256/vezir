package position

import (
	"fmt"
	"strconv"
	"strings"
)

//TODO find a better one of these. I think I saw one in sunfish?
var chars = map[byte]rune{
	'R': '♜',
	'N': '♞',
	'B': '♗',
	'Q': '♕',
	'K': '♔',
	'P': '♙',
	'r': '♜',
	'n': '♞',
	'b': '♝',
	'q': '♛',
	'k': '♚',
	'p': '♙',
}

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
	return f.simpleString()
}

func (f Fen) simpleString() string {
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

func (f Fen) colorString() string {
	const (
		ww  = "\033[97;47m"
		wb  = "\033[97;100m"
		bw  = "\033[30;47m"
		bb  = "\033[30;100m"
		end = "\033[0m"
	)
	s := strings.Split(f.simpleString(), "\n")
	var sb strings.Builder
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if (i+j)%2 == 0 {
				if strings.IndexByte("KQRBNP", s[i][j]) != -1 {
					fmt.Fprintf(&sb, "%s%c ", ww, chars[s[i][j]])
				} else if strings.IndexByte("kqrbnp", s[i][j]) != -1 {
					fmt.Fprintf(&sb, "%s%c ", bw, chars[s[i][j]])
				} else {
					fmt.Fprintf(&sb, "%s  ", ww)
				}
			} else {
				if strings.IndexByte("KQRBNP", s[i][j]) != -1 {
					fmt.Fprintf(&sb, "%s%c ", wb, chars[s[i][j]])
				} else if strings.IndexByte("kqrbnp", s[i][j]) != -1 {
					fmt.Fprintf(&sb, "%s%c ", bb, chars[s[i][j]])
				} else {
					fmt.Fprintf(&sb, "%s  ", bb)
				}
			}
		}
		fmt.Fprintf(&sb, "%s\n", end)
	}
	return sb.String()

}

func Wtf() {
	println("why")
}
