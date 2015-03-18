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
      grad[i] = -grad[i] / norm
    }
  }

  return grad
}
