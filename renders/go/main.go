package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("graph.dot")
	if err != nil {
		log.Fatal(err)
	}
	g := readGraph(f)
	fmt.Println(g)
	fmt.Println(findPath(g, "TY07", "H04"))
}
