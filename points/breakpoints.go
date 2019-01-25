package points

import (
	"fmt"
	"log"
	"sync"

	"github.com/h2non/bimg"
)

type Option struct {
	Filename string
	Budget   int
	Verbose  bool
	MinWidth int
	MaxWidth int
}

type Point struct {
	Width    int
	FileSize int
}

var verbose bool

func CalcBreakpoints(opt Option) ([]*Point, error) {
	verbose = opt.Verbose

	img, err := readImage(opt.Filename)
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
	logger("number the parts from original image, %d", parts)
	var sv []int
	if thereSufficientNumberStrongValues(parts) {
		logger("sufficient number of values for create cubic spline")
		sv = calcStrongValues(opt.MinWidth, maxWidth)
	} else {
		logger("insufficient number of values for create cubic spline")
		sv = createStrongValues(opt.MinWidth, maxWidth)
	}

	points := calcRealFileSize(img.Buffer, sv)

	s := createCubicSpline(points)
	logger("cubic spline was created")

	return calcBreakpoints(s, img.FileSize, opt.Budget), nil
}

type image struct {
	Buffer   []byte
	Width    int
	Height   int
	FileSize int
}

func readImage(filename string) (*image, error) {
	buf, err := bimg.Read(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read image, %s", err)
	}
	logger("done: read image %s", filename)

	size, err := bimg.NewImage(buf).Size()
	if err != nil {
		return nil, fmt.Errorf("failed getting size image, %s", err)
	}
	logger("done: getting size of image width: %d, height: %d", size.Width, size.Height)

	return &image{
		Buffer:   buf,
		Width:    size.Width,
		Height:   size.Height,
		FileSize: len(buf),
	}, nil
}

func calcRealFileSize(buf []byte, points []int) []*Point {
	var (
		ch  = make(chan *Point, 10)
		ret []*Point
	)

	wg := sync.WaitGroup{}
	for _, v := range points {
		wg.Add(1)
		go func(width int) {
			p, err := aspectResizeByWidth(buf, width)
			if err == nil {
				logger("calculate data for strong-values, width: %d, filesize: %d", p.Width, p.FileSize)
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

func aspectResizeByWidth(buf []byte, width int) (*Point, error) {
	resultBuffer, err := bimg.NewImage(buf).Process(bimg.Options{
		Width: width,
		Embed: true,
	})
	if err != nil {
		return nil, err
	}
	logger("done: resize image by width: %d", width)

	size, err := bimg.NewImage(resultBuffer).Size()
	if err != nil {
		return nil, err
	}
	logger("done: after resize image has dimensions width: %d, height: %d", size.Width, size.Height)

	return &Point{
		Width:    size.Width,
		FileSize: len(resultBuffer),
	}, nil
}

func logger(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}
