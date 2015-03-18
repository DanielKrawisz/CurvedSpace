package polynomialsurfaces

import "github.com/DanielKrawisz/CurvedSpace/surface"

//The torus corresponding to the equation 
//
//  (x.x + R^2 - r^2)^2 == 4 R^2 x.(I - v v).x
//
// v should be a unit vector and r < R. 
// p is a translation vector that takes x -> x - p
//
// This torus is three-dimensional!!! 
//
//May return nil
func NewTorus(p []float64, v []float64, R, r float64) surface.Surface {
  if p == nil || v == nil {
    return nil
  }

  if 3 != len(v) || 3 != len(p) {
    return nil
  }

  rr  := r*r
  RR  := R*R
  rR  := 2 * (rr + RR)
  RR4 := -RR * 4
  var t float64 = -1./3.

  return translateQuartic(&quarticSurface{3,
    [][][][]float64{
      [][][]float64{
        [][]float64{[]float64{-1}}},
      [][][]float64{
        [][]float64{[]float64{0}},
        [][]float64{[]float64{t}, []float64{0, -1}}},
      [][][]float64{
        [][]float64{[]float64{0}},
        [][]float64{[]float64{0}, []float64{0, 0}},
        [][]float64{[]float64{t}, []float64{0, t}, []float64{0, 0, -1}}}},
    [][][]float64{
      [][]float64{[]float64{0}},
      [][]float64{[]float64{0}, []float64{0, 0}},
      [][]float64{[]float64{0}, []float64{0, 0}, []float64{0, 0, 0}}}, 
    [][]float64{
      []float64{rR + RR4 * v[0] * v[0]},
      []float64{RR4 * v[1] * v[0], rR + RR4 * v[1] * v[1]},
      []float64{RR4 * v[2] * v[0], RR4 * v[2] * v[1], rR + RR4 * v[2] * v[2]}},
    []float64{0, 0, 0}, -rr*rr - RR*RR + 2 * rr * RR}, p)
}
