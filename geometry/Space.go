package geometry

//There may be regions of space in which a particular coordinate system is invalid. 
//This error is returned when that happens. 
type InvalidCoordinateError struct {
}

func (e *InvalidCoordinateError) Error() string {
  return "Invalid coordinates."
}

//A Space is a pseudo-Riemannian manifold. 
//Because different regions of the space may be most conveniently
//metrized with different coordinate systems, the Space struct is
//mainly concerned with managing a set of coordinate systems.
// 
//A space is also a Derivative for a particle that parallel
//transports itself along its surface. 
type Space interface {
  //The space is divided into one or more regions, which may overlap.
  //Each is given a name.
  RegionName(c int) *string

  //The coordinate system that is used for a particular region of space.
  CoordinateSystem(c int) CoordinateSystem

  //Transform a coordinate point which is represented in the coordinates
  //of one region to that of another. This is only allowed when the two
  //regions overlap at the given point. 
  TransformCoordinates(x CoordinatePoint, to int) (*InvalidCoordinateError)

  //At the CoordinatePoint x, which is the region that the space prefers to use?
  PreferredRegion(x CoordinatePoint) int

  //Transforms coordinate point to the preferred region for its location. 
  //No effect if already in the preferred region. 
  TransformCoordinatesToPreferredRegion(x CoordinatePoint)

  //Creates a CoordinatePoint struct at a given location. Only allowed
  //when the given coordinates are valid in the given region. 
  CoordinatePoint(space *Space, t, x, y, z float64, region int) (CoordinatePoint, *InvalidCoordinateError) 
}
