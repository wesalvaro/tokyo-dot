package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/appengine"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var tokyo *trainGraph

func main() {
	f, err := os.Open("graph.dot")
	checkErr(err)
	tokyo = readGraph(f)
	http.HandleFunc("/route/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	parts := strings.Split(r.URL.String()[1:], "/")
	numParts := len(parts)
	fmt.Fprintln(w, `<!DOCTYPE html><meta charset="UTF-8"><h1>Routing</h1>`)
	if numParts == 0 {
		fmt.Fprintf(w, "try a command")
		return
	} else if parts[0] == "route" {
		if numParts < 3 {
			fmt.Fprintf(w, "you need at least two stations")
			return
		}
		route := findMultiRoute(tokyo, parts[1:]...)
		for _, d := range parts[1 : len(parts)-1] {
			fmt.Fprintf(w, "%s to ", d)
		}
		fmt.Fprint(w, parts[len(parts)-1])
		fmt.Fprintf(w, "<div>Time: %.0f</div>", route.time)
		fmt.Fprintf(w, "<ol>")
		for i := 0; i < len(route.stations)-1; i++ {
			f := route.stations[i]
			t := route.stations[i+1]
			fmt.Fprintf(w, "<li>%s -> %s</li>", f.nameEn, t.nameEn)
		}
		fmt.Fprintf(w, "</ol>")
	}
	fmt.Fprintln(w)
}

func cli() {
	f, err := os.Open("graph.dot")
	checkErr(err)
	tokyo := readGraph(f)
	fmt.Println(tokyo)
	if os.Args[1] == "route" {
		cliRoute(tokyo, os.Args[2:]...)
	} else if os.Args[1] == "explore" {
		min, err := strconv.Atoi(os.Args[3])
		checkErr(err)
		limit, err := strconv.Atoi(os.Args[4])
		checkErr(err)
		cliExplore(tokyo, os.Args[2], float64(min), float64(limit))
	}
}

func cliRoute(g *trainGraph, stations ...string) {
	fmt.Printf("Routing: %s\n", stations)
	render(g, findMultiRoute(g, stations...))
}

func cliExplore(g *trainGraph, station string, min, limit float64) {
	routes := exploreFrom(g, station, min, limit)
	for dest, route := range routes {
		fmt.Println(dest)
		render(g, route)
	}
}
