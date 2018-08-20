package main

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	edot "gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/formats/dot"
	"gonum.org/v1/gonum/graph/formats/dot/ast"
	"gonum.org/v1/gonum/graph/simple"
)

func newTrainGraph() *trainGraph {
	return &trainGraph{
		DirectedGraph: simple.NewDirectedGraph(),
		Lines:         make(map[string][]*station),
	}
}

type trainGraph struct {
	*simple.DirectedGraph
	Lines map[string][]*station
}

func (g *trainGraph) StationNode(sid string) *station {
	for _, n := range g.DirectedGraph.Nodes() {
		s := n.(*station)
		if sid == s.StationID {
			return s
		}
	}
	return nil
}

func (g *trainGraph) Stations() chan *station {
	ch := make(chan *station)
	go func() {
		for _, n := range g.Nodes() {
			ch <- n.(*station)
		}
		close(ch)
	}()
	return ch
}

func (g *trainGraph) Weight(s, d int64) (float64, bool) {
	if !g.HasEdgeBetween(s, d) {
		return 0, false
	}
	e := g.Edge(s, d).(*edge)
	return e.time, true
}

func (g *trainGraph) NewNode() graph.Node {
	return &station{Node: g.DirectedGraph.NewNode(), g: g}
}

func (g *trainGraph) NewEdge(from, to graph.Node) graph.Edge {
	return &edge{Edge: g.DirectedGraph.NewEdge(from, to)}
}

func (g *trainGraph) String() string {
	return fmt.Sprintf("%d platforms", len(g.Nodes()))
}

var quotedLabel = regexp.MustCompile(`"(?P<content>.*?)"`)

type edge struct {
	graph.Edge
	time   float64
	oneWay bool
	cars   []string
}

func (e *edge) SetAttribute(attr encoding.Attribute) error {
	if attr.Key == "len" {
		time, err := strconv.ParseFloat(attr.Value, 32)
		if err != nil {
			return err
		}
		e.time = time
	} else if attr.Key == "return" && attr.Value == "no" {
		e.oneWay = true
	} else if attr.Key == "label" {
		value := quotedLabel.ReplaceAllString(attr.Value, "${content}")
		if value != "-" && value != "車内" {
			transfer := strings.SplitN(value, "|", 2)
			dir := strings.SplitN(transfer[0], ";", 2)
			if len(dir) == 2 {
				log.Println("transfer car by direction")
			}
			e.cars = strings.Split(dir[0], ",")
		}
	}
	return nil
}

type station struct {
	graph.Node `json:"-"`
	StationID  string `json:"i"`
	NameEn     string `json:"e"`
	NameJp     string `json:"n"`
	Line       string `json:"b"`
	Lat, Lng   float64
	g          *trainGraph
}

func (s *station) String() string {
	return s.StationID + " " + s.NameEn
}

func (s *station) SetDOTID(id string) {
	s.StationID = id
}

func (s *station) SetAttribute(attr encoding.Attribute) error {
	switch attr.Key {
	case "label":
		re := regexp.MustCompile(`"{(?P<jp>.*?)\|(?P<en>.*?)}\|{(?P<line>.*?)\|.*?}"`)
		s.NameEn = re.ReplaceAllString(attr.Value, "${en}")
		s.NameJp = re.ReplaceAllString(attr.Value, "${jp}")
		s.Line = re.ReplaceAllString(attr.Value, "${line}")
		s.g.Lines[s.Line] = append(s.g.Lines[s.Line], s)
	case "pos":
		latLng := strings.Split(strings.Trim(attr.Value, "!\""), ",")
		lat, err := strconv.ParseFloat(latLng[0], 64)
		if err != nil {
			return err
		}
		lng, err := strconv.ParseFloat(latLng[1], 64)
		if err != nil {
			return err
		}
		s.Lat, s.Lng = lat, lng
	}
	return nil
}

func convertToTrainGraph(graph *ast.Graph) *trainGraph {
	dst := newTrainGraph()
	if err := edot.Unmarshal([]byte(graph.String()), dst); err != nil {
		log.Fatal(err)
	}
	for _, e := range dst.Edges() {
		// Make edges two-way unless specified individually:
		if f := e.(*edge); !f.oneWay {
			// Return edge already exists in graph:
			if dst.HasEdgeFromTo(f.To().ID(), f.From().ID()) {
				continue
			}
			r := dst.NewEdge(f.To(), f.From()).(*edge)
			r.time = f.time
			dst.SetEdge(r)
		}
	}
	return dst
}

func readGraph(f io.Reader) *trainGraph {
	g, err := dot.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	return convertToTrainGraph(g.Graphs[0])
}
