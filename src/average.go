package polecalc

// Collects values passed on the newValues channel to find an average
type Accumulator struct {
	value      float64      // sum of points seen so far
	compensate float64      // used in Kahan summation to correct floating-point error
	points     uint64       // number of points seen
	newValues  chan float64 // channel for new values to be summed
}

// Handle data when it shows up on newValues.
// Use Kahan summation algorithm to reduce error: implementation cribbed from Wikipedia
func (accum *Accumulator) listen() {
	for {
		y := <-accum.newValues - accum.compensate
		t := accum.value + y
		accum.compensate = (t - accum.value) - y
		accum.value = t
		accum.points += 1
	}
}

// Average of points passed in through newValues
func (accum *Accumulator) average() float64 {
	return accum.value / float64(accum.points)
}

// Create a new accumulator listening on its channel
func BuildAccumulator() *Accumulator {
	accum := new(Accumulator)
	accum.newValues = make(chan float64)
	go accum.listen()
	return accum
}

// A function which does calculations based on data passed in on cmesh and returns results through accum
type Consumer func(point []float64) float64

// Find the average over a square grid of the function given by worker.
// Spawn number of goroutines given by numWorkers.
// pointsPerSide is uint32 so that accum.points will fit in a uint64.
// numWorkers is uint16 to avoid spawning a ridiculous number of processes.
func Average(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
	cmesh, done := Square(pointsPerSide)
	accum := BuildAccumulator()
	var i uint16 = 0
	for i = 0; i < numWorkers; i++ {
		go func() {
			for {
				point := <-cmesh
				accum.newValues <- worker(point)
			}
		}()
	}
	<-done
	return accum.average()
}
