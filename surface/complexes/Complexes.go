package complexes

//import "github.com/DanielKrawisz/CurvedSpace/vector"
import "github.com/DanielKrawisz/CurvedSpace/surface/"
import "github.com/DanielKrawisz/CurvedSpace/surface/polynomialsurfaces"
import "github.com/DanielKrawisz/CurvedSpace/surface/booleans"

//A simplex given as a list of points. This should
//be an n * (n + 1) matrix. 
//Note: the simplex will be inside-out if the points
//are not given in the right order. 
func NewSimplex(p [][]float64) surface.Surface {
  if p == nil { return nil }
  l := len(p)
  dim := l - 1
  if dim < 1 { return nil }

  var s, g surface.Surface = nil, nil

  p_sub := make([][]float64, dim)
  var sig int = 1
  for i := 0; i < l; i ++ {
    q := 0
    for j := (i + 1) % l; j != i; j = (j + 1) % l {
      //TODO This step may not get all points in the correct order.
      //Needs to be tested to make sure. 
      p_sub[q] = p[j]
      q ++ 
    }

    g = polynomialsurfaces.NewPlaneByPointsAndSignature(p_sub, 2 * ((sig + dim) % 2) - 1)
    sig += dim
    if g == nil {
      return nil
    }
    if g.Dimension() != dim {
      return nil
    }
    if s == nil {
      s = g
    } else {
      s = booleans.NewIntersection(g, s)
    }
  }

  return s
}

//A parallelpiped is given here by a corner point and 
//an n * n matrix. 
func NewParallelpipedByCornerAndEdges(P []float64, V [][]float64, right_side_out bool) surface.Surface {
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
  //halve the vectors and move the center to the correct place.
  for i := 0; i < dim; i ++ {
    v[i] = make([]float64, dim) 
    p[i] = P[i]
    for j := 0; j < dim; j ++ {
      v[i][j] = V[i][j] / 2
      p[i] += v[i][j]
    }
  }

  var inside float64
  if right_side_out { inside = -1 } else { inside = 1 }

  //TODO get this working again. 
  /*m := vector.Inverse(vector.Transpose(v))

  var s, g surface.Surface = nil, nil

  for i := 0; i < dim; i ++ {
    c := make([][]float64, dim)
    b := make([]float64, dim)
    for j := 0; j < dim; j ++ {
      c[j] = make([]float64, j + 1)
      for k := 0; k <= j; k ++ {
        if j == k && j == i {
          c[j][k] = inside
        } else {
          c[j][k] = 0
        }
      }
    }

    q := &quadraticSurface{dim, c, b, -inside}
    if q == nil {
      return nil
    }
    coordinateShiftQuadratic(q, m)
    translateQuadratic(q, p)
    g = q

    if s == nil {
      s = g
    } else {
      s = booleans.NewIntersection(g, s)
    }
  }

  return s*/

  return nil
}
