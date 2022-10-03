package emu

import (
	"fmt"
	"reflect"
	"testing"
)

func meminit(memory []int) []uint32 {
	return sliceinit(memory, MEM_SIZE)
}

func stinit(stack []int) []uint32 {
	return sliceinit(stack, STACK_LIMIT)
}

func sliceinit(content []int, size int) []uint32 {
	sl := make([]uint32, size)
	for i, v := range content {
		sl[i] = uint32(v)
	}
	return sl
}

func eq(expected Cpu, got Cpu) (bool, error) {
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
		n uint32
	}
	tests := []struct {
		name string
		c    Cpu
		args args
		want Cpu
	}{
		{
			name: "push should add value to a stack",
			args: args{1},
			c: Cpu{
				sp:     0,
				memory: meminit([]int{}),
				stack:  stinit([]int{}),
			},
			want: Cpu{
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
		c     Cpu
		want  Cpu
		want1 uint32
	}{
		{
			name: "pop should return top value",
			c: Cpu{
				sp:    1,
				stack: stinit([]int{1}),
			},
			want: Cpu{
				sp:    0,
				stack: stinit([]int{1}), // stack memory is not cleared
			},
			want1: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.pop(); got != tt.want1 {
				t.Errorf("Cpu.pop() = %v, want %v", got, tt.want1)
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop two elements and push their sum",
			c: Cpu{
				sp:    2,
				stack: stinit([]int{1, 2}),
			},
			want: Cpu{
				sp:    1,
				stack: stinit([]int{3, 2}), // stack memory is not cleared
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop two elements from stack and push their difference",
			c: Cpu{
				sp:    2,
				stack: stinit([]int{2, 1}),
			},
			want: Cpu{
				sp:    1,
				stack: stinit([]int{1, 1}), // stack memory is not cleared
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise and",
			c: Cpu{
				sp:    2,
				stack: stinit([]int{7, 5}),
			},
			want: Cpu{
				sp:    1,
				stack: stinit([]int{5, 5}), // stack memory is not cleared
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise and",
			c: Cpu{
				sp:    2,
				stack: stinit([]int{7, 5}),
			},
			want: Cpu{
				sp:    1,
				stack: stinit([]int{7, 5}), // stack memory is not cleared
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop two elements from stack and push bitwise xor",
			c: Cpu{
				sp:    2,
				stack: stinit([]int{7, 5}),
			},
			want: Cpu{
				sp:    1,
				stack: stinit([]int{2, 5}), // stack memory is not cleared
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop elements from stack and push bitwise not",
			c: Cpu{
				sp:    1,
				stack: stinit([]int{2}),
			},
			want: Cpu{
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should load value from memory onto the stack",
			c: Cpu{
				sp:     1,
				memory: meminit([]int{5}),
				stack:  stinit([]int{0}),
			},
			want: Cpu{
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should store value from stack into memory",
			c: Cpu{
				sp:     2,
				stack:  stinit([]int{34, 1}),
				memory: meminit([]int{1, 2, 3}),
			},
			want: Cpu{
				sp:     0,
				stack:  stinit([]int{34, 1}),
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop value from stack and goto there",
			c: Cpu{
				sp:    1,
				ip:    0,
				stack: stinit([]int{42}),
			},
			want: Cpu{
				sp:    0,
				ip:    42 + MEM_SIZE/2,
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop value from stack and goto there",
			c: Cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 0}),
			},
			want: Cpu{
				sp:    0,
				ip:    42 + MEM_SIZE/2,
				stack: stinit([]int{42, 0}),
			},
		},
		{
			name: "should pop value from stack and not goto there",
			c: Cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 1}),
			},
			want: Cpu{
				sp:    0,
				ip:    0,
				stack: stinit([]int{42, 1}),
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should push next word onto the stack",
			c: Cpu{
				sp:     0,
				ip:     0,
				memory: meminit([]int{42}),
				stack:  stinit([]int{}),
			},
			want: Cpu{
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should duplicate stack top",
			c: Cpu{
				sp:    1,
				stack: stinit([]int{42}),
			},
			want: Cpu{
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should swap two top stack values",
			c: Cpu{
				sp:    2,
				stack: stinit([]int{24, 42}),
			},
			want: Cpu{
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
		c    Cpu
		want Cpu
	}{
		{
			name: "(a, b, c) -> (b, c, a)",
			c: Cpu{
				sp:    3,
				stack: stinit([]int{24, 42, 86}),
			},
			want: Cpu{
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should pop value from stack and goto there",
			c: Cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 1}),
			},
			want: Cpu{
				sp:    0,
				ip:    42 + MEM_SIZE/2,
				stack: stinit([]int{42, 1}),
			},
		},
		{
			name: "should pop value from stack and not goto there",
			c: Cpu{
				sp:    2,
				ip:    0,
				stack: stinit([]int{42, 0}),
			},
			want: Cpu{
				sp:    0,
				ip:    0,
				stack: stinit([]int{42, 0}),
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should drop stack top",
			c: Cpu{
				sp:    1,
				stack: stinit([]int{42}),
			},
			want: Cpu{
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
		c    Cpu
		want Cpu
	}{
		{
			name: "should push top complement",
			c: Cpu{
				sp:    1,
				stack: stinit([]int{42}),
			},
			want: Cpu{
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
		c    Cpu
		want []uint32
	}{
		{
			name: "should return copy of memory dump",
			c: Cpu{
				memory: []uint32{1, 2, 3},
			},
			want: []uint32{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ttwant := make([]uint32, MEM_SIZE)
			copy(ttwant, tt.want)
			tt.want = ttwant

			if got := tt.c.MemDump(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cpu.MemDump() = %v, want %v", got, tt.want)
			}
		})
	}
}
