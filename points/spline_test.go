package points

import (
	"reflect"
	"sort"
	"testing"
)

func Test_sortPointsByWidth(t *testing.T) {
	type args struct {
		points []*Point
	}
	tests := []struct {
		name string
		args args
		want []*Point
	}{
		{
			name: "sort points",
			args: struct {
				points []*Point
			}{
				points: []*Point{
					{Width: 200},
					{Width: 100},
				},
			},
			want: []*Point{
				{Width: 100},
				{Width: 200},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sort.Sort(sortedValuesByWidth(tt.args.points))
			if !reflect.DeepEqual(tt.args.points, tt.want) {
				t.Errorf("sortPointsByWidth() = %v, want %v", tt.args.points, tt.want)
			}
		})
	}
}
