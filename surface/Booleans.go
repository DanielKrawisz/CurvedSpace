package surface

import "math"

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

//Assumes that a and b have been checked to have the 
//same dimension when the object was created. 
func (s *addition) Dimension() int {
  return s.a.Dimension()
}

func (s *addition) F(x []float64) float64 {
  return math.Max(s.a.F(x), s.b.F(x))
}

func (s *addition) Interior(x []float64) bool {
  if s.a.Interior(x) {
    return true
  } else if s.b.Interior(x) {
    return true
  }
  return false
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

    if !s.b.Interior(p) {
      z[zi] = inta[i]
      zi ++
    }
  }

  for i := 0; i < len(intb); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + intb[i] * v[j]
    }

    if !s.a.Interior(p) {
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

//Assumes that a and b have been checked to have the 
//same dimension when the object was created. 
func (s *intersection) Dimension() int {
  return s.a.Dimension()
}

func (s *intersection) F(x []float64) float64 {
  return math.Min(s.a.F(x), s.b.F(x))
}

func (s *intersection) Interior(x []float64) bool {
  if !s.a.Interior(x) {
    return false
  } else if !s.b.Interior(x) {
    return false
  }
  return true
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
  intb := s.b.Intersection(x, v)

  z := make([]float64, len(inta) + len(intb))

  var zi int = 0
  p := make([]float64, len(x))

  for i := 0; i < len(inta); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + inta[i] * v[j]
    }

    if s.b.Interior(p) {
      z[zi] = inta[i]
      zi ++
    }
  }

  for i := 0; i < len(intb); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + intb[i] * v[j]
    }

    if s.a.Interior(p) {
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

//Assumes that a and b have been checked to have the 
//same dimension when the object was created. 
func (s *subtraction) Dimension() int {
  return s.a.Dimension()
}

func (s *subtraction) F(x []float64) float64 {
  return math.Min(s.a.F(x), -s.b.F(x))
}

func (s *subtraction) Interior(x []float64) bool {
  if s.b.Interior(x) {
    return false
  } else if s.a.Interior(x) {
    return true
  }
  return false
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

    if !s.b.Interior(p) {
      z[zi] = inta[i]
      zi ++
    }
  }

  for i := 0; i < len(intb); i++ {
    for j := 0; j < len(x); j++ {
      p[j] = x[j] + intb[i] * v[j]
    }

    if s.a.Interior(p) {
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
