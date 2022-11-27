package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	NoOutput bool `short:"n" long:"no_output" description:"No output"`
}

func main() {
	if _, err := flags.ParseArgs(&opts, os.Args); err != nil {
		log.Fatal(err)
	}

	var n1 int
	fmt.Print("First matr n: ")
	if _, err := fmt.Scanln(&n1); err != nil {
		log.Fatal(err)
	}

	var m1 int
	fmt.Print("First matr m: ")
	if _, err := fmt.Scanln(&m1); err != nil {
		log.Fatal(err)
	}

	var n2 int
	fmt.Print("Second matr n: ")
	if _, err := fmt.Scanln(&n2); err != nil {
		log.Fatal(err)
	}

	var m2 int
	fmt.Print("Second matr m: ")
	if _, err := fmt.Scanln(&m2); err != nil {
		log.Fatal(err)
	}

	if m1 != n2 {
		log.Fatalf("Could not multiply matrixes with given dimensions")
	}

	log.Printf("Generating random matrixes %dx%d and %dx%d", n1, m1, n2, m2)

	matr1 := make([][]int, n1)
	for i := range matr1 {
		matr1[i] = make([]int, m1)
		for j := range matr1[i] {
			matr1[i][j] = rand.Intn(10)
		}
	}

	matr2 := make([][]int, n2)
	for i := range matr2 {
		matr2[i] = make([]int, m2)
		for j := range matr2[i] {
			matr2[i][j] = rand.Intn(10)
		}
	}

	maxTheads := runtime.NumCPU()
	for threads := 1; threads <= maxTheads; threads *= 2 {
		runtime.GOMAXPROCS(threads)

		func() {
			res := make([][]int, n1)
			for i := range res {
				res[i] = make([]int, m2)
			}

			now := time.Now()

			for i := range matr1 {
				for j := range matr1[i] {
					go func(i, j int) {
						temp := 0
						for k := range matr1[i] {
							temp += matr1[i][k] * matr2[k][j]
						}
						res[i][j] = temp
					}(i, j)
				}
			}

			fmt.Printf("threads: %d, time: %d ms\n", threads, time.Since(now).Milliseconds())
		}()

		fmt.Println()
	}

}
