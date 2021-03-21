package main

import (
	"math"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
)

type route struct {
	Stations []*station `json:"stations"`
	Time     float64    `json:"time"`
}

type stationDistance struct {
	Station  *station `json:"s"`
	Distance int      `json:"m"`
}

func makeRoute(nodes []graph.Node, time float64) route {
	var stations []*station
	for _, s := range nodes {
		stations = append(stations, s.(*station))
	}
	return route{stations, time}
}

func nearest(g *trainGraph, lat, lng float64) []*stationDistance {
	var stations []*stationDistance
	for s := range g.Stations() {
		φ1 := radians(lat)
		φ2 := radians(s.Lat)
		Δλ := radians(s.Lng - lng)
		R := 6371e3 // gives d in metres
		d := math.Acos(math.Sin(φ1)*math.Sin(φ2)+math.Cos(φ1)*math.Cos(φ2)*math.Cos(Δλ)) * R
		if d < 2e3 {
			stations = append(stations, &stationDistance{s, int(d)})
		}
	}
	return stations
}

func radians(d float64) float64 {
	return float64(d) * (math.Pi / 180.0)
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
			stations = route.Stations[1:]
		} else {
			stations = route.Stations
		}
		fullRoute.Stations = append(fullRoute.Stations, stations...)
		fullRoute.Time += route.Time
	}
	return fullRoute
}

func exploreFrom(g *trainGraph, s string, min, lim float64) map[string]route {
	shortest := path.DijkstraFrom(g.StationNode(s), g)
	routes := make(map[string]route)
	for s := range g.Stations() {
		route, time := shortest.To(s.ID())
		if time >= min && time <= lim {
			routes[s.StationID] = makeRoute(route, time)
		}
	}
	return routes
}
