package main

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"

	"github.com/cnkei/gospline"
)

type Points struct {
	Width    float64
	FileSize float64
}

func TestSpline(t *testing.T) {
	var points = []Points{
		{Width: 162, FileSize: 46775},
		{Width: 324, FileSize: 73185},
		{Width: 648, FileSize: 165460},
		//{Width: 800, FileSize: 227262},
		{Width: 1296, FileSize: 460826},

		//{Width: 1400, FileSize: 528475},
		//{Width: 1800, FileSize: 772338},
		//{Width: 2400, FileSize: 1188984},

		{Width: 2592, FileSize: 1329808},
		//{Width: 2800, FileSize: 1507885},
		//{Width: 3000, FileSize: 1675946},
		//{Width: 3500, FileSize: 2140865},
		//{Width: 4000, FileSize: 2655643},
		//{Width: 4800, FileSize: 3588614},
		{Width: 2754, FileSize: 1472716}, // расчет по мелкой картине

		{Width: 3240, FileSize: 1886549},

		{Width: 3564, FileSize: 2203766}, // расчет по мелкой картине

		{Width: 3888, FileSize: 2532748},
		{Width: 4212, FileSize: 2890745}, // расчет по мелкой картине
		{Width: 4536, FileSize: 3273282},

		{Width: 4698, FileSize: 3464308}, // расчет по мелкой картине
		{Width: 5022, FileSize: 3876272}, // расчет по мелкой картине

		{Width: 5184, FileSize: 18302321},
	}

	x := []float64{}
	y := []float64{}

	for _, p := range points {
		x = append(x, p.Width)
		y = append(y, p.FileSize)
	}

	s := gospline.NewCubicSpline(x, y)
	for _, w := range []float64{2900, 300} {
		fmt.Printf("spline: width %f size %f\n", w, s.At(w))
	}

	//spline: width 3500.000000 size 2152999.061237  real 2140865
	//spline: width 3000.000000 size 1664140.821622  real 1675946
	//spline: width 2800.000000 size 1495250.300687  real 1507885
	//spline: width 2400.000000 size 1182459.559021  real 1188984
	//spline: width 1800.000000 size 758734.593815   real 772338
	//spline: width 1400.000000 size 517095.886749   real 528475

	//spline: width 2900.000000 size 1592438.787599  real 1591479
	// spline: width 2900.000000 size 1595466.637406 //по мелкой

	//spline: width 2900.000000 size 1504805.477399
	//spline: width 2900.000000 size 1579500.739937
	//spline: width 300.000000 size 68593.099554 67151
}

func TestOut(t *testing.T) {
	const padding = 3
	fmt.Println("sdfsdf")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, '-', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "a\tb\taligned\t")
	fmt.Fprintln(w, "aa\tbb\taligned\t")
	fmt.Fprintln(w, "aaa\tbbb\tunaligned") // no trailing tab
	fmt.Fprintln(w, "aaaa\tbbbb\taligned\t")
	w.Flush()
}
