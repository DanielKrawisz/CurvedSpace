package diffeq

type CompoundDerivative interface {
  Plus(Derivative) 
  Length() int
  Derivative
}

type compoundDerivative struct {
  dim int
  l []Derivative
}

func (cd *compoundDerivative) Dimension() int {
  return cd.dim
}

func (cd *compoundDerivative) Length() int {
  return len(cd.l)
}

func (cd *compoundDerivative) Plus(f Derivative) {
  if f == nil {
    return
  }
  if f.Dimension() != cd.dim {
    return
  }

  for i := 0; i < len(cd.l); i ++ {
    if cd.l[i] == f {
      cd.l = append(cd.l[:i], cd.l[i+1:]...)
      return 
    }
  }

  cd.l = append(cd.l, f)
}

func (cd *compoundDerivative) DxDs(x []float64, v []float64) {
  for i := 0; i < len(v); i ++ {
    v[i] = 0
  }
  v2 := make([]float64, len(v))
  for e := 0; e < len(cd.l); e ++ {
    for i := 0; i < len(v); i ++ {
      v2[i] = 0
    }

    cd.l[e].DxDs(x, v2)

    for i := 0; i < len(v); i ++ { 
      v[i] += v2[i]
    }
  }
}

func NewCompoundDerivative(dim int) CompoundDerivative {
  return &compoundDerivative{dim, []Derivative{}}
}
