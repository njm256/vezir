package main

import (
	//"os"

	"github.com/njm256/vezir/engine"
	"github.com/njm256/vezir/position"
)

func main() {

	//a := position.NewFen(os.Args[1])
	//fmt.Println(a)
	var g engine.GameState
	g = position.NewGame()
	engine.MCTSLoop(&g)
}
