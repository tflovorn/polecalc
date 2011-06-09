package mesh

import "math"

func Square(pointsPerSide uint32) (chan []float64, chan bool) {
    cmesh := make(chan []float64)
    done := make(chan bool)
    go helpSquare(cmesh, done, pointsPerSide)
    return cmesh, done
}

func helpSquare(cmesh chan []float64, done chan bool, pointsPerSide uint32) {
    var x, y, step, length float64
    length = 2 * math.Pi
    step = length / float64(pointsPerSide)
    x, y = -math.Pi, -math.Pi
    for y < math.Pi {
        for x < math.Pi {
            cmesh <- []float64{x, y}
            x += step
        }
        y += step
        x = -math.Pi
    }
    for {
        done <- true;
    }
}
