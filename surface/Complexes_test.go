package surface

import "testing"
import "../test"
import "../vector"

//For both these next two, I only have to test that
//a given point is interior or exterior to the complex.
func TestNewSimplex(t *testing.T) {
  /*if NewSimplex(nil) != nil {
    t.Error("simplex error 1")
  }
  if NewSimplex([][]float64{}) != nil {
    t.Error("simplex error 2")
  }
  if NewSimplex([][]float64{[]float64{}}) != nil {
    t.Error("simplex error 3")
  }
  if NewSimplex([][]float64{[]float64{1}}) != nil {
    t.Error("simplex error 4")
  }
  if NewSimplex([][]float64{[]float64{1}, []float64{2}}) == nil {
    t.Error("simplex error 5")
  }
  if NewSimplex([][]float64{[]float64{0, 0}, []float64{0, 0}, []float64{0, 1}}) != nil {
    t.Error("simplex error 6")
  }
  if NewSimplex([][]float64{[]float64{0, 0}, []float64{1, 0}, []float64{0, 1, 0}}) != nil {
    t.Error("simplex error 7")
  }*/

  //Test whether the interior of the simplex is correct.
  //Go from dimension zero to dimension 5. 
  for dim := 1; dim < 5; dim ++ {
    //create the list of points. 
    points := make([][]float64, dim + 1)

    points[0] = make([]float64, dim)
    for j := 1; j <= dim; j ++ {
      points[j] = make([]float64, dim)
      points[j][j - 1] = 1
    }

    simplex := NewSimplex(points)

    if simplex == nil {
      t.Error("simplex error 8, dimension", dim)
    } else {
      inside := make([]float64, dim)
      for i := 0; i < dim; i ++ {
        inside[i] = 1./(float64(dim) + 1.)
      }

      //Ensure the center point is inside the simplex.
      if !SurfaceInterior(simplex, inside) {
        t.Error("simplex error 9, dimension", dim)
      }

      //Test to ensure the centers of all faces are on the surface.
      for i := 0; i < dim; i ++ {
        for j := 0; j < dim; j ++ {
          inside[j] = 1./float64(dim)
        }
        inside[i] = 0

        if !test.CloseEnough(simplex.F(inside), 0, .000001) {
          t.Error("simplex error 10, dimension ", dim, " point ", inside)
        }
      }

      if dim > 1 {
        for i := 0; i < dim; i ++ {
          inside[i] = 1./float64(dim)
        }

        if !test.CloseEnough(simplex.F(inside), 0, .000001) {
          t.Error("simplex error 10, dimension ", dim, " point ", inside, " got ", simplex.F(inside))
        }
      }
    }
  }
}

func TestNewParallelpiped(t *testing.T) {
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

    m := vector.Inverse(v)

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
