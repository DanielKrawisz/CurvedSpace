package functions

import "../vector"

//These two are actually the same type, but the spherical color function
//is supposed to depend only on the direction of the vector given it.
//There is no real way to enforce that though. 
type SpacialFunction func([]float64) float64

type SphericalFunction func([]float64) float64

func ConstantFunction(z float64) (func([]float64) float64) {
  return func([]float64) float64 {
   return z
  }
}

// v - a position vector that gives the origin point of the checks.
// c - coordinate vectors describing the orientation of the checks.
//
//may return nil.
func Checks(v[] float64, c[][] float64) SpacialFunction {
  if v == nil || c == nil { return nil}

  if len(v) != len(c) { return nil }

  m := vector.Inverse(c)
  p := vector.MatrixMultiply(m, v)

  return func(position []float64) float64{
    var i int = 0
    for _, d := range vector.Minus(vector.MatrixMultiply(m, position), p) {
      i += int(d) % 2
    }

    if i % 2 == 0 {
      return 1
    } else {
      return -1
    } 
  }
}
