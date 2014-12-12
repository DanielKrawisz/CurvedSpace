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

//TODO: create functions to test gradients and intersections. 
