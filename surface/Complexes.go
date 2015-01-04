package surface

import "../vector"

//A simplex given as a list of points. This should
//be an n * (n + 1) matrix. 
func NewSimplex(p [][]float64) Surface {
  if p == nil {return nil}
  dim := len(p) - 1
  if dim < 1 {return nil}

  var s, g Surface = nil, nil

  p_sub := make([][]float64, dim - 1)
  for i := 0; i < dim; i ++ {
    q := 0
    for j := 0; j < dim; j ++ {
      if j != i {
        p_sub[q] = p[j]
        q ++ 
      }
    }

    g = NewPlaneByPoints(p_sub)
    if g == nil {
      return nil
    }
    if s == nil {
      s = g
    } else {
      s = NewIntersection(g, s)
    }
  }

  return s
}

//A parallelpiped is given here by a corner point and 
//an n * n matrix. 
func NewParallelPipedByCenterAndEdges(p []float64, v [][]float64) Surface {
  if p == nil || v == nil {
    return nil
  }
  dim := len(p)
  if len(v) != dim {
    return nil
  }
  for i := 0; i < dim; i ++ {
    if v[i] == nil { return nil }
    if len(v[i]) != dim { return nil }
  }

  var s, g Surface = nil, nil

  v_sub := make([][]float64, dim - 1)
  for i := 0; i < dim; i ++ {
    q := 0
    for j := 0; j < dim; j ++ {
      if j != i {
        v_sub[q] = v[j]
        q ++ 
      }
    }

    g = NewInfiniteCylinder(p, [][]float64{vector.Cross(v_sub)})
    if g == nil {
      return nil
    }
    if s == nil {
      s = g
    } else {
      s = NewIntersection(g, s)
    }
  }

  return s
}
