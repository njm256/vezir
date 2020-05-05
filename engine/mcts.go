package engine

import (
	"math/rand"
)

const (
	c = 2 //exploration parameter
)

type gameState interface {
	Result() (bool, int) //true if game over, then -1,0,1 for game result
	Moves() []gameState
}

type gameTree struct {
	root *treeNode
}

type treeNode struct {
	data     *gameState
	parent   *treeNode
	children []treeNode
}

func (t *treeNode) playout() int {
	currentState := *t.data
	for {
		if r, res := currentState.Result(); r {
			return res
		}
		nextMoves := currentState.Moves()
		i := rand.Intn(len(nextMoves))
		currentState = nextMoves[i]
	}
}
