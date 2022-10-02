package example

import "testing"

func TestArraySum(t *testing.T) {
	type args struct {
		arr []uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "sum([]) = 0",
			args: args{
				arr: []uint32{},
			},
			want: 0,
		},
		{
			name: "sum(3, 5, 4, 7) = 19",
			args: args{
				arr: []uint32{3, 5, 4, 7},
			},
			want: 19,
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
