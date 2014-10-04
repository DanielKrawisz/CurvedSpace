package geometry

type CoordinateSystem interface {
  Name() string
  Metric(p *CoordinatePoint)
}

type CoordinatePoint interface {
  Space() Space
  CoordinateSystem() CoordinateSystem
  Region() int
  RegionName() *string
  Metric() Metric
  T() float64
  X() float64
  Y() float64
  Z() float64
  x() []float64
}

type coordinatePoint struct {
  space Space
  region int
  //metric Metric
  point []float64
}

func (c *coordinatePoint) Space() Space {
  return c.space
}

func (c *coordinatePoint) Region() int {
  return c.region
}

func (c *coordinatePoint) RegionName() *string {
  return c.space.RegionName(c.region)
}

func (c *coordinatePoint) CoordinateSystem() CoordinateSystem {
  return c.space.CoordinateSystem(c.region)
}

func (c *coordinatePoint) T() float64 {
  return c.point[0]
}

func (c *coordinatePoint) X() float64 {
  return c.point[1]
}

func (c *coordinatePoint) Y() float64 {
  return c.point[2]
}

func (c *coordinatePoint) Z() float64 {
  return c.point[3]
}

func (c *coordinatePoint) x() []float64 {
  return c.point
}

func NewCoordinatePoint(space Space, region int, x []float64) (*coordinatePoint, *InvalidCoordinateError) {
  //TODO: check that the region is good, do the coordinate transformation
  return nil, nil
}

type Vector interface {
  Dt() float64
  Dx() float64
  Dy() float64
  Dz() float64
  Norm() float64
  Location() CoordinatePoint
  dx() []float64
}

type vector struct {
  space Space
  region int
  metric Metric
}

func (c *vector) Space() Space {
  return c.space
}

func (c *vector) Region() int {
  return c.region
}

func (c *vector) RegionName() *string {
  return c.space.RegionName(c.region)
}

func (c *vector) Metric() Metric {
  return c.metric
}

func (v *vector) Norm() float64 {
  return v.Metric().InnerProduct(v, v)
}

func (v *vector) Location() CoordinatePoint {
  //TODO -- it should create a new coordinate point for its location.
  return nil;
}

func (v *vector) Dt() float64 {
  //TODO
  return 0
}

func (v *vector) Dx() float64 {
  //TODO
  return 0
}

func (v *vector) Dy() float64 {
  //TODO
  return 0
}

func (v *vector) Dz() float64 {
  //TODO
  return 0
}

func (v *vector) dx() []float64 {
  //TODO
  return nil
}

type CoordinateTransformation interface {
  From() CoordinateSystem
  To() CoordinateSystem
  Transform(v Vector)
}

type Metric interface {
  InnerProduct(v, w Vector) float64
}
