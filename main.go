package main

import (
	//"os"

	"fmt"
	"github.com/njm256/vezir/engine"
	//"github.com/njm256/vezir/position"
)

func main() {
	/*
		a := position.NewFen(os.Args[1])
		fmt.Println(a)
		g := engine.GameState(position.NewGame())
		calc := (*engine.MCTS(&g, 10)).(position.Game)
		fmt.Println(position.GameToFen(calc))
	*/
	g := engine.NewGame()
	calc := engine.MCTS(g, 1).(*engine.Game)
	fmt.Println(calc.String())
}
