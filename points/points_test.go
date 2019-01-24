package points

import (
	"reflect"
	"testing"
)

func Test_getNumberPartsFromOriginal(t *testing.T) {
	type args struct {
		p int
		w int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "image has width 1000 and min-width 320",
			args: struct {
				p int
				w int
			}{
				p: 320,
				w: 1000,
			},
			want: 3,
		},
		{
			name: "image has width 2000 and min-width 320",
			args: struct {
				p int
				w int
			}{
				p: 320,
				w: 2000,
			},
			want: 6,
		},
		{
			name: "image has width 3000 and min-width 320",
			args: struct {
				p int
				w int
			}{
				p: 320,
				w: 3000,
			},
			want: 9,
		},
		{
			name: "image has width 3168 and min-width 320",
			args: struct {
				p int
				w int
			}{
				p: 320,
				w: 3168,
			},
			want: 9,
		},
		{
			name: "image has width 320 and min-width 320",
			args: struct {
				p int
				w int
			}{
				p: 320,
				w: 320,
			},
			want: 1,
		},
		{
			name: "image has width 200 and min-width 320",
			args: struct {
				p int
				w int
			}{
				p: 320,
				w: 200,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNumberPartsFromOriginal(tt.args.p, tt.args.w); got != tt.want {
				t.Errorf("getNumberPartsFromOriginal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_thereSufficientNumberStrongValues(t *testing.T) {
	type args struct {
		numberParts int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "7 is insufficient number of strong values",
			args: struct{ numberParts int }{numberParts: 7},
			want: false,
		},
		{
			name: "8 is sufficient number of strong values",
			args: struct{ numberParts int }{numberParts: 8},
			want: true,
		},
		{
			name: "9 is sufficient number of strong values",
			args: struct{ numberParts int }{numberParts: 9},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := thereSufficientNumberStrongValues(tt.args.numberParts); got != tt.want {
				t.Errorf("thereSufficientNumberStrongValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createStrongValues(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "create strong values between 320..500",
			args: struct {
				min int
				max int
			}{
				min: 320,
				max: 500,
			},
			want: []int{320, 342, 364, 386, 408, 430, 452, 474},
		},
		{
			name: "create strong values between 320..2500",
			args: struct {
				min int
				max int
			}{
				min: 320,
				max: 2500,
			},
			want: []int{320, 592, 864, 1136, 1408, 1680, 1952, 2224},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createStrongValues(tt.args.min, tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createStrongValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStrongValuesLessHalf(t *testing.T) {
	type args struct {
		max int
		min int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "create strong values between 320..2500",
			args: struct {
				max int
				min int
			}{
				max: 2500,
				min: 320,
			},
			want: []int{1250, 624},
		},
		{
			name: "create strong values between 320..5000",
			args: struct {
				max int
				min int
			}{
				max: 5000,
				min: 320,
			},
			want: []int{2500, 1250, 624},
		},
		{
			name: "create strong values between 0..100",
			args: struct {
				max int
				min int
			}{
				max: 100,
				min: 0,
			},
			want: []int{50, 24, 12, 6, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStrongValuesLessHalf(tt.args.max, tt.args.min); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStrongValuesLessHalf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStrongValuesOverHalf(t *testing.T) {
	type args struct {
		max  int
		min  int
		step int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "create strong values between 500..1000",
			args: struct {
				max  int
				min  int
				step int
			}{
				max:  1000,
				min:  500,
				step: 20,
			},
			want: []int{980, 940, 880, 800, 700, 580},
		},
		{
			name: "create strong values between 0..100",
			args: struct {
				max  int
				min  int
				step int
			}{
				max:  200,
				min:  100,
				step: 10,
			},
			want: []int{190, 170, 140, 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStrongValuesOverHalf(tt.args.max, tt.args.min, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStrongValuesOverHalf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMinimalWidth(t *testing.T) {
	type args struct {
		points []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: struct{ points []int }{points: []int{1, 2, 3}},
			want: 1,
		},
		{
			name: "",
			args: struct{ points []int }{points: []int{3, 2, 1}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMinimalWidth(tt.args.points); got != tt.want {
				t.Errorf("getMinimalWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcStrongValues(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "create strong values between 320..5184",
			args: struct {
				min int
				max int
			}{
				min: 320,
				max: 5184,
			},
			want: []int{2592, 1296, 648, 324, 4860, 4212, 3240},
		},
		{
			name: "create strong values between 320..10000",
			args: struct {
				min int
				max int
			}{
				min: 320,
				max: 10000,
			},
			want: []int{5000, 2500, 1250, 624, 9376, 8128, 6256},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcStrongValues(tt.args.min, tt.args.max); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcStrongValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
