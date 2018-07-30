package main

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
)

type route struct {
	stations []*station
	time     float64
}

func makeRoute(nodes []graph.Node, time float64) route {
	var stations []*station
	for _, s := range nodes {
		stations = append(stations, s.(*station))
	}
	return route{stations, time}
}

func findRoute(g *trainGraph, s, d string) route {
	dest := g.StationNode(d)
	shortest, _ := path.AStar(
		g.StationNode(s),
		dest,
		g, nil)
	return makeRoute(shortest.To(dest.ID()))
}

func findMultiRoute(g *trainGraph, stations ...string) route {
	var fullRoute route
	for i := 0; i < len(stations)-1; i++ {
		route := findRoute(g, stations[i], stations[i+1])
		var stations []*station
		if i > 0 {
			stations = route.stations[1:]
		} else {
			stations = route.stations
		}
		fullRoute.stations = append(fullRoute.stations, stations...)
		fullRoute.time += route.time
	}
	return fullRoute
}

func exploreFrom(g *trainGraph, s string, min, lim float64) map[string]route {
	shortest := path.DijkstraFrom(g.StationNode(s), g)
	routes := make(map[string]route)
	for s := range g.Stations() {
		route, time := shortest.To(s.ID())
		if time >= min && time <= lim {
			routes[s.StationID()] = makeRoute(route, time)
		}
	}
	return routes
}
