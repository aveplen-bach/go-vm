package internal

import "testing"

func TestArraySum(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "sum(3, 5, 4, 7) = 19",
			args: args{
				arr: []int{3, 5, 4, 7},
			},
			want: 19,
		},
		{
			name: "sum([]) = 0",
			args: args{
				arr: []int{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArraySum(tt.args.arr); got != tt.want {
				t.Errorf("ArraySum() = %v, want %v", got, tt.want)
			}
		})
	}
}
