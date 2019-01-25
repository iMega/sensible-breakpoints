package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/imega/sensible-breakpoints/points"
)

const (
	MIN_WIDTH_IMAGE = 320
	BUDGET          = 20000
)

var (
	writer = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
)

func main() {
	ptn, err := points.CalcBreakpoints(points.Option{
		//Filename: "times-square.jpg",
		Filename: "kettering-sky.jpg",
		//Filename: "src.jpg",
		Budget:   BUDGET,
		MinWidth: MIN_WIDTH_IMAGE,
		//MaxWidth: 2400,
		//Verbose: true,
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
