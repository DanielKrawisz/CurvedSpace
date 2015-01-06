package surface

import "testing"
import "../test"
import "../vector"

//For both these next two, I only have to test that
//a given point is interior or exterior to the complex.
func TestNewSimplex(t *testing.T) {
  //TODO
}

func TestNewParallelpiped(t *testing.T) {
  test.SetSeed(89999)

  var dim int = 2
  for i := 0; i < 5; i ++ {
    //p := test.RandFloatVector(-1, 1, dim)
    p := make([]float64, dim)
    v := make([][]float64, dim)
    for j := 0; j < dim; j ++ {
      v[j] = test.RandFloatVector(-.4, .4, dim)
    }
    for j := 0; j < dim; j ++ {
      v[j][j] += 2
    }

    pp := NewParallelpipedByCornerAndEdges(p, v)

    m := vector.Inverse(vector.Transpose(v))

    //5 test points. 
    for j := 0; j < 5; j ++ {
      point := test.RandFloatVector(-.5, 2.5, dim)
      test_p := make([]float64, dim)
      for k := 0; k < dim; k ++ {
        test_p[k] = point[k] - p[k]
      }
      test_inside := SurfaceInterior(pp, point)

      inverse := vector.MatrixMultiply(m, test_p)

      var inside bool = true
      for k := 0; k < dim; k ++ {
        inside = inside && inverse[k] > 0 && inverse[k] < 1
      }
      if inside != test_inside {
        t.Error("Parallelpiped error! ", pp.String(), " corner = ", p, ", v = ",
          v, ", point = ", point, ", inverse = ", m, " inverse vec = ", inverse, ", inside = ", test_inside, ", expected = ", inside)
      }
    }
  }
}
