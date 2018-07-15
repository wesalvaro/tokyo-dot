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
	tokyo := readGraph(f)
	fmt.Println(tokyo)
	var fullRoute []*station
	var totalTime float64
	for i := 1; i < len(os.Args)-1; i++ {
		route, time := findPath(tokyo, os.Args[i], os.Args[i+1])
		if i > 1 {
			route = route[1:]
		}
		fullRoute = append(fullRoute, route...)
		totalTime += time
	}
	render(tokyo, fullRoute, totalTime)
}
