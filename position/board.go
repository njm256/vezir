package position

type squares uint64

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
	hclock      int
	move        int
}
