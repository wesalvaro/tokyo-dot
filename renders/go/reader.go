package main

import (
	"fmt"
	"log"
	"os"

	"gonum.org/v1/gonum/graph/formats/dot"
	"gonum.org/v1/gonum/graph/formats/dot/ast"
)

func main() {
	f, err := os.Open("graph.dot")
	if err != nil {
		log.Fatal(err)
	}
	g, err := dot.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	for _, s := range g.Graphs[0].Stmts {
		if sg, ok := s.(*ast.Subgraph); ok {
			fmt.Println(sg.ID)
			if sg.ID == "HibiyaStations" {
				fmt.Println(sg.String())
			}
		}
	}
}
