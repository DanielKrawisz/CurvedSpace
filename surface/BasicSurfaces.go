package surface

import "./polynomials"

//Use polynomial surfaces and solid constructive geometry to
//make some primitive shapes.

//TODO plane, cylinder, cone, torus, elipse, paraboloid, hyperboloid, box, elliptic curve.

type Sphere interface {
  Surface
  X() []float64
  R2() float64
}

type sphere struct {
  dim int
  p []float64
  p2 float64
  r2 float64 //Must not be negative; may be zero or infinity.
}

func (s *sphere) R2() float64 {
  return s.r2
}

func (s *sphere) X() []float64 {
  return s.p
}

func (s *sphere) Dimension() int {
  return s.dim
}

func (s *sphere) F(x []float64) float64 {
  var d float64 = 0
  for i := 0; i < s.dim; i++ {
    z := x[i] - s.p[i]
    d += z*z
  }
  return s.r2 - d
}

func (s *sphere) Intersection(x, v []float64) []float64 {
  var px, xx, vv, vpx float64
  for i := 0; i < s.dim; i++ {
    vpx += v[i] * (x[i] - s.p[i])
    vv += v[i] * v[i]
    xx += x[i] * x[i]
    px += s.p[i] * x[i]
  }

  //v needs to have a length greater than zero, so don't need to take account of divide by zero.
  return polynomials.QuadraticFormula((s.p2 - s.r2 - 2 * px + xx) / vv, 2 * vpx / vv)
}

//(x - p).(x - p) < r2
func (s *sphere) Interior(x []float64) bool {
  return s.F(x) > 0
}

// -2 (x - p) 
func (s *sphere) Gradient(x []float64) []float64 {
  z := make([]float64, s.dim)
  for i := 0; i < s.dim; i++ {
    z[i] = 2 * (x[0] - s.p[0])
  }
  return z
}

//May return nil
func NewSphere(p []float64, r float64) Sphere {
  if p == nil {return nil}
  var p2 float64 = 0.0
  for i := 0; i < len(p); i++ {
    p2 += p[i] * p[i]
  }
  return &sphere{len(p), p, p2, r*r}
}

//May return nil
func NewPlaneByPointAndNormal(point, norm []float64) Surface {
  if point == nil || norm == nil {return nil}
  if len(point) != len(norm) {return nil}
  var b2, a float64 = 0, 0
  for i := 0; i < len(norm); i++ {
    b2 += norm[i] * norm[i]
    a += norm[i] * point[i]
  }
  //Ensure that norm is actually a normalizable vector. 
  if b2 == 0.0 {return nil}
  a /= b2
  b := make([]float64, len(norm))
  for i := 0; i < len(point); i++ {
    b[i] = norm[i] / b2
  }
  return &LinearCurve{len(norm), b, a}
}
