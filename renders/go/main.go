package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/appengine/log"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
)

func checkErr(ctx context.Context, err error) {
	if err != nil {
		log.Errorf(ctx, "%v", err)
		panic(err.Error())
	}
}

func main() {
	http.HandleFunc("/route/", handle)
	appengine.Main()
}

func loadGraphFromGs(ctx context.Context) *trainGraph {
	// `dev_appserver.py --default_gcs_bucket_name GCS_BUCKET_NAME`
	bucketName, err := file.DefaultBucketName(ctx)
	log.Debugf(ctx, "BUCKET: %s", bucketName)
	checkErr(ctx, err)
	client, err := storage.NewClient(ctx)
	checkErr(ctx, err)
	defer client.Close()
	bucket := client.Bucket(bucketName)
	f, err := bucket.Object("graph.dot").NewReader(ctx)
	checkErr(ctx, err)
	defer f.Close()
	checkErr(ctx, err)
	return readGraph(f)
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	w.Header().Set("Content-Type", "text/html")
	parts := strings.Split(r.URL.Path[1:], "/")
	numParts := len(parts)
	log.Debugf(ctx, "PARTS: %s", parts)
	if numParts == 0 {
		fmt.Fprintf(w, "try a command")
		return
	} else if parts[0] == "route" {
		if numParts < 3 {
			fmt.Fprintf(w, "you need at least two stations")
			return
		}
		tokyo := loadGraphFromGs(ctx)
		log.Debugf(ctx, "STATIONS: %d", len(tokyo.Nodes()))
		route := findMultiRoute(tokyo, parts[1:]...)
		w.Header().Add("Content-Type", "application/json")
		output, err := json.Marshal(route)
		if err != nil {
			log.Criticalf(ctx, "Could not marshal route: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(output)
	} else {
		fmt.Fprintf(w, "Bad command")
	}
	fmt.Fprintln(w)
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
