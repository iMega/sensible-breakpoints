package points

import (
	"github.com/cnkei/gospline"
)

func createCubicSpline(points []*Point) gospline.Spline {
	var w, fs []float64

	for _, p := range points {
		w = append(w, float64(p.Width))
		fs = append(fs, float64(p.FileSize))
	}

	return gospline.NewCubicSpline(fs, w)
}

func calcBreakpoints(s gospline.Spline, fileSizeSrc int, budget int) []*Point {
	var (
		i int
		p []*Point

		fs = float64(fileSizeSrc)
		b  = float64(budget)
	)

	for {
		i++
		fs = fs - b
		if fs < b {
			return p
		}

		p = append(p, &Point{
			Width:    int(s.At(fs)),
			FileSize: int(fs),
		})
	}

}
