package surface

import "math"
import "errors"
//import "fmt"
import "../vector"

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
//the F method is correct. 
//x and x + v should be on different sides of the surface. The
//intersection is assumed to be between these two points. 
//TODO improve this function for use in isosurfaces eventually. 
//There should be a way to make this work more like the original
//version and I bet it would be quicker. 
func testIntersection(s Surface, x []float64, v []float64, tolerance float64, max_steps int) ([]float64, error) {
  var err error = nil

//  fmt.Println("Running intersection test for ", s.String(), " with (p, v) = ", x, v)

  var u float64 = 0.0
  p0 := make([]float64, len(x))
  p1 := make([]float64, len(x))
  p := make([]float64, len(x)) 
  
  for i := 0; i < len(x); i ++ {
    p0[i] = x[i]
    p1[i] = x[i] + v[i]
  }

  f0 := s.F(p0)
  f1 := s.F(p1)
  var f, f_last float64

  u = f0 / (f0 - f1) //Estimate a good initial value for u.
  f_last = math.Inf(1)

  //Use Newton's method to find the intersection. 
  for i := 0; i < max_steps; i ++ {

    //New estimated intersection point.
    for j := 0; j < len(x); j ++ {
      p[j] = x[j] + u * v[j]
    }

    f_last = f
    f = s.F(p)

//    fmt.Println("step ", i, "; u = ", u, "; p = ", p, "; f = ", f)

    if f == f_last {
      goto convergence
    }

//    fmt.Println("step ", i, "; grad = ", s.Gradient(p),
//      "; g.v = ", vector.Dot(s.Gradient(p), v), "; -f / g.v = ", -f / vector.Dot(s.Gradient(p), v))

    //new u calculated with Lie derivative wrt to v at p.
    u -= f / vector.Dot(s.Gradient(p), v)
  }

  if math.Abs(f_last - f) > tolerance {
    err = errors.New("testIntersection: max_steps reached!")
  }

  convergence:

  return []float64{u}, err
}
