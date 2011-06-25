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
	accum.value, accum.compensate = KahanSum(newValue, accum.value, accum.compensate)
	accum.points += 1
	return accum
}

func KahanSum(extraValue, oldValue, compensate float64) (float64, float64) {
	y := extraValue - compensate
	newValue := oldValue + y
	newCompensate := (newValue - oldValue) - y
	return newValue, newCompensate
}

// Average of points passed in through grab()
func (accum Accumulator) result() interface{} {
	return accum.value / float64(accum.points)
}

// Create a new accumulator
func NewAccumulator(worker Consumer) *Accumulator {
	accum := new(Accumulator)
	accum.worker = worker
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

func NewMinimumData(worker Consumer) *MinimumData {
	minData := new(MinimumData)
	minData.worker = worker
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

func NewMaximumData(worker Consumer) *MaximumData {
	maxData := new(MaximumData)
	maxData.worker = worker
	return maxData
}

// --- accumulator for (discrete approximation) delta functions ---

// returns pair of slices of bin variable values and their associciated 
// coefficients which are affected at the given point
// (for Gc0, bin variable is omega)
type DeltaTermsFunc func(point []float64) ([]float64, []float64)

type DeltaBinner struct {
	deltaTerms        DeltaTermsFunc
	binStart, binStop float64
	numBins           uint
	bins              []float64 // value of the function at various omega values
	compensates       []float64 // compensation values for Kahan summation
	numPoints         uint64
}

func (binner DeltaBinner) initialize() GridListener {
	binner.numPoints = 0
	for i, _ := range binner.bins {
		binner.bins[i] = 0.0
		binner.compensates[i] = 0.0
	}
	return binner
}

func (binner DeltaBinner) grab(point []float64) GridListener {
	omegas, coeffs := binner.deltaTerms(point)
	for i, omega := range omegas {
		n := binner.BinVarToIndex(omega)
		binner.bins[n], binner.compensates[n] = KahanSum(coeffs[i], binner.bins[n], binner.compensates[n])
	}
	binner.numPoints += 1
	return binner
}

func (binner DeltaBinner) result() interface{} {
	result := make([]float64, binner.numBins)
	for i, val := range binner.bins {
		result[i] = val / float64(binner.numPoints)
	}
	return result
}

func (binner DeltaBinner) Step() float64 {
	return math.Fabs(binner.binStop-binner.binStart) / float64(binner.numBins)
}

func (binner DeltaBinner) BinVarToIndex(binVar float64) int {
	return int(math.Floor((binVar - binner.binStart) / binner.Step()))
}

func (binner DeltaBinner) IndexToBinVar(index int) float64 {
	return binner.binStart + binner.Step()*float64(index)
}

func (binner DeltaBinner) BinVarValues() []float64 {
	values := make([]float64, binner.numBins)
	for i, _ := range values {
		values[i] = binner.IndexToBinVar(i)
	}
	return values
}

func NewDeltaBinner(deltaTerms DeltaTermsFunc, binStart, binStop float64, numBins uint) *DeltaBinner {
	if binStart > binStop {
		binStart, binStop = binStop, binStart
	}
	bins, compensates := make([]float64, numBins), make([]float64, numBins)
	binner := &DeltaBinner{deltaTerms, binStart, binStop, numBins, bins, compensates, 0}
	return binner
}

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
	accum := NewAccumulator(worker)
	return DoGridListen(pointsPerSide, numWorkers, *accum).(float64)
}

func Minimum(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
	minData := NewMinimumData(worker)
	return DoGridListen(pointsPerSide, numWorkers, *minData).(float64)
}

func Maximum(pointsPerSide uint32, worker Consumer, numWorkers uint16) float64 {
	maxData := NewMaximumData(worker)
	return DoGridListen(pointsPerSide, numWorkers, *maxData).(float64)
}

// Instead of taking a worker directly, this functions takes a *DeltaBinner
// (to avoid passing in all the params for DeltaBinner)
func DeltaBin(pointsPerSide uint32, deltaTerms *DeltaBinner, numWorkers uint16) []float64 {
	return DoGridListen(pointsPerSide, numWorkers, deltaTerms).([]float64)
}
