package internal

import (
	"fmt"
	"reflect"
	"testing"
)

func meminit(memory []int) []uint16 {
	return sliceinit(memory, MemSize)
}

func stinit(stack []int) []uint16 {
	return sliceinit(stack, StackLimit)
}

func sliceinit(content []int, size int) []uint16 {
	sl := make([]uint16, size)
	for i, v := range content {
		sl[i] = uint16(v)
	}
	return sl
}

func eq(expected cpu, got cpu) (bool, error) {
	if got.sp != expected.sp {
		return false, fmt.Errorf(
			"wrong stack pointer value: %d, expected: %d",
			expected.sp, got.sp,
		)
	}

	if got.ip != expected.ip {
		return false, fmt.Errorf(
			"wrong instruction pointer value: %d, expected: %d",
			expected.ip, got.ip,
		)
	}

	if !reflect.DeepEqual(expected.stack, got.stack) {
		return false, fmt.Errorf(
			"stacks are not equal: %v, expeced: %v",
			expected.ip, got.ip,
		)
	}

	if !reflect.DeepEqual(expected.memory, got.memory) {
		return false, fmt.Errorf(
			"memsets are not equal: %v, expeced: %v",
			expected.ip, got.ip,
		)
	}

	return true, nil
}

func TestCpu_push(t *testing.T) {
	type args struct {
		n uint16
	}
	tests := []struct {
		name string
		c    cpu
		args args
		want cpu
	}{
		{
			name: "push should add value to a stack",
			args: args{1},
			c: cpu{
				sp:     0,
				memory: meminit([]int{}),
				stack:  stinit([]int{}),
			},
			want: cpu{
				sp:     1,
				memory: meminit([]int{}),
				stack:  stinit([]int{1}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.push(tt.args.n)

			equal, reason := eq(tt.want, tt.c)
			if !equal {
				t.Error(reason)
			}
		})
	}
}

func TestCpu_pop(t *testing.T) {
	tests := []struct {
		name  string
		c     cpu
		want  cpu
		want1 uint16
	}{
		{
			name: "pop should return top value",
			c: cpu{
				sp:    1,
				stack: stinit([]int{1}),
			},
			want: cpu{
				sp:    0,
				stack: stinit([]int{1}),
			},
			want1: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.pop(); got != tt.want1 {
				t.Errorf("cpu.pop() = %v, want %v", got, tt.want1)
			}

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_iadd(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop two elements and push their sum",
			c: cpu{
				sp:    2,
				stack: stinit([]int{1, 2}),
			},
			want: cpu{
				sp:    1,
				stack: stinit([]int{3, 0}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.iadd()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_isub(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop two elements from stack and push their difference",
			c: cpu{
				sp:    2,
				stack: stinit([]int{2, 1}),
			},
			want: cpu{
				sp:    1,
				stack: stinit([]int{1, 0}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.isub()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_iand(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise and",
			c: cpu{
				sp:    2,
				stack: stinit([]int{7, 5}),
			},
			want: cpu{
				sp:    1,
				stack: stinit([]int{5, 0}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.iand()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_ior(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise and",
			c: cpu{
				sp:    2,
				stack: stinit([]int{7, 5}),
			},
			want: cpu{
				sp:    1,
				stack: stinit([]int{7, 0}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ior()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_ixor(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise xor",
			c: cpu{
				sp:    2,
				stack: stinit([]int{7, 5}),
			},
			want: cpu{
				sp:    1,
				stack: stinit([]int{2, 0}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ixor()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_inot(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop elements from stack and push bitwise not",
			c: cpu{
				sp:    1,
				stack: stinit([]int{2}),
			},
			want: cpu{
				sp:    1,
				stack: stinit([]int{^2}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.inot()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_iload(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should load value from memory onto the stack",
			c: cpu{
				sp:     1,
				memory: meminit([]int{5}),
				stack:  stinit([]int{0}),
			},
			want: cpu{
				sp:     1,
				memory: meminit([]int{5}),
				stack:  stinit([]int{5}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.iload()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_istor(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should store value from stack into memory",
			c: cpu{
				sp:     2,
				stack:  stinit([]int{34, 1}),
				memory: meminit([]int{1, 2, 3}),
			},
			want: cpu{
				sp:     0,
				stack:  stinit([]int{34}),
				memory: meminit([]int{1, 34, 3}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.istor()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_ijmp(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop value from stack and goto there",
			c: cpu{
				sp:    1,
				ip:    0,
				stack: stinit([]int{42}),
			},
			want: cpu{
				sp:    0,
				ip:    42 + MemSize/2,
				stack: stinit([]int{42}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ijmp()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_ijz(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop value from stack and goto there",
			c: cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 0}),
			},
			want: cpu{
				sp:    0,
				ip:    42 + MemSize/2,
				stack: stinit([]int{42}),
			},
		},
		{
			name: "should pop value from stack and not goto there",
			c: cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 1}),
			},
			want: cpu{
				sp:    0,
				ip:    0,
				stack: stinit([]int{42}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ijz()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_ipush(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should push next word onto the stack",
			c: cpu{
				sp:     0,
				ip:     0,
				memory: meminit([]int{42}),
				stack:  stinit([]int{}),
			},
			want: cpu{
				sp:     1,
				ip:     1,
				memory: meminit([]int{42}),
				stack:  stinit([]int{42}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ipush()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_idup(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should duplicate stack top",
			c: cpu{
				sp:    1,
				stack: stinit([]int{42}),
			},
			want: cpu{
				sp:    2,
				stack: stinit([]int{42, 42}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.idup()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_iswap(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should swap two top stack values",
			c: cpu{
				sp:    2,
				stack: stinit([]int{24, 42}),
			},
			want: cpu{
				sp:    2,
				stack: stinit([]int{42, 24}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.iswap()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_irol3(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "(a, b, c) -> (b, c, a)",
			c: cpu{
				sp:    3,
				stack: stinit([]int{24, 42, 86}),
			},
			want: cpu{
				sp:    3,
				stack: stinit([]int{42, 86, 24}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.irol3()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_ijnz(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should pop value from stack and goto there",
			c: cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 1}),
			},
			want: cpu{
				sp:    0,
				ip:    42 + MemSize/2,
				stack: stinit([]int{42}),
			},
		},
		{
			name: "should pop value from stack and not goto there",
			c: cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 0}),
			},
			want: cpu{
				sp:    0,
				ip:    0,
				stack: stinit([]int{42}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ijnz()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_idrop(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should drop stack top",
			c: cpu{
				sp:    1,
				stack: stinit([]int{42}),
			},
			want: cpu{
				sp:    0,
				stack: stinit([]int{42}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.idrop()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_icomp(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want cpu
	}{
		{
			name: "should push top complement",
			c: cpu{
				sp:    1,
				stack: stinit([]int{42}),
			},
			want: cpu{
				sp:    1,
				stack: stinit([]int{-42}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.icomp()

			equal, problem := eq(tt.want, tt.c)
			if !equal {
				t.Error(problem)
			}
		})
	}
}

func TestCpu_MemDump(t *testing.T) {
	tests := []struct {
		name string
		c    cpu
		want []uint16
	}{
		{
			name: "should return copy of memory dump",
			c: cpu{
				memory: []uint16{1, 2, 3},
			},
			want: []uint16{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ttwant := make([]uint16, MemSize)
			copy(ttwant, tt.want)
			tt.want = ttwant

			if got := tt.c.MemDump(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cpu.MemDump() = %v, want %v", got, tt.want)
			}
		})
	}
}
