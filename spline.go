package sensbreakpoints

import (
	"sort"

	"github.com/cnkei/gospline"
)

type sortedValuesByWidth []*Point

func (sv sortedValuesByWidth) Len() int {
	return len(sv)
}

func (sv sortedValuesByWidth) Swap(i, j int) {
	sv[i], sv[j] = sv[j], sv[i]
}

func (sv sortedValuesByWidth) Less(i, j int) bool {
	return sv[i].Width < sv[j].Width
}

func createCubicSpline(points []*Point) gospline.Spline {
	var w, fs []float64

	sort.Sort(sortedValuesByWidth(points))

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

		w := int(s.At(fs))
		logger("cubic spline returns value of width is %d for size of file %d", w, int(fs))

		p = append(p, &Point{
			Width:    w,
			FileSize: int(fs),
		})
	}
}
