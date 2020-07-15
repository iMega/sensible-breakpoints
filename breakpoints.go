package main

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/h2non/bimg"
)

type Option struct {
	Filename    string
	Budget      int
	Verbose     bool
	MinWidth    int
	MaxWidth    int
	Breakpoints []int
}

type Point struct {
	Width    int
	FileSize int
}

var verbose bool

func CalcBreakpoints(opt Option) ([]*Point, error) {
	verbose = opt.Verbose

	img, err := ReadImage(opt.Filename)
	if err != nil {
		return nil, err
	}

	if img.FileSize <= opt.Budget {
		return nil, fmt.Errorf("image file size (%d) must be greater than option budget (%d)", img.FileSize, opt.Budget)
	}

	maxWidth := img.Width
	if opt.MaxWidth > 0 {
		logger("max-width from option %d", opt.MaxWidth)
		maxWidth = opt.MaxWidth
	}

	parts := getNumberPartsFromOriginal(opt.MinWidth, maxWidth)
	logger("number the parts from original Image, %d", parts)
	var sv []int
	if thereSufficientNumberStrongValues(parts) {
		logger("sufficient number of values for create cubic spline")
		sv = calcStrongValues(opt.MinWidth, maxWidth)
	} else {
		logger("insufficient number of values for create cubic spline")
		sv = createStrongValues(opt.MinWidth, maxWidth)
	}

	points := ProcessCalc(img, sv, aspectResizeByWidth)

	s := createCubicSpline(points)
	logger("cubic spline was created")

	result := calcBreakpoints(s, img.FileSize, opt.Budget)
	sort.Sort(sortedValuesByWidth(result))

	return filterBreakpoint(opt.Breakpoints, result), nil
}

type Image struct {
	Buffer   []byte
	Width    int
	Height   int
	FileSize int
	Filename string
}

func ReadImage(filename string) (*Image, error) {
	buf, err := bimg.Read(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read Image, %s", err)
	}
	logger("done: read Image %s", filename)

	size, err := bimg.NewImage(buf).Size()
	if err != nil {
		return nil, fmt.Errorf("failed getting size Image, %s", err)
	}
	logger("done: getting size of Image width: %d, height: %d", size.Width, size.Height)

	return &Image{
		Buffer:   buf,
		Width:    size.Width,
		Height:   size.Height,
		FileSize: len(buf),
		Filename: filename,
	}, nil
}

type imageResizeByWidth func(img *Image, width int) (*Point, error)

func ProcessCalc(img *Image, points []int, resize imageResizeByWidth) []*Point {
	var (
		ch  = make(chan *Point, 10)
		ret []*Point
	)

	wg := sync.WaitGroup{}
	for _, v := range points {
		wg.Add(1)
		go func(width int) {
			p, err := resize(img, width)
			if err == nil {
				//logger("it calculated data for strong-values, width: %d, filesize: %d", p.Width, p.FileSize)
				ch <- p
			} else {
				log.Printf("failed to resize Image, %s", err)
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

func aspectResizeByWidth(img *Image, width int) (*Point, error) {
	resultBuffer, err := bimg.NewImage(img.Buffer).Process(bimg.Options{
		Width: width,
		Embed: true,
	})
	if err != nil {
		return nil, err
	}
	logger("done: resize Image by width: %d", width)

	size, err := bimg.NewImage(resultBuffer).Size()
	if err != nil {
		return nil, err
	}
	logger("done: after resize Image has dimensions width: %d, height: %d", size.Width, size.Height)

	return &Point{
		Width:    size.Width,
		FileSize: len(resultBuffer),
	}, nil
}

func filterBreakpoint(breakpoints []int, points []*Point) []*Point {
	var result []*Point

	if len(breakpoints) < 1 {
		return points
	}

	var start int
	for _, bp := range breakpoints {
		for i := start; i < len(points); i++ {
			if bp < points[i].Width {
				result = append(result, points[i])
				start = i + 1
				break
			}
		}
	}

	return result
}

func logger(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}
