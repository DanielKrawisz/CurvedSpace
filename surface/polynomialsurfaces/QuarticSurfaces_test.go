package polynomialsurfaces

import "testing"
import "github.com/DanielKrawisz/CurvedSpace/distributions"
import "github.com/DanielKrawisz/CurvedSpace/test"
import "github.com/DanielKrawisz/CurvedSpace/vector"
import "github.com/DanielKrawisz/CurvedSpace/surface"

//Strategy: just pick a bunch of rays somewhat randomly and test that the
//intersection points which are returned obey the torus equation. 
func TestTorus(t *testing.T) {
  if NewTorus(nil, []float64{0, 0, 1}, 3, 1) != nil { t.Error("torus error 1") }
  if NewTorus([]float64{0, 0, 0}, nil, 3, 1) != nil { t.Error("torus error 2") }
  if NewTorus([]float64{0, 0}, []float64{0, 0, 1}, 3, 1) != nil { t.Error("torus error 3") }
  if NewTorus([]float64{0, 0, 0}, []float64{0, 0}, 3, 1) != nil { t.Error("torus error 4") }
  if NewTorus([]float64{0, 0, 0, 0}, []float64{0, 0, 1}, 3, 1) != nil { t.Error("torus error 5") }
  if NewTorus([]float64{0, 0, 0}, []float64{0, 0, 1, 0}, 3, 1) != nil { t.Error("torus error 6") }

  //TODO Only can do three dimensions until the unit sphere function
  //is generalized to more dimensions. 
  for dim := 3; dim <= 3; dim ++ {
    //Test 3 cases in each dimension. 
    for i := 0; i < 3; i ++ { 
      //We already know we can translate polynomial surfaces, so
      //can leave this at zero. 
      zero := make([]float64, dim)

      z := distributions.RandomUnitSphereSurfacePoint()[:]

      R := test.RandFloat(3, 4)
      r := test.RandFloat(1, 2)

      torus := NewTorus(zero, z, R, r)

      for j := 0; j < 9; j ++ { 
        p := vector.LinearSum(1, R * 1.1, zero, distributions.RandomUnitSphereSurfacePoint()[:])
        v := vector.LinearSum(-1, r * 1.2, p, distributions.RandomUnitSphereSurfacePoint()[:])

        us := torus.Intersection(p, v)

        V := make([][]float64, dim)
        for a := 0; a < dim; a ++ {
          V[a] = make([]float64, a + 1)
          for b := 0; b < a; b ++ {
            V[a][b] = - z[a]*z[b]
          }
          V[a][a] = 1 - z[a]*z[a]
        }

        for _, u := range us {

          x := vector.LinearSum(1, u, p, v) 
          t1 := vector.Dot(x, x) + R*R - r*r
          F := - t1*t1 + 4 * R*R * vector.Dot(vector.ContractSymmetricTensor(V, x), x)

          if !test.CloseEnough(F, 0, .000001) {
            t.Error("torus error! got ", F, " for ", torus, " at point ", x)
          }
        }
      }
    }
  }
}

//Test some interior and exterior points. 
func TestTorusInterior(t *testing.T) {
  p := []float64{0, 0, 0}
  v := []float64{0, 0, 1}
  var R, r float64 = 2, 1

  torus := NewTorus(p, v, R, r)

  if torus == nil {
    t.Error("torus error: torus should exist but is nil!!!!") 
    return 
  }

  if surface.SurfaceInterior(torus, p) {
    t.Error("torus error: torus is inside-out 1! ", torus.F(p)) 
  }

  if !surface.SurfaceInterior(torus, []float64{2, 0, 0}) {
    t.Error("torus error: torus is inside-out 2! ", torus.F([]float64{2, 0, 0})) 
  }
}
