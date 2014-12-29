package distributions

type SampleStatistics struct {
  n int
  k, ex, ex2 float64
}
 
func (v *SampleStatistics) AddVariable(x float64) {
  if (v.n == 0) {
    v.k = x
  }
  v.n ++
  v.ex = v.ex + (x - v.k)
  v.ex2 = v.ex2 + (x - v.k) * (x - v.k)
}
 
func (v *SampleStatistics) RemoveVariable(x float64) {
  v.n --
  v.ex = v.ex - (x - v.k)
  v.ex = v.ex2 - (x - v.k) * (x - v.k)
}
 
func (v *SampleStatistics) Mean() float64 {
  return v.k + v.ex / float64(v.n)
}
 
func (v *SampleStatistics) Variance() float64 {
  return (v.ex2 - (v.ex * v.ex) / float64(v.n)) / float64(v.n-1)
}

func (v *SampleStatistics) MeanVariance() float64 {
  return v.Variance() / float64(v.n)
}

func NewSampleStatistics() *SampleStatistics {
  return &SampleStatistics{0, 0, 0, 0}
}
