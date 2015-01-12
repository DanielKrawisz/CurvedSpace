package surface

import "../vector"

//v should be a unit vector.
//May return nil
func NewTorus(p []float64, v []float64, R, r float64) Surface {
  if p == nil || v == nil {
    return nil
  }

  dim := len(p)
  if dim != len(v) {
    return nil
  }

  rr := r*r
  RR := R*R
  rR := 2*rr + 4*RR

  return NewQuarticSurface(p, [][]float64{},
    [][]float64{[]float64{1, 0, 0}, []float64{0, 1, 0}, []float64{0, 0, 1}}, [][]float64{},
    [][]float64{[]float64{rR, 0, 0}, []float64{0, rR, 0}, []float64{0, 0, rR}}, 
    [][]float64{[]float64{2*RR, 0, 0}, []float64{0, 2*RR, 0}, []float64{0, 0, 2*RR},
      vector.Times(4*RR, v)}, 
    []float64{0, 0, 0}, -rr*rr - RR*RR + 2 * rr * RR)
}
