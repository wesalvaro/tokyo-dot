package main

import (
	"fmt"
	"io"
)

var circleNums = []rune{
	'⓪', '①', '②', '③', '④', '⑤', '⑥', '⑦', '⑧', '⑨', '⑩',
	'⑪', '⑫', '⑬', '⑭', '⑮', '⑯', '⑰', '⑱', '⑲', '⑳',
}

func render(g *trainGraph, r route) {
	for i := 0; i < len(r.Stations)-1; i++ {
		s := r.Stations[i]
		d := r.Stations[i+1]
		edge := g.Edge(s.ID(), d.ID()).(*edge)
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
		fmt.Printf("\n%40s ➜ %s ➜ %-40s", s, string(circleNums[int(edge.time)]), d)
	}
	fmt.Printf("\nTotal time: %.0f\n", r.Time)
}

func renderRouteHTML(w io.Writer, r route) {
	fmt.Fprintf(w, "<div>Time: %.0f</div>", r.Time)
	fmt.Fprintf(w, "<ol>")
	for i := 0; i < len(r.Stations)-1; i++ {
		f := r.Stations[i]
		t := r.Stations[i+1]
		fmt.Fprintf(w, "<li>%s -> %s</li>", f.NameEn, t.NameEn)
	}
	fmt.Fprintf(w, "</ol>")
}
