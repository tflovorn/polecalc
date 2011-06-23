package polecalc

import "math"

// A function which does calculations based on data passed in on cmesh and returns results through accum
type Consumer func(point []float64) float64

// A type which can absorb grid points and return a result
type GridListener interface {
	initialize() GridListener
	grab(point []float64) GridListener
	result() interface{}
}

// --- Accumulator ---
// Collects values passed through grab() to find an average
type Accumulator struct {
	worker     Consumer // function to average
	value      float64  // sum of points seen so far
	compensate float64  // used in Kahan summation to correct floating-point error
	points     uint64   // number of points seen
}

func (accum Accumulator) initialize() GridListener {
	accum.value = 0.0
	accum.compensate = 0.0
	accum.points = 0
	return accum
}

// Handle new data.
// Use Kahan summation algorithm to reduce error: implementation cribbed from Wikipedia
func (accum Accumulator) grab(point []float64) GridListener {
	newValue := accum.worker(point)
	y := newValue - accum.compensate
	t := accum.value + y
	accum.compensate = (t - accum.value) - y
	accum.value = t
	accum.points += 1
	return accum
}

// Average of points passed in through grab()
func (accum Accumulator) result() interface{} {
	return accum.value / float64(accum.points)
}

// Create a new accumulator
func BuildAccumulator(worker Consumer) *Accumulator {
	accum := new(Accumulator)
	accum.worker = worker
	accum.initialize()
	return accum
}

// --- accumulator for minima ---
type MinimumData struct {
	worker  Consumer // function to minimize
	minimum float64  // minimum value seen so far
}

func (minData MinimumData) initialize() GridListener {
	minData.minimum = math.MaxFloat64
	return minData
}

func (minData MinimumData) grab(point []float64) GridListener {
	newValue := minData.worker(point)
	if newValue < minData.minimum {
		minData.minimum = newValue
	}
	return minData
}

func (minData MinimumData) result() interface{} {
	return minData.minimum
}

func BuildMinimumData(worker Consumer) *MinimumData {
	minData := new(MinimumData)
	minData.worker = worker
	minData.initialize()
	return minData
}

// --- accumulator for maximua ---
// it'd be nice to combine this with MaximumData but maybe would lose some
// speed - most common (?) use case is minimizing epsilon after changing D1
type MaximumData struct {
	worker  Consumer
	maximum float64
}

func (maxData MaximumData) initialize() GridListener {
	maxData.maximum = -math.MaxFloat64
	return maxData
}

func (maxData MaximumData) grab(point []float64) GridListener {
	newValue := maxData.worker(point)
	if newValue > maxData.maximum {
		maxData.maximum = newValue
	}
	return maxData
}

func (maxData MaximumData) result() interface{} {
	return maxData.maximum
}

func BuildMaximumData(worker Consumer) *MaximumData {
	maxData := new(MaximumData)
	maxData.worker = worker
	maxData.initialize()
	return maxData
}
/*
// --- accumulator for (discrete approximation) delta functions ---
type DeltaBinner struct {
	DeltaTerms func(q []float64) ([]float64, []float64)
	BinsStart, BinsStop float64
	NumBins uint
	Bins []float64
}

func (binner DeltaBinner) initialize() GridListener {

}

func (binner DeltaBinner) grab(point []float64) GridListener {

}

func (binner DeltaBinner) result(
*/
// -- utility functions --
// assumes numWorkers > 0
func DoGridListen(pointsPerSide uint32, numWorkers uint16, listener GridListener) interface{} {
	cmesh := Square(pointsPerSide)
	done := make(chan bool)
	listener = listener.initialize()
	var i uint16 = 0
	for i = 0; i < numWorkers; i++ {
		go func() {
			for point, ok := <-cmesh; ok; point, ok = <-cmesh {
				listener = listener.grab(point)
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
	accum := BuildAccumulator(worker)
	return DoGridListen(pointsPerSide, numWorkers, *accum).(float64)
}

func Minimum(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
	minData := BuildMinimumData(worker)
	return DoGridListen(pointsPerSide, numWorkers, *minData).(float64)
}

func Maximum(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
	maxData := BuildMaximumData(worker)
	return DoGridListen(pointsPerSide, numWorkers, *maxData).(float64)
}
