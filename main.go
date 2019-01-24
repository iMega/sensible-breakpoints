package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"text/tabwriter"

	"github.com/cnkei/gospline"
	"github.com/dustin/go-humanize"
	"github.com/h2non/bimg"
	points2 "github.com/imega/image-resize/points"
)

const (
	MIN_WIDTH_IMAGE = 320
	BUDGET          = 200000
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

	buf, err := bimg.Read("src.jpg")
	if err != nil {
		log.Fatalf("failed to read image, %s", err)
	}

	srcSize, err := bimg.NewImage(buf).Size()
	if err != nil {
		log.Fatalf("failed getting size image, %s", err)
	}

	width := srcSize.Width
	fileSize := len(buf)

	var breakpoints []int

	breakpoints = append(breakpoints, getBreakpoints(width, MIN_WIDTH_IMAGE)...)
	points := calcRealFileSize(buf, breakpoints)

	minWidth := getMinimalWidth(points)
	breakpoints = breakpoints[:0]
	breakpoints = append(breakpoints, getStrongValuesOverHalf(width, width/2, minWidth)...)

	points = append(points, calcRealFileSize(buf, breakpoints)...)

	for _, v := range points {
		fmt.Printf("width: %f, size: %f\n", v.Width, v.FileSize)
	}

	s := createCubicSpline(points)

	fmt.Println("Breakpoints")
	fmt.Fprintln(writer, "#\tWidth\tSize")

	bp := calcBreakpoints(s, fileSize, BUDGET)
	for _, p := range bp {
		i++
		fmt.Fprintf(writer, "%d\t%d\t%s\n", i, int(p.Width), humanize.Bytes(uint64(p.FileSize)))
	}

	writer.Flush()
}

func getMinimalWidth(points []*Point) int {
	min := points[0].Width

	for _, p := range points {
		if p.Width < min {
			min = p.Width
		}
	}

	return int(min)
}

func calcRealFileSize(buf []byte, breakpoints []int) []*Point {
	var (
		ch  = make(chan *Point, 10)
		ret []*Point
	)

	wg := sync.WaitGroup{}
	for _, v := range breakpoints {
		wg.Add(1)
		go func(width int) {
			p, err := aspectResizeByWith(buf, width)
			if err == nil {
				ch <- p
			} else {
				log.Printf("failed to resize image, %s", err)
			}
			wg.Done()
		}(v)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for {
		select {
		case p, ok := <-ch:
			if !ok {
				return ret
			}
			ret = append(ret, p)
		}
	}
}

func getBreakpoints(max, min int) []int {
	var (
		ret         []int
		strongValue = max
	)

	for {
		strongValue = strongValue / 2
		if strongValue < min {
			return ret
		}

		if isEven(strongValue) == false {
			strongValue--
		}

		ret = append(ret, strongValue)
	}
}

func getStrongValuesOverHalf(max, min, step int) []int {
	var (
		ret         []int
		strongValue = max
		i           int
	)

	for {
		i++

		strongValue = strongValue - step*i

		if strongValue < min {
			return ret
		}

		if isEven(strongValue) == false {
			strongValue--
		}

		ret = append(ret, strongValue)
	}
}

func isEven(n int) bool {
	return n%2 == 0
}

func aspectResizeByWith(buf []byte, width int) (*Point, error) {
	resultBuffer, err := bimg.NewImage(buf).Process(bimg.Options{
		Width: width,
		Embed: true,
	})
	if err != nil {
		return nil, err
	}

	size, err := bimg.NewImage(resultBuffer).Size()
	if err != nil {
		return nil, err
	}

	//bimg.Write(fmt.Sprintf("img-%d.jpg", size.Width), resultBuffer)

	return &Point{
		Width:    float64(size.Width),
		FileSize: float64(len(resultBuffer)),
	}, nil
}

func createCubicSpline(points []*Point) gospline.Spline {
	var w, fs []float64

	for _, p := range points {
		w = append(w, p.Width)
		fs = append(fs, p.FileSize)
	}

	fmt.Println(w)
	fmt.Println(fs)

	return gospline.NewCubicSpline(fs, w)
}

func calcBreakpoints(s gospline.Spline, fileSizeSrc int, budget int) []*Point {
	var (
		i  int
		fs float64
		b  = float64(budget)
		p  []*Point
	)

	fs = float64(fileSizeSrc)

	for {
		i++
		fs = fs - b
		if fs < b {
			return p
		}

		p = append(p, &Point{
			Width:    s.At(fs),
			FileSize: fs,
		})
	}

}
