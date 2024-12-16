package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

func dft(x []float64) []complex128 {
	fx := make([]complex128, len(x))
	var temp1, temp2 float64
	for i := range x {
		fx[i] = complex(0, 0)
		temp1 = 2. * math.Pi * float64(i) / float64(len(x))
		for j := range x {
			temp2 = temp1 * float64(j)
			fx[i] += complex(x[j] * math.Cos(temp2), -x[j] * math.Sin(temp2))
		}
	}
	return fx
} 

func dft_gr(x []float64, rnum uint16) []complex128 {
	fx := make([]complex128, len(x))
	var seglen int = len(x) / (int(rnum) - 1)
	if (seglen <= 0) { 
		fmt.Println("too many goroutines")
		os.Exit(1)
	}
	/* Defining wait group (to know when all routines end) */ 
	var wg sync.WaitGroup
	wg.Add(int(rnum)) 
	/* Every routine will calculate segment of fx array
	 * I do not use mutexes here because it will affect performance 
	 * and I'm sure that different coroutines will not be able to write
	 * to the same fx element */
	routine := func(start int, end int) {
		defer wg.Done() // Changes wait group counter
		var temp1, temp2 float64
		for ; start < end; start++ {
			fx[start] = complex(0, 0)
			temp1 = 2. * math.Pi * float64(start) / float64(len(x))
			for j := range x {
				temp2 = temp1 * float64(j)
				fx[start] += complex(x[j] * math.Cos(temp2), -x[j] * math.Sin(temp2)) 
			}
		}
	}
	for start := 0; start < seglen * (int(rnum) - 1); start += seglen {
		go routine(start, start + seglen)
	}
	go routine(seglen * (int(rnum) - 1), len(x)) 
	/* Yes, main routine does nothing, just waits. I don't think it will 
	 * have big impact on program performance */
	wg.Wait() // Wait for all goroutines to end

	return fx
}

func dft_reverse(fx []complex128) []float64 {
	xnew := make([]float64, len(fx))
	var temp1, temp2 float64
	for i := range xnew {
		xnew[i] = 0
		temp1 = 2. * math.Pi * float64(i) / float64(len(xnew))
		for j := range xnew {
			temp2 = temp1 * float64(j)
			xnew[i] += real(fx[j] * complex(math.Cos(temp2), math.Sin(temp2)))
		}
		xnew[i] /= float64(len(xnew))
	}
	return xnew
}

func dft_reverse_gr(fx []complex128, rnum uint16) []float64 {
	xnew := make([]float64, len(fx))
	var seglen int = len(fx) / (int(rnum) - 1)
	if (seglen <= 0) { 
		fmt.Println("too many goroutines")
		os.Exit(1)
	}
	/* Defining wait group (to know when all routines end) */ 
	var wg sync.WaitGroup
	wg.Add(int(rnum)) 
	routine := func(start int, end int) {
		defer wg.Done() // Changes wait group counter
		var temp1, temp2 float64
		for ; start < end; start++ {
			xnew[start] = complex(0, 0)
			temp1 = 2. * math.Pi * float64(start) / float64(len(fx))
			for j := range fx {
				temp2 = temp1 * float64(j)
				xnew[start] += real(fx[j] * complex(math.Cos(temp2), math.Sin(temp2)))
			}
			xnew[start] /= float64(len(xnew))
		}
	}
	for start := 0; start < seglen * (int(rnum) - 1); start += seglen {
		go routine(start, start + seglen)
	}
	go routine(seglen * (int(rnum) - 1), len(fx)) 
	wg.Wait() // Wait for all goroutines to end

	return xnew
}

func compare(x1 []float64, x2 []float64) float64 {
	if len(x1) != len(x2) {
		fmt.Println("compare: len(x1) != len(x2)")
		os.Exit(1)
	}
	var avgdiff float64 = 0
	for i := range x1 {
		avgdiff += math.Abs(x1[i] - x2[i])
	}
	avgdiff /= float64(len(x1))
	
	return avgdiff
}

func main() {
	const max_goroutines = 15
	/* Declaring and reading array (x) 
	 * We suppose that file/CLI input will contain only correct values */
	var datalen int32
	fmt.Fscan(os.Stdin, &datalen)
	var x []float64 = make([]float64, datalen)
	for i := range x {
		fmt.Fscan(os.Stdin, &x[i])
	}

	start := time.Now()
	fx := dft(x)
	xnew := dft_reverse(fx)
	duration := time.Since(start)
	diff := compare(x, xnew)
	fmt.Println(" Single goroutine:\n* Time", duration.Seconds(), "s.\n* Average accuracy:", diff)

	for n := 2; n <= max_goroutines; n++ {
		start := time.Now()
		fx := dft_gr(x, uint16(n))
		xnew := dft_reverse_gr(fx,uint16(n))
		duration := time.Since(start)
		diff := compare(x, xnew)
		fmt.Println("\n", n, "goroutines:\n* Time", duration.Seconds(), "s.\n* Average accuracy:", diff)
	}
}