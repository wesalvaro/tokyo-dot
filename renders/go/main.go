package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
)

func checkErr(err error) {
	if err != nil {
		log.Panicf("%v", err)
		panic(err.Error())
	}
}

func main() {
	http.HandleFunc("/route/", handle)
	http.HandleFunc("/explore/", handle)
	http.HandleFunc("/near/", handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func loadGraph(ctx context.Context) *trainGraph {
	graphPath := os.Getenv("GRAPH_PATH")
	if graphPath != "" {
		f, err := os.Open(graphPath)
		checkErr(err)
		return readGraph(f)
	} else {
		return loadGraphFromGs(ctx)
	}
}

func loadGraphFromGs(ctx context.Context) *trainGraph {
	client, err := storage.NewClient(ctx)
	checkErr(err)
	defer client.Close()
	bucket := client.Bucket("train-lines.appspot.com")
	f, err := bucket.Object("graph.dot").NewReader(ctx)
	checkErr(err)
	defer f.Close()
	checkErr(err)
	return readGraph(f)
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html")
	parts := strings.Split(r.URL.Path[1:], "/")
	numParts := len(parts)
	log.Printf("PARTS: %s", parts)
	switch {
	case numParts == 0:
		fmt.Fprintf(w, "try a command")
	case parts[0] == "route":
		handleRoute(ctx, w, parts[1:])
	case parts[0] == "explore":
		handleExplore(ctx, w, parts[1])
	case parts[0] == "near":
		handlerNear(ctx, w, parts[1], parts[2])
	default:
		fmt.Fprintf(w, "Bad command")
	}
}

func handlerNear(ctx context.Context, w http.ResponseWriter, latPart, lngPart string) {
	lat, err := strconv.ParseFloat(latPart, 64)
	if err != nil {
		log.Printf("could not parse float lat: %s", err)
		return
	}
	lng, err := strconv.ParseFloat(lngPart, 64)
	if err != nil {
		log.Printf("could not parse float lng: %s", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	tokyo := loadGraph(ctx)
	stations := nearest(tokyo, lat, lng)
	if stations == nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "no stations found nearby")
		return
	}
	respond(w, stations)
}

func handleExplore(ctx context.Context, w http.ResponseWriter, line string) {
	w.Header().Add("Content-Type", "application/json")
	tokyo := loadGraph(ctx)
	stations := tokyo.Lines[line]
	if stations == nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "line not found")
		return
	}
	respond(w, stations)
}

func handleRoute(ctx context.Context, w http.ResponseWriter, ss []string) {
	w.Header().Add("Content-Type", "application/json")
	if len(ss) < 2 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "you need at least two stations")
		return
	}
	tokyo := loadGraph(ctx)
	log.Printf("STATIONS: %d", tokyo.Nodes().Len())
	respond(w, findMultiRoute(tokyo, ss...))
}

func respond(w http.ResponseWriter, response interface {}) {
	output, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Could not marshal response: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(output)
}

func cli() {
	f, err := os.Open("graph.dot")
	if err != nil {
		fmt.Println(err)
		return
	}
	tokyo := readGraph(f)
	fmt.Println(tokyo)
	if os.Args[1] == "route" {
		cliRoute(tokyo, os.Args[2:]...)
	} else if os.Args[1] == "explore" {
		min, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println(err)
			return
		}
		limit, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
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
