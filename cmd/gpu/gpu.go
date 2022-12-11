package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"sync"

	"github.com/aveplen/sm/internal"
)

const matrixSize = 10

const sourceCode = `
/* 00 */   start:             //
/* 01 */     push             // load arr len
/* 02 */     0                //
/* 03 */     load             //
                              //
/* 04 */     dup              // save first length into arr_len1
/* 05 */     push             //
/* 06 */     0                // arr len
/* 07 */     stor             //
                              //
/* 08 */     dup              // goto final_routine if arr len == 0
/* 09 */     push             //
/* 10 */     &final_routine   //
/* 11 */     swap             //
/* 12 */     jz               //
                              //
/* 13 */     stc              // counter = arr len
                              //
/* 14 */   mult_routine:      //
/* 15 */     cts              // get counter
/* 16 */     load             // load element of the first array (counter)
                              //
/* 17 */     cts              // get counter
/* 18 */     push             //
/* 19 */     0                // push arr len addr
/* 20 */     load             // load arr len
/* 21 */     push             //
/* 22 */     1                //
/* 23 */     add              //
/* 24 */     add              //
/* 25 */     load             // load element of the second array (counter + len1 + 1)
                              // additional 1 is for arr2 len, which is not used in program
                              //
/* 26 */     mul              // arr1[i] * arr2[i]
                              //
/* 27 */     cdec             // counter --
                              //
/* 28 */     cts              // jump to sum_routine if counter == 0
/* 29 */     push             //
/* 30 */     &sum_routine     //
/* 31 */     swap             //
/* 32 */     jz               //
                              //
/* 33 */     push             // if counter != 0 continue multiplying on a stack
/* 34 */     &mult_routine    //
/* 35 */     jmp              //
                              //
/* 36 */   sum_routine:       //
/* 37 */     push             // counter = len(arr)-1
/* 38 */     0                // arr len
/* 39 */     load             //
/* 40 */     stc              //
/* 41 */     cdec             //
                              //
/* 42 */   while:             //
/* 43 */     add              // arr1[i-1]*arr2[i-1] + arr1[i]*arr2[i]
/* 44 */     cdec             // counter --
                              //
/* 45 */     cts              // if counter == 0 goto final_routine
/* 46 */     push             //
/* 47 */     &final_routine   //
/* 48 */     swap             //
/* 49 */     jz               //
                              //
/* 50 */     push             // goto while
/* 51 */     &while           //
/* 52 */     jmp              //
                              //
/* 53 */   final_routine:     //
/* 54 */     term             //
`

var program = func() []uint16 {
	buffer := make([]byte, len(sourceCode))
	copy(buffer, []byte(sourceCode))
	reader := bytes.NewReader(buffer)
	prog, _ := internal.Compile(*bufio.NewReader(reader), false)
	return prog
}()

func main() {
	matr1 := generateMatrix(matrixSize)
	outputMatrix(matr1)

	matr2 := generateMatrix(matrixSize)
	outputMatrix(matr2)

	res := make([][]int, matrixSize)
	for i := range res {
		res[i] = make([]int, matrixSize)
	}

	wg := sync.WaitGroup{}
	for i := 0; i < matrixSize; i++ {
		for j := 0; j < matrixSize; j++ {

			data := make([]uint16, (matrixSize+1)*2)

			data[0] = matrixSize
			for k := 0; k < matrixSize; k++ {
				data[k+1] = uint16(matr1[i][k])
			}

			data[matrixSize] = matrixSize
			for k := 0; k < matrixSize; k++ {
				data[matrixSize+k+1] = uint16(matr2[k][j])
			}

			wg.Add(1)
			go func(program, data []uint16, i, j int) {
				wg.Done()

				cpu := internal.WithMemProg(program, data)
				cpu.Run()
				res[i][j] = int(cpu.StackDump()[0])
			}(program, data, i, j)
		}
	}

	wg.Wait()

	outputMatrix(res)
}

func generateMatrix(size int) [][]int {
	matr := make([][]int, matrixSize)
	for i := range matr {
		matr[i] = make([]int, matrixSize)
		for j := range matr[i] {
			matr[i][j] = rand.Intn(10)
		}
	}
	return matr
}

func outputMatrix(matr [][]int) {
	for i := range matr {
		for j := range matr[i] {
			fmt.Printf("%d ", matr[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}
