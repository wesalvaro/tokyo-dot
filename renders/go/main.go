package main

import (
	"fmt"
	"log"
	"os"
)

var circleNums = []rune{
	'⓪',
	'①',
	'②',
	'③',
	'④',
	'⑤',
	'⑥',
	'⑦',
	'⑧',
	'⑨',
	'⑩',
	'⑪',
	'⑫',
	'⑬',
	'⑭',
	'⑮',
	'⑯',
	'⑰',
	'⑱',
	'⑲',
	'⑳',
}

func main() {
	f, err := os.Open("graph.dot")
	if err != nil {
		log.Fatal(err)
	}
	tokyo := readGraph(f)
	fmt.Println(tokyo)
	route, time := findPath(tokyo, os.Args[1], os.Args[2])
	for i := 0; i < len(route)-1; i++ {
		s := route[i]
		d := route[i+1]
		edge := tokyo.Edge(s.ID(), d.ID()).(*edge)
		for i, c := range edge.cars {
			comma := ""
			if i > 0 {
				comma = ", "
			}
			fmt.Printf("%s%s", comma, c)
		}
		if len(edge.cars) > 0 {
			fmt.Print("号車")
		}
		fmt.Printf("\n%5s ➜ %s ➜ %-5s", s, string(circleNums[int(edge.time)]), d)
	}
	fmt.Printf("\nTotal time: %.0f\n", time)
}
