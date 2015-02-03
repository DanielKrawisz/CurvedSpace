package surface

//Booleans allow several different objects to be merged
//into a single object in various ways. 

import "math"
import "strings"
import "github.com/DanielKrawisz/CurvedSpace/vector"

type Boolean interface {
  Surface
  SurfaceA() Surface
  SurfaceB() Surface
}

//Addition objects include the points from both surfaces. 
type boolean struct {
  a, b Surface 
}

func (s *boolean) SurfaceA() Surface {
  return s.a 
}

func (s *boolean) SurfaceB() Surface {
  return s.b 
}

//Assumes that a and b have been checked to have the 
//same dimension when the object was created. 
func (s *boolean) Dimension() int {
  return s.a.Dimension()
}

//Addition objects include the points from both surfaces. 
type addition struct {
  boolean 
}

func (s *addition) String() string {
  return strings.Join([]string{"addition{", s.a.String(), ", ", s.b.String(), "}"}, "")
}

func (s *addition) F(x []float64) float64 {
  return math.Max(s.a.F(x), s.b.F(x))
}

func (s *addition) Gradient(x []float64) []float64 {
  if s.a.F(x) >= s.b.F(x) {
    return s.a.Gradient(x)
  } else {
    return s.b.Gradient(x)
  }
}

func (s *addition) Intersection(x, v []float64) []float64 {
  inta := s.a.Intersection(x, v)
  intb := s.b.Intersection(x, v)

  z := make([]float64, len(inta) + len(intb))

  var zi int = 0
  p := make([]float64, len(x))

  for i := 0; i < len(inta); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + inta[i] * v[j]
    }

    if !SurfaceInterior(s.b, p) {
      z[zi] = inta[i]
      zi ++
    }
  }

  for i := 0; i < len(intb); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + intb[i] * v[j]
    }

    if !SurfaceInterior(s.a, p) {
      z[zi] = intb[i]
      zi ++
    }
  }

  return z[0:zi]
}

//Intersection objects include only the points that are
//common to both surfaces. 
type intersection struct {
  boolean
}

func (s *intersection) String() string {
  return strings.Join([]string{"intersection{", s.a.String(), ", ", s.b.String(), "}"}, "")
}

func (s *intersection) F(x []float64) float64 {
  return math.Min(s.a.F(x), s.b.F(x))
}

func (s *intersection) Gradient(x []float64) []float64 {
  if s.a.F(x) < s.b.F(x) {
    return s.a.Gradient(x)
  } else {
    return s.b.Gradient(x)
  }
}

//This is set off in its own function so that the bounding object can use it.
func (s *intersection) findCommonIntersectionPoints(x, v, inta, intb []float64) []float64 {
  z := make([]float64, len(inta) + len(intb))

  var zi int = 0
  p := make([]float64, len(x))

  for i := 0; i < len(inta); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + inta[i] * v[j]
    }

    if SurfaceInterior(s.b, p) {
      z[zi] = inta[i]
      zi ++
    }
  }

  for i := 0; i < len(intb); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + intb[i] * v[j]
    }

    if SurfaceInterior(s.a, p) {
      z[zi] = intb[i]
      zi ++
    }
  }

  return z[0:zi]
}

func (s *intersection) Intersection(x, v []float64) []float64 {
  inta := s.a.Intersection(x, v)
  intb := s.b.Intersection(x, v)

  return s.findCommonIntersectionPoints(x, v, inta, intb)
}

//Bounding objects are the same as intersection objects 
//EXCEPT that it is assumed that if there are no intersection
//points for object a, then it does not bother checking
//object b. This allows a to be a bounding box or bounding
//sphere. Something simple that defines the limits of a more
//complicated object. 
type bounding struct {
  intersection
}

func (s *bounding) Intersection(x, v []float64) []float64 {
  inta := s.a.Intersection(x, v)

  //The only difference between intersection and bounding right here.
  //This allows for object a to be a very simple bounding object that
  //entirely surrounds the real object and prevents excepssive 
  //computation when testing for intersections. 
  if len(inta) == 0 {
    return inta
  }

  intb := s.b.Intersection(x, v)

  return s.intersection.findCommonIntersectionPoints(x, v, inta, intb)
}

//Open bounding objects are the same as bounding objects 
//EXCEPT that only intersection points for object b are
//returned. This can be used for creating a misty object
//that is bounded by a surface but such that light rays
//only interact with the interior, not with the surface. 
type openBounding struct {
  intersection
}

func (s *openBounding) Intersection(x, v []float64) []float64 {
  intb := s.b.Intersection(x, v)

  z := make([]float64, len(intb))
  var zi int = 0

  for i, u := range intb {
    if SurfaceInterior(s.a, vector.LinearSum(1, u, x, v)) {
      z[i] = u
      zi ++
    }
  }

  return z[0:zi]
}

//Subtraction objects allow one object to be cut out of another.
type subtraction struct {
  boolean
}

func (s *subtraction) String() string {
  return strings.Join([]string{"subtraction{", s.a.String(), ", ", s.b.String(), "}"}, "")
}

func (s *subtraction) F(x []float64) float64 {
  return math.Min(s.a.F(x), -s.b.F(x))
}

func (s *subtraction) Gradient(x []float64) []float64 {
  if s.a.F(x) < -s.b.F(x) {
    return s.a.Gradient(x)
  } else {
    z := s.b.Gradient(x)
    for i := 0; i < len(z); i ++ { 
      z[i] = - z[i]
    }
    return z
  }
}

//The intersection boolean can also be used to make a bounding
//box or bounding sphere because the intersection function
//returns right away if there are not intersection points 
//for either subsurface. 
func (s *subtraction) Intersection(x, v []float64) []float64 {
  inta := s.a.Intersection(x, v)
  intb := s.b.Intersection(x, v)

  z := make([]float64, len(inta) + len(intb))

  var zi int = 0
  p := make([]float64, len(x))

  for i := 0; i < len(inta); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + inta[i] * v[j]
    }

    if !SurfaceInterior(s.b, p) {
      z[zi] = inta[i]
      zi ++
    }
  }

  for i := 0; i < len(intb); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + intb[i] * v[j]
    }

    if SurfaceInterior(s.a, p) {
      z[zi] = intb[i]
      zi ++
    }
  }

  return z[0:zi]
}

//Can return nil
func NewAddition(a, b Surface) Surface {
  if a == nil || b == nil {return nil}
  if a.Dimension() != b.Dimension() {return nil}

  return &addition{boolean{a, b}}
}

//Can return nil
func NewSubtraction(a, b Surface) Surface {
  if a == nil || b == nil {return nil}
  if a.Dimension() != b.Dimension() {return nil}

  return &subtraction{boolean{a, b}}
}

//Can return nil
func NewIntersection(a, b Surface) Surface {
  if a == nil || b == nil {return nil}
  if a.Dimension() != b.Dimension() {return nil}

  return &intersection{boolean{a, b}}
}

//Can return nil
func NewBounding(a, b Surface) Surface {
  if a == nil || b == nil {return nil}
  if a.Dimension() != b.Dimension() {return nil}

  return &bounding{intersection{boolean{a, b}}}
}

//Can return nil
func NewOpenBounding(a, b Surface) Surface {
  if a == nil || b == nil {return nil}
  if a.Dimension() != b.Dimension() {return nil}

  return &openBounding{intersection{boolean{a, b}}}
}
