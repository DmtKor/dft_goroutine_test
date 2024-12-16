package main

import (
	"fmt"
	"os"
	"strconv"
	"math"
)

func usage() {
	fmt.Println("Usage:", os.Args[0], "<data length>")
}

func generate(n int) {
	var val float64
	for ; n > 0; n-- {
		val = math.Sin(float64(n % 10) / 10)
		fmt.Println(val)
	}
}

func main() {
	if (len(os.Args[1:]) != 1) {
		usage()
		return
	}
	if datalen, err := strconv.Atoi(os.Args[1]); err == nil {
		if (datalen <= 0) {
			fmt.Println("data length should be greater than zero")
			return
		}
		fmt.Println(datalen)
		generate(datalen)
	} else {
		usage()
		return
	}
}