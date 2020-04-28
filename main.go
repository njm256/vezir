package main

import (
	"fmt"
	"github.com/njm256/vezir/position"
	"os"
)

func main() {
	a := position.NewFen(os.Args[1])
	fmt.Println(a)
}
