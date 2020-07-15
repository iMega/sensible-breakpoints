package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	points "github.com/iMega/sensible-breakpoints"
	"github.com/iMega/sensible-breakpoints/demo"
	"github.com/spf13/cobra"
)

type RootOptions struct {
	MinWidth    int
	MaxWidth    int
	Budget      int
	Verbose     bool
	Demo        bool
	Breakpoints []int
}

var (
	version     = "0.0.0"
	rootOptions = RootOptions{}
	rootCmd     = &cobra.Command{
		Use:     "sensible-breakpoints",
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE:    run,
	}
	writer = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
)

func Execute() {
	flag := rootCmd.PersistentFlags()

	flag.IntVarP(&rootOptions.MinWidth, "min-width", "m", 320, "The minimum width for which to make calculations")
	flag.IntVarP(&rootOptions.MaxWidth, "max-width", "x", 0, "The maximum width for which to make calculations")
	flag.IntVarP(&rootOptions.Budget, "budget", "b", 20000, "Set performance budget")
	flag.BoolVarP(&rootOptions.Verbose, "verbose", "v", false, "Enable verbose mode")
	flag.IntSliceVarP(&rootOptions.Breakpoints, "breakpoints", "p", []int{}, "Set your breakpoints")
	flag.BoolVarP(&rootOptions.Demo, "demo", "d", false, "Make demo")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func run(cmd *cobra.Command, args []string) error {
	ptn, err := points.CalcBreakpoints(points.Option{
		Filename:    args[0],
		Budget:      rootOptions.Budget,
		MinWidth:    rootOptions.MinWidth,
		MaxWidth:    rootOptions.MaxWidth,
		Breakpoints: rootOptions.Breakpoints,
		Verbose:     rootOptions.Verbose,
	})
	if err != nil {
		return err
	}

	fmt.Println("Breakpoints")
	fmt.Fprintln(writer, "#\tWidth\tSize")

	var i int
	for _, p := range ptn {
		i++
		fmt.Fprintf(writer, "%d\t%d\t%s\n", i, int(p.Width), humanize.Bytes(uint64(p.FileSize)))
	}

	writer.Flush()

	if rootOptions.Demo {
		demo.MakeDemo(args[0], ptn)
	}

	return nil
}
