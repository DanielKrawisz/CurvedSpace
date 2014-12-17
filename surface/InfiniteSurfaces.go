package surface

//TODO: test
//A general quadratic surface from a central point and a list of vectors
//defining a quadratic form on the coordinates. The vectors do not need
//to satisfy any particular properties, but a set which is all normal
//to one another is the most general possibility. 
//
// (using some mixed notation here, but it works.)
//  (x - p) (v_i v_i) (x - p) == 0
//  p.v_i v_i.p - x.v_i v_i.p + x.v_i v_i.x == 0
//
//May return nil
func NewQuadraticSurfaceByCenterVectorList(p []float64, v [][]float64, r2 float64) Surface {
  if len(p) != len(v) {
    return nil
  }
  dim := len(p)
  for i := 0; i < dim; i++ {
    if len(v[i]) != dim {
      return nil
    }
  }

  var a float64 = r2
  var b []float64 = make([]float64, len(p))
  var c [][]float64 = make([][]float64, len(p))

  var vp []float64 = make([]float64, len(p)) 

  for i := 0; i < len(p); i ++ {
    vp[i] = 0
    for j := 0; j < len(p); j ++ {
      vp[i] += v[i][j] * p[j]
    }
  }

  for i := 0; i < len(p); i ++ {
    c[i] = make([]float64, i + 1)
    for j := 0; j <= i; j ++ {
      c[i][j] = 0.0
      for k := 0; k < len(p); j ++ {
        c[i][j] += v[k][i] * v[k][j]
      }
    }
  }

  return &quadraticSurface{dim, c, b, a}
}

//A surface that consists of two parallel planes.
//May return nil. 
func NewBiplane(point []float64, vector []float64, dist float64) Surface {
  if point == nil || vector == nil {
    return nil
  }
  if len(point) != len(vector) {
    return nil
  }
  dim := len(point)
  v := make([][]float64, dim) 
  v[0] = vector
  for i := 1; i < dim; i ++ {
    v[i] = make([]float64, dim)
  }

  return NewQuadraticSurfaceByCenterVectorList(point, v, dist*dist)
}

//A surface that consists of two parallel planes.
//May return nil. 
/*func NewInfiniteCylinder(point []float64, vector []float64, dist float64) Surface {
  
}

//Vectors are made to be orthonormal. 
func NewInfinite2SheetHyperboloid(point []float64, vector []float64) Surface {
  
}

//Given by the point at the center and the 
func NewInfinite1SheetHyperboloid(p []float64, a, b, c []float64) Surface {
  
}

//Given by the point at the apex of the paraboloid, 
//a set of vectors defining the symmetric tensor and a set
//defining the vector part of the quadratic surface. 
//May return nil.
func NewInfiniteParaboloid(p []float64, vc [][]float64, vb [][]float64) Surface {
  if p == nil || vc == nil || vb == nil {
    return nil
  }
  if len(p) != len(vc) + len(vb) {
    return nil
  }

  
}

func NewInfiniteCone(p []float64, vp [][]float64, vn [][]float64) Surface {
  if p == nil || vp == nil || vn == nil {
    return nil
  }
  if len(p) != len(vp) + len(vn) {
    return nil
  }

  dim := len(p)
  for i := 0; i < len(vp); i ++ {
    if len(vp[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vn); i ++ {
    if len(vn[i]) != dim {
      return nil
    }
  }

  
}*/
