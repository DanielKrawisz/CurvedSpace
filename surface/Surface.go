package surface

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
  //Whether a given point is on the interior of the surface. 
  Interior(x []float64) bool
  //The gradianThe normal to a surface at a given point on the surface.
  Gradient(x []float64) []float64
}

//A function to test gradients with numerical approximation against
//a surface's given formula. Assumes that the F method is correct.
//For testing purposes only! 
//TODO need higher-order methods or more digits.
func testGradient(s Surface, x []float64, e float64) []float64 {
  z := make([]float64, s.Dimension())
  delta := make([]float64, s.Dimension())
  f := s.F(x)

  for i := 0; i < s.Dimension(); i ++ {
    for j := 0; j < s.Dimension(); j ++ {
      if i == j {
        delta[j] = x[j] + e
      } else {
        delta[j] = x[j]
      }
    }

    z[i] = (s.F(delta) - f) / e
  }

  return z
}

//A function to test intersections by numerical approximation with
//Newton's method against a surface's given formula. Assumes that
//the F method is correct. For testing purposes only! 
//x and x + v should be on different sides of the surface. The
//intersection is assumed to be between these two points. 
func testIntersection(s Surface, x []float64, v []float64, max_steps int) []float64 {
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

    u2 = u0 + (u0 - u1) * f0 / (f0 - f1)

    for j := 0; j < len(x); j ++ {
      p2[i] = p0[i] + u2 * p1[i]
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
    return []float64{1 - u0}
  } else {
    return []float64{u0}
  }
}
