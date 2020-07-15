package main

import (
	"reflect"
	"testing"
)

func Test_filterBreakpoint(t *testing.T) {
	type args struct {
		breakpoints []int
		points      []*Point
	}

	points := []*Point{
		{Width: 300},
		{Width: 350},
		{Width: 450},
		{Width: 500},
		{Width: 550},
		{Width: 650},
		{Width: 700},
		{Width: 750},
		{Width: 850},
		{Width: 900},
	}

	tests := []struct {
		name string
		args args
		want []*Point
	}{
		{
			name: "breakpoints is empty, returns points",
			args: args{
				breakpoints: []int{},
				points:      points,
			},
			want: points,
		},
		{
			name: "filtering points by breakpoints",
			args: args{
				breakpoints: []int{400, 600, 800},
				points:      points,
			},
			want: []*Point{
				points[2],
				points[5],
				points[8],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterBreakpoint(tt.args.breakpoints, tt.args.points); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterBreakpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
