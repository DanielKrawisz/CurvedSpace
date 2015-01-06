package surface

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
      //TODO This step may not get all points in the correct order.
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
func NewParallelpipedByCornerAndEdges(P []float64, V [][]float64) Surface {
  if P == nil || V == nil {
    return nil
  }
  dim := len(P)
  if len(V) != dim {
    return nil
  }
  for i := 0; i < dim; i ++ {
    if V[i] == nil { return nil }
    if len(V[i]) != dim { return nil }
  }

  p := make([]float64, dim)
  v := make([][]float64, dim)
  //halve the vectors. 
  for i := 0; i < dim; i ++ {
    v[i] = make([]float64, dim) 
    for j := 0; j < dim; j ++ {
      v[i][j] = V[i][j] / 2
    }
  }
  //Move the position to the center of the object.
  for i := 0; i < dim; i ++ {
    for j := 0; j < dim; j ++ {
      p[i] += v[j][i]
    }
  }

  var s, g Surface = nil, nil

  for i := 0; i < dim; i ++ {

    g = NewInfiniteCylinder(p, [][]float64{v[i]})
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


