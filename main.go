package main

import (
	"fmt"
	"os"

	"github.com/njm256/vezir/position"
)

func main() {
	a := position.NewFen(os.Args[1])
	fmt.Println(a)
}
