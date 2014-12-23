package surface

import "math"

//A surface is here defined as an equation f(x_i) == 0, although
//this is slightly different from what a surface normally is
//mathematically. Here a surface is always an (n-1)-dimensional
//object embedded in n-dimensional space. A two-surface in three
//-space is the ordinary case, but we'll have occaision to work
//with more than that. 
type Surface interface {
  //The number of dimensions of the space in which a surface is embedded.
  //This is the length of the vectors expected by Intersection, 
  //Interior, and Gradiant, and the length of the vectors returned by
  //Intersection and Gradiant.
  Dimension() int
  //The value of the function at a given point. 
  F(x []float64) float64
  //Returns intersection parameters which interprets x and v as a
  //parametric line. Returns infinity if the line does not intersect
  //the surface. Returns an empty list if there is no intersection.
  //May return several intersection parameters. 
  Intersection(x, v []float64) []float64
  //The gradianThe normal to a surface at a given point on the surface.
  Gradient(x []float64) []float64
  //Turn the surface into a string for testing purposes.
  String() string
}

//Whether a given point is on the interior of the surface. 
func SurfaceInterior(s Surface, x []float64) bool {
  return s.F(x) >= 0
}

//The normal vector to the surface. 
func SurfaceNormal(s Surface, x []float64) []float64 {
  grad := s.Gradient(x)
  var norm float64

  for i := 0; i < len(x); i ++ {
    norm += grad[i] * grad[i]
  }

  if norm == 0 {
    for i := 0; i < len(x); i ++ {
      grad[i] = 0.
    }
  } else {
    norm = math.Sqrt(norm)
    for i := 0; i < len(x); i ++ {
      grad[i] /= norm
    }
  }

  return grad
}

//A function to test gradients with numerical approximation against
//a surface's given formula. 
//For testing purposes only! 
//Uses a higher-order method. 
func testGradient(s Surface, x []float64, e float64) []float64 {
  z := make([]float64, s.Dimension())
  delta := make([]float64, s.Dimension())

  coefficients := []float64{3, -32, 168, -672, 0, 672, -168, 32, -3}

  for i := 0; i < s.Dimension(); i ++ {
    for j := -4; j <= 4; j ++ {
      for k := 0; k < s.Dimension(); k ++ {
        if i == k {
          delta[k] = x[k] + float64(j) * e
        } else {
          delta[k] = x[k]
        }
      }

      z[i] += coefficients[j + 4] * s.F(delta)
    }

    z[i] /= (840. * e)
  }

  return z
}

//A function to test intersections by numerical approximation with
//Newton's method against a surface's given formula. Assumes that
//the F method is correct. For testing purposes only! 
//x and x + v should be on different sides of the surface. The
//intersection is assumed to be between these two points. 
func testIntersection(s Surface, x []float64, v []float64, max_steps int) []float64 {
  //max_steps = 5

  var u0, u1, u2 float64 = 0.0, 1.0, 0.0
  p0 := make([]float64, len(x))
  p1 := make([]float64, len(x))
  p2 := make([]float64, len(x)) 
  
  for i := 0; i < len(x); i ++ {
    p0[i] = x[i] + v[i] * u0
    p1[i] = x[i] + v[i] * u1
  }

  f0 := s.F(p0)
  f1 := s.F(p1)
  var f2 float64

  var swap bool = false

  if f0 > 0 {
    if f1 > 0 {
      return []float64{}
    } else {
      pswap := p0
      p0 = p1
      p1 = pswap
      fswap := f0
      f0 = f1
      f1 = fswap
      swap = true
    }
  } else if f1 < 0 {
    return []float64{}
  }

  for i := 0; i < max_steps; i ++ {
    if u0 == u1 {return []float64{u0}}

    u2 = u0 + (u1 - u0) * f0 / (f0 - f1)

    for j := 0; j < len(x); j ++ {
      p2[j] = p0[j] + u2 * v[j]
    }

    f2 = s.F(p2)

    if f2 < 0 {
      u0 = u2
      f0 = f2
    } else {
      u1 = u2
      f1 = f2
    }
  }

  if swap {
    return []float64{1 - u2}
  } else {
    return []float64{u2}
  }
}
