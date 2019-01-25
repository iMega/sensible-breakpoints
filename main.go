package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	points2 "github.com/imega/sensible-breakpoints/points"
)

const (
	MIN_WIDTH_IMAGE = 320
	BUDGET          = 20000
)

var (
	writer = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
)

type Point struct {
	Width    float64
	FileSize float64
}

func main() {
	ptn, err := points2.CalcBreakpoints(points2.Option{
		Filename: "times-square.jpg",
		//Filename: "src.jpg",
		Budget:   BUDGET,
		MinWidth: MIN_WIDTH_IMAGE,
		//MaxWidth: 1200,
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
