package internal

type Mat struct {
	Slice  []float64
	width  int
	height int
}

// Set - set value
func (m *Mat) Set(x, y int, value float64) {
	m.Slice[m.width*y+x] = value
}

// At - get value at (x, y)
func (m *Mat) At(x, y int) float64 {
	return m.Slice[m.width*y+x]
}

// MapToRange - map all values to TargetMin/TargetMax (Normalize)
func (m *Mat) MapToRange(TargetMin, TargetMax float64) {
	// find min/max
	min, max := m.Slice[0], m.Slice[0]
	for _, v := range m.Slice {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	// Normalize
	for i := range m.Slice {
		m.Slice[i] = TargetMin + ((TargetMax-TargetMin)/(max-min))*(m.Slice[i]-min)
	}

}

func NewMatrix(w, h int) Mat {
	var mat = Mat{
		width:  w,
		height: h,
		Slice:  make([]float64, w*h),
	}
	return mat
}
