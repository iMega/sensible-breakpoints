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

	maxWidth := img.Width
	if opt.MaxWidth > 0 {
		maxWidth = opt.MaxWidth
	}

	parts := getNumberPartsFromOriginal(opt.MinWidth, maxWidth)
	var sv []int
	if thereSufficientNumberStrongValues(parts) {
		sv = calcStrongValues(opt.MinWidth, maxWidth)
	} else {
		sv = createStrongValues(opt.MinWidth, maxWidth)
	}

	points := calcRealFileSize(img.Buffer, sv)

	s := createCubicSpline(points)

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

	size, err := bimg.NewImage(buf).Size()
	if err != nil {
		return nil, fmt.Errorf("failed getting size image, %s", err)
	}

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
