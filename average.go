package average

import "./mesh2d"

type Accumulator struct {
    value float64
    points uint64
    newValues chan float64
}

func (accum *Accumulator) listen() {
    for {
        accum.value += <-accum.newValues
        accum.points += 1
    }
}

func (accum *Accumulator) average() float64 {
    return accum.value / float64(accum.points)
}

func BuildAccumulator() *Accumulator {
    accum := new(Accumulator)
    accum.value = 0
    accum.newValues = make(chan float64)
    go accum.listen()
    return accum
}

type Consumer func(cmesh chan []float64, accum chan float64)

// pointsPerSide is uint32 so that accum.points will fit in a uint64
// numWorkers is uint16 to avoid spawning a ridiculous number of processes
func Average(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
    cmesh, done := mesh2d.Square(pointsPerSide)
    accum := BuildAccumulator()
    var i uint16 = 0
    for i = 0; i < numWorkers; i++ {
        go worker(cmesh, accum.newValues)
    }
    <-done
    return accum.average()
}
