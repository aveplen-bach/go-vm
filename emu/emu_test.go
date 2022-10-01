package emu

import (
	"testing"
)

func TestCpu_push(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		c    *Cpu
		args args
		want *Cpu
	}{
		{
			name: "push should add value to a stack",
			args: args{1},
			c: &Cpu{
				sp:    0,
				stack: [STACK_LIMIT]int32{},
			},
			want: &Cpu{
				sp:    1,
				stack: [STACK_LIMIT]int32{1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.push(tt.args.n)

			if tt.c.sp != tt.want.sp {
				t.Errorf("sp wrong value: %d, expected: %d", tt.c.sp, tt.want.sp)
			}

			for i := 0; i < STACK_LIMIT; i++ {
				if tt.c.stack[i] != tt.want.stack[i] {
					t.Errorf(
						"wrong stack value on index %d: %d, expected: %d",
						i, tt.c.stack[i], tt.want.stack[i],
					)
				}
			}
		})
	}
}

func TestCpu_pop(t *testing.T) {
	tests := []struct {
		name  string
		c     *Cpu
		want  *Cpu
		want1 int32
	}{
		{
			name: "pop should return top value",
			c: &Cpu{
				sp:    1,
				stack: [STACK_LIMIT]int32{1},
			},
			want: &Cpu{
				sp:    0,
				stack: [STACK_LIMIT]int32{1},
			},
			want1: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.pop(); got != tt.want1 {
				t.Errorf("Cpu.pop() = %v, want %v", got, tt.want1)
			}

			if tt.c.sp != tt.want.sp {
				t.Errorf("sp wrong value: %d, expected: %d", tt.c.sp, tt.want.sp)
			}

			for i := 0; i < STACK_LIMIT; i++ {
				if tt.c.stack[i] != tt.want.stack[i] {
					t.Errorf(
						"wrong stack value on index %d: %d, expected: %d",
						i, tt.c.stack[i], tt.want.stack[i],
					)
				}
			}
		})
	}
}

func TestCpu_iadd(t *testing.T) {
	tests := []struct {
		name string
		c    *Cpu
		want *Cpu
	}{
		{
			name: "iadd should pop two elements and push their sum",
			c: &Cpu{
				sp:    2,
				stack: [STACK_LIMIT]int32{1, 2},
			},
			want: &Cpu{
				sp:    1,
				stack: [STACK_LIMIT]int32{3, 2}, // leaving waste in stack
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.iadd()

			if tt.c.sp != tt.want.sp {
				t.Errorf("sp wrong value: %d, expected: %d", tt.c.sp, tt.want.sp)
			}

			for i := 0; i < STACK_LIMIT; i++ {
				if tt.c.stack[i] != tt.want.stack[i] {
					t.Errorf(
						"wrong stack value on index %d: %d, expected: %d",
						i, tt.c.stack[i], tt.want.stack[i],
					)
				}
			}
		})
	}
}

func TestCpu_isub(t *testing.T) {
	tests := []struct {
		name string
		c    *Cpu
		want *Cpu
	}{
		{
			name: "sub should pop two elements from stack and push their difference",
			c: &Cpu{
				sp:    2,
				stack: [STACK_LIMIT]int32{2, 1},
			},
			want: &Cpu{
				sp:    1,
				stack: [STACK_LIMIT]int32{1, 1}, // leaving waste in stack
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.isub()

			if tt.c.sp != tt.want.sp {
				t.Errorf("sp wrong value: %d, expected: %d", tt.c.sp, tt.want.sp)
			}

			for i := 0; i < STACK_LIMIT; i++ {
				if tt.c.stack[i] != tt.want.stack[i] {
					t.Errorf(
						"wrong stack value on index %d: %d, expected: %d",
						i, tt.c.stack[i], tt.want.stack[i],
					)
				}
			}
		})
	}
}

func TestCpu_iand(t *testing.T) {
	tests := []struct {
		name string
		c    *Cpu
		want *Cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise and",
			c: &Cpu{
				sp:    2,
				stack: [STACK_LIMIT]int32{7, 5},
			},
			want: &Cpu{
				sp:    1,
				stack: [STACK_LIMIT]int32{5, 5}, // leaving waste in stack
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.iand()

			if tt.c.sp != tt.want.sp {
				t.Errorf("sp wrong value: %d, expected: %d", tt.c.sp, tt.want.sp)
			}

			for i := 0; i < STACK_LIMIT; i++ {
				if tt.c.stack[i] != tt.want.stack[i] {
					t.Errorf(
						"wrong stack value on index %d: %d, expected: %d",
						i, tt.c.stack[i], tt.want.stack[i],
					)
				}
			}
		})
	}
}

func TestCpu_ior(t *testing.T) {
	tests := []struct {
		name string
		c    *Cpu
		want *Cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise and",
			c: &Cpu{
				sp:    2,
				stack: [STACK_LIMIT]int32{7, 5},
			},
			want: &Cpu{
				sp:    1,
				stack: [STACK_LIMIT]int32{7, 5}, // leaving waste in stack
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ior()

			if tt.c.sp != tt.want.sp {
				t.Errorf("sp wrong value: %d, expected: %d", tt.c.sp, tt.want.sp)
			}

			for i := 0; i < STACK_LIMIT; i++ {
				if tt.c.stack[i] != tt.want.stack[i] {
					t.Errorf(
						"wrong stack value on index %d: %d, expected: %d",
						i, tt.c.stack[i], tt.want.stack[i],
					)
				}
			}
		})
	}
}

func TestCpu_ixor(t *testing.T) {
	tests := []struct {
		name string
		c    *Cpu
		want *Cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise xor",
			c: &Cpu{
				sp:    2,
				stack: [STACK_LIMIT]int32{7, 5},
			},
			want: &Cpu{
				sp:    1,
				stack: [STACK_LIMIT]int32{4, 5}, // leaving waste in stack
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ixor()

			if tt.c.sp != tt.want.sp {
				t.Errorf("sp wrong value: %d, expected: %d", tt.c.sp, tt.want.sp)
			}

			for i := 0; i < STACK_LIMIT; i++ {
				if tt.c.stack[i] != tt.want.stack[i] {
					t.Errorf(
						"wrong stack value on index %d: %d, expected: %d",
						i, tt.c.stack[i], tt.want.stack[i],
					)
				}
			}
		})
	}
}
