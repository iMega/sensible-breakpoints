package demo

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/h2non/bimg"
	points "github.com/iMega/sensible-breakpoints"
	"github.com/iMega/sensible-breakpoints/demo/template"
)

func MakeDemo(filename string, p []*points.Point) error {
	var (
		srcset []template.Set
		sizes  []template.Size
		widths []int
	)

	for _, w := range p {
		widths = append(widths, w.Width)
	}

	imgs, err := makeImages(filename, widths)
	if err != nil {
		return err
	}

	for _, i := range imgs {
		srcset = append(srcset, template.Set{
			Url:   i.Filename,
			Width: i.Width,
		})
	}

	sizes = append(sizes, template.Size{
		Width: 320,
		Vw:    100,
	})

	sizes = append(sizes, template.Size{
		Width: 700,
		Vw:    50,
	})

	buf := new(bytes.Buffer)
	template.GetDemo(
		srcset,
		sizes,
		template.Default{
			Src:   imgs[0].Filename,
			Width: 1000,
		},
		buf,
	)

	ioutil.WriteFile(filepath.Dir(filename)+"/index.html", buf.Bytes(), 0755)

	return nil
}

func makeImages(filename string, widths []int) ([]*points.Image, error) {
	var imgs []*points.Image

	img, err := points.ReadImage(filename)
	if err != nil {
		return nil, err
	}

	points.ProcessCalc(img, widths, func(img *points.Image, width int) (*points.Point, error) {
		resultBuffer, err := bimg.NewImage(img.Buffer).Process(bimg.Options{
			Width: width,
			Embed: true,
		})
		if err != nil {
			return nil, err
		}

		dir := filepath.Dir(img.Filename)
		base := filepath.Base(img.Filename)
		ext := filepath.Ext(img.Filename)

		postfix := strconv.Itoa(width)

		filename := filepath.Join(dir, base+"-"+postfix+ext)

		imgs = append(imgs, &points.Image{
			Width:    width,
			Filename: filename,
		})

		bimg.Write(filename, resultBuffer)

		return &points.Point{}, nil
	})

	return imgs, nil
}
