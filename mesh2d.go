package mesh2d

import "math"

func Square(pointsPerSide uint32) chan []float64 {
    cmesh := make(chan []float64)
    go helpSquare(cmesh, pointsPerSide)
    return cmesh
}

func helpSquare(cmesh chan []float64, pointsPerSide uint32) {
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
}
