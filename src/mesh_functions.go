package polecalc

import "math"

// A function which does calculations based on data passed in on cmesh and returns results through accum
type Consumer func(point []float64) float64

// A type which can absorb grid points and return a result
type GridListener interface {
	initialize() GridListener
	grab(newValue float64) GridListener
	result() float64
}

// --- Accumulator ---
// Collects values passed through grab() to find an average
type Accumulator struct {
	value      float64 // sum of points seen so far
	compensate float64 // used in Kahan summation to correct floating-point error
	points     uint64  // number of points seen
}

func (accum Accumulator) initialize() GridListener {
	accum.value = 0.0
	accum.compensate = 0.0
	accum.points = 0
	return accum
}

// Handle new data.
// Use Kahan summation algorithm to reduce error: implementation cribbed from Wikipedia
func (accum Accumulator) grab(newValue float64) GridListener {
	y := newValue - accum.compensate
	t := accum.value + y
	accum.compensate = (t - accum.value) - y
	accum.value = t
	accum.points += 1
	return accum
}

// Average of points passed in through grab()
func (accum Accumulator) result() float64 {
	return accum.value / float64(accum.points)
}

// Create a new accumulator
func BuildAccumulator() *Accumulator {
	accum := new(Accumulator)
	accum.initialize()
	return accum
}

// --- MinimumData ---
type MinimumData struct {
	minimum float64
}

func (minData MinimumData) initialize() GridListener {
	minData.minimum = math.MaxFloat64
	return minData
}

func (minData MinimumData) grab(newValue float64) GridListener {
	if newValue < minData.minimum {
		minData.minimum = newValue
	}
	return minData
}

func (minData MinimumData) result() float64 {
	return minData.minimum
}

func BuildMinimumData() *MinimumData {
	minData := new(MinimumData)
	minData.initialize()
	return minData
}

// -- utility functions --
// assumes numWorkers > 0
func DoGridListen(pointsPerSide uint32, worker Consumer, numWorkers uint16, listener GridListener) float64 {
	cmesh := Square(pointsPerSide)
	done := make(chan bool)
	listener = listener.initialize()
	var i uint16 = 0
	for i = 0; i < numWorkers; i++ {
		go func() {
			for point, ok := <-cmesh; ok; point, ok = <-cmesh {
				listener = listener.grab(worker(point))
			}
			done <- true
		}()
	}
	for doneCount := 0; doneCount < int(numWorkers); doneCount++ {
		<-done
	}
	return listener.result()
}

// Find the average over a square grid of the function given by worker.
// Spawn number of goroutines given by numWorkers.
// pointsPerSide is uint32 so that accum.points will fit in a uint64.
// numWorkers is uint16 to avoid spawning a ridiculous number of processes.
// Consumer is defined in utility.go
func Average(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
	accum := BuildAccumulator()
	return DoGridListen(pointsPerSide, worker, numWorkers, *accum)
}

func Minimum(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
	minData := BuildMinimumData()
	return DoGridListen(pointsPerSide, worker, numWorkers, *minData)
}
