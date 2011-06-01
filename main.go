package main

import "fmt"
import "math"
import "./average"

func worker(cmesh chan []float64, accum chan float64) {
    for {
        k := <-cmesh
        kx, ky := k[0], k[1]
        accum <- math.Sin(kx) * math.Sin(ky)
    }
}

func main() {
    fmt.Println(average.Average(128, worker, 8))
}
