package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	f, err := os.Open("graph.dot")
	checkErr(err)
	tokyo := readGraph(f)
	fmt.Println(tokyo)
	if os.Args[1] == "route" {
		doRoute(tokyo, os.Args[2:]...)
	} else if os.Args[1] == "explore" {
		min, err := strconv.Atoi(os.Args[3])
		checkErr(err)
		limit, err := strconv.Atoi(os.Args[4])
		checkErr(err)
		doExplore(tokyo, os.Args[2], float64(min), float64(limit))
	}
}

func doRoute(g *trainGraph, stations ...string) {
	fmt.Printf("Routing: %s\n", stations)
	var fullRoute route
	for i := 0; i < len(stations)-1; i++ {
		route := findRoute(g, stations[i], stations[i+1])
		var stations []*station
		if i > 1 {
			stations = route.stations[1:]
		} else {
			stations = route.stations
		}
		fullRoute.stations = append(stations, stations...)
		fullRoute.time += route.time
	}
	render(g, fullRoute)
}

func doExplore(g *trainGraph, station string, min, limit float64) {
	fmt.Printf("Exploring from: %s\n", station)
	routes := exploreFrom(g, station, min, limit)
	for dest, route := range routes {
		fmt.Println(dest)
		render(g, route)
	}
}
