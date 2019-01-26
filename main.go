package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/imega/sensible-breakpoints/points"
)

var (
	writer = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
)

func main() {
	ptn, err := points.CalcBreakpoints(points.Option{
		Filename:    "times-square.jpg",
		Breakpoints: []int{320, 700, 900},
		MinWidth:    320,
		Budget:      20000,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Breakpoints")
	fmt.Fprintln(writer, "#\tWidth\tSize")

	var i int
	for _, p := range ptn {
		i++
		fmt.Fprintf(writer, "%d\t%d\t%s\n", i, int(p.Width), humanize.Bytes(uint64(p.FileSize)))
	}

	writer.Flush()
}
