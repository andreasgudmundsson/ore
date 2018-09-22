package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	i := float64(0)
	for {
		fmt.Printf("%0.4f\n", 30*math.Sin(i/10.0*math.Pi))
		i += 1
		time.Sleep(200 * time.Millisecond)
	}
}
