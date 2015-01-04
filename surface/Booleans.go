package surface

import "math"
import "strings"

type Boolean interface {
  Surface
  SurfaceA() Surface
  SurfaceB() Surface
}

type addition struct {
  a, b Surface 
}

func (s *addition) SurfaceA() Surface {
  return s.a 
}

func (s *addition) SurfaceB() Surface {
  return s.b 
}

func (s *addition) String() string {
  return strings.Join([]string{"addition{", s.a.String(), ", ", s.b.String(), "}"}, "")
}

//Assumes that a and b have been checked to have the 
//same dimension when the object was created. 
func (s *addition) Dimension() int {
  return s.a.Dimension()
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

type intersection struct {
  a, b Surface 
}

func (s *intersection) SurfaceA() Surface {
  return s.a 
}

func (s *intersection) SurfaceB() Surface {
  return s.b 
}

func (s *intersection) String() string {
  return strings.Join([]string{"intersection{", s.a.String(), ", ", s.b.String(), "}"}, "")
}

//Assumes that a and b have been checked to have the 
//same dimension when the object was created. 
func (s *intersection) Dimension() int {
  return s.a.Dimension()
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

func (s *intersection) Intersection(x, v []float64) []float64 {
  inta := s.a.Intersection(x, v)
  if len(inta) == 0 { return inta }
  intb := s.b.Intersection(x, v)
  if len(intb) == 0 { return intb }

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

type subtraction struct {
  a, b Surface 
}

func (s *subtraction) SurfaceA() Surface {
  return s.a 
}

func (s *subtraction) SurfaceB() Surface {
  return s.b 
}

func (s *subtraction) String() string {
  return strings.Join([]string{"subtraction{", s.a.String(), ", ", s.b.String(), "}"}, "")
}

//Assumes that a and b have been checked to have the 
//same dimension when the object was created. 
func (s *subtraction) Dimension() int {
  return s.a.Dimension()
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
  if len(inta) == 0 { return inta }
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

  return &addition{a, b}
}

//Can return nil
func NewSubtraction(a, b Surface) Surface {
  if a == nil || b == nil {return nil}
  if a.Dimension() != b.Dimension() {return nil}

  return &subtraction{a, b}
}

//Can return nil
func NewIntersection(a, b Surface) Surface {
  if a == nil || b == nil {return nil}
  if a.Dimension() != b.Dimension() {return nil}

  return &intersection{a, b}
}
