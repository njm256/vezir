package engine

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	c = math.Sqrt2 //exploration parameter
)

type gameState interface {
	Result() (bool, int) //true if game over, then -1,0,1 for game result
	Moves() []gameState
}

type gameTree struct {
	root *treeNode
}

type treeNode struct {
	data     gameState
	parent   *treeNode
	children []treeNode
	expanded bool
	score    int
	sims     int
}

func MCTS(g *gameState, iter int) *gameState {
	tree := gameTree{
		root: &treeNode{
			data:     *g,
			parent:   nil,
			children: []treeNode{},
			expanded: false,
			score:    0,
			sims:     0,
		},
	}
	for i := 0; i < iter; i++ {
		nextNode := tree.root.selectNode()
		v := nextNode.rollout()
		nextNode.backPropagate(v)
	}
	best := &tree.root.children[0]
	for _, v := range tree.root.children {
		if float64(v.score)/float64(v.sims) > float64(best.score)/float64(best.sims) {
			best = &v
		}
	}
	return &best.data
}

func MCTSLoop(g *gameState) {
	tree := gameTree{
		root: &treeNode{
			data:     *g,
			parent:   nil,
			children: []treeNode{},
			expanded: false,
			score:    0,
			sims:     0,
		},
	}
	for {
		nextNode := tree.root.selectNode()
		v := nextNode.rollout()
		nextNode.backPropagate(v)

		best := &tree.root.children[0]
		for _, v := range tree.root.children {
			if float64(v.score)/float64(v.sims) > float64(best.score)/float64(best.sims) {
				best = &v
			}
		}
		fmt.Println(best.data)
		fmt.Printf("score: %d, sims: %d", best.score, best.sims)
	}
}

func (t treeNode) rollout() int {
	currentState := t.data
	for {
		if r, res := currentState.Result(); r {
			return res
		}
		nextMoves := currentState.Moves()
		i := rand.Intn(len(nextMoves))
		currentState = nextMoves[i]
	}
}

func (t *treeNode) backPropagate(res int) {
	for t != nil {
		t.score += res
		t.sims++
		t = t.parent
	}
	return
}

func (t *treeNode) selectNode() *treeNode {
	if t.sims == 0 {
		return t
	}
	if !t.expanded {
		t.expand()
		return &t.children[0]
	}
	best := &t.children[0]
	for _, v := range t.children {
		if v.sims == 0 {
			return &v
		}
		if v.uct() > best.uct() {
			best = &v
		}
	}
	return best.selectNode()
}

func (t *treeNode) uct() float64 {
	return float64(t.score)/float64(t.sims) + c*math.Sqrt(math.Log(float64(t.parent.sims))/float64(t.sims))
}

func (t *treeNode) expand() {
	m := t.data.Moves()
	t.children = make([]treeNode, len(m), len(m))
	for i, v := range m {
		t.children[i] = treeNode{
			data:     v,
			parent:   t,
			children: []treeNode{},
			expanded: false,
			score:    0,
			sims:     0,
		}
	}
	t.expanded = true
}
