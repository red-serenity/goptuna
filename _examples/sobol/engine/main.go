package main

import (
	"fmt"

	"github.com/red-serenity/goptuna/sobol"
)

func main() {
	g := sobol.NewEngine(3)
	for i := 0; i < 10; i++ {
		points := g.Draw()
		fmt.Println(points)
	}
}
