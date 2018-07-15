package main

import (
	"gonum.org/v1/gonum/graph/path"
)

func findPath(graph *trainGraph, s, d string) ([]*station, float64) {
	dest := graph.StationNode(d)
	path, _ := path.AStar(
		graph.StationNode(s),
		dest,
		graph, nil)
	route, weight := path.To(dest.ID())
	var stations []*station
	for _, s := range route {
		stations = append(stations, s.(*station))
	}
	return stations, weight
}
