package main

import (
	//"os"

	"fmt"
	"github.com/njm256/vezir/engine"
	"github.com/njm256/vezir/position"
)

func main() {

	//a := position.NewFen(os.Args[1])
	//fmt.Println(a)
	var g engine.GameState
	g = position.NewGame()
	calc := (*engine.MCTS(&g, 10)).(position.Game)
	fmt.Println(position.GameToFen(calc))
}
