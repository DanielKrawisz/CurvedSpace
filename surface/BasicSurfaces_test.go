package surface

import "testing"
import "../test"
import "sort"

var err_bs float64 = 0.00001

//TODO a lot of tests! 
//TODO intersection tests

func TestNewSphere(t *testing.T) {
  if NewSphere(nil, 1) != nil {
    t.Error("sphere position cannot be nil.")
  }

  s1 := NewSphere([]float64{1,2,3}, 2)
  s2 := NewSphere([]float64{1,2,3,4}, -2) 

  if s1 == nil || s2 == nil {
    t.Error("nil returned for valid values.")
  }

  if s1.Dimension() != 3 || s2.Dimension() != 4 {
    t.Error("Invalid values returned for sphere dimensions")
  }

  if !test.CloseEnough(s2.R2(), 4, err_bs) {
    t.Error("Invalid sphere value 1: got ", s2.R2(), ", expected 4.")
  }  

  if !test.CloseEnough(s1.R2(), 4, err_bs) {
    t.Error("Invalid sphere value 2.")
  }  

  if !test.CloseEnough(s1.X()[0], 1, err_bs) {
    t.Error("Invalid sphere value 3.")
  }  

  if !test.CloseEnough(s1.X()[1], 2, err_bs) {
    t.Error("Invalid sphere value 4.")
  }  
}

func TestSphereF(t *testing.T) {
  for i := 0; i < 10; i++ {
    p := make([]float64, 2)

    for j := 0; j < 2; j ++ {
      p[j] = test.RandFloat(-3, 3)
    }
    r := test.RandFloat(1, 5)

    s := NewSphere(p, r)

    test_point := make([]float64, 2)
    //Test 5 random points.
    for j := 0; j < 5; j++ {
      test_point[0] = test.RandFloat(-5, 5)
      test_point[1] = test.RandFloat(-5, 5)

      f_exp := r*r - (p[0] - test_point[0]) * (p[0] - test_point[0]) - (p[1] - test_point[1]) * (p[1] - test_point[1])
      f := s.F(test_point)
      grad_exp := []float64{-2 * (p[0] - test_point[0]), -2 * (p[1] - test_point[1])}
      grad := s.Gradient(test_point)

      if !test.CloseEnough(f_exp, f, err_bs) {
        t.Error("Sphere error p = ", p, ", r = ", r,
          ". F error at point ", test_point, ". Expected ", f_exp, "; got ", f)
      }

      if (f > 0) != (SurfaceInterior(s, test_point)) {
        t.Error("Sphere error p = ", p, ", r = ", r, ". Interior error at point ", test_point)
      }

      if !test.CloseEnough(grad[0], grad_exp[0], err_bs) && !test.CloseEnough(grad[1], grad_exp[1], err_bs) {
        t.Error("Sphere error p = ", p, ", r = ", r,
          ". Grad error at point ", test_point, ". Expected ", grad_exp, "; got ", grad)
      }
    }
  }
}

//Tests some preset cases for sphere intersection. 
func TestSphereIntersectionSetCases(t *testing.T) {
  //These values were generated with the Mathematica code:
  //
  //spherex = {0, 1.4}
  //spherey = {0, .6}
  //spherer = {1, 2}
  //raypos = {{-3, 0}, {-.5, 0}}
  //raydir = {{1, .5}, {1, .9}, {.4, .9}}
  //f[n_] := Function[
  //  StringJoin[(Sequence @@ Table["[]", {n}]), "float64{", 
  //   Riffle[#, ", "], "}"]]
  //Fold[Function[{a, b}, Map[f[6 - b], a, {b}]], 
  // Outer[Function[{sx, sy, sr, rp, rd}, 
  //   If[#[[1]] \[Element] Reals, ToString /@ #, {}] &[
  //    Last[Last[#]] & /@ 
  //     Solve[#.# &[rp + t rd - {sx, sy}] == sr^2, t]]], spherex, 
  //  spherey, spherer, raypos, raydir, 1], {5, 4, 3, 2, 1, 0}]
  //

  spherex := []float64{0, 1.4}
  spherey := []float64{0, .6}
  spherer := []float64{1, 2}
  raypos := [][]float64{[]float64{-3, 0}, []float64{-.5,0}}
  raydir := [][]float64{[]float64{1, .5}, []float64{1, .9}, []float64{.4, .9}}
  index := make([]int, 5)

  answers := [][][][][][]float64{[][][][][]float64{[][][][]float64{[][][]float64{[][]float64{[]float64{}, []float64{}, []float64{}}, [][]float64{[]float64{-0.47178, 1.27178}, []float64{-0.424239, 0.976725}, []float64{-0.69698, 1.10935}}}, [][][]float64{[][]float64{[]float64{1.07335, 3.72665}, []float64{}, []float64{}}, [][]float64{[]float64{-1.37764, 2.17764}, []float64{-1.18941, 1.74189}, []float64{-1.77081, 2.18318}}}}, [][][][]float64{[][][]float64{[][]float64{[]float64{2.10934, 3.17066}, []float64{}, []float64{}}, [][]float64{[]float64{-0.20947, 1.48947}, []float64{-0.164074, 1.31325}, []float64{-0.22911, 1.75488}}}, [][][]float64{[][]float64{[]float64{1.00244, 4.27756}, []float64{1.02638, 2.88523}, []float64{}}, [][]float64{[]float64{-1.12681, 2.40681}, []float64{-0.909691, 2.05886}, []float64{-1.25623, 2.78201}}}}}, [][][][][]float64{[][][][]float64{[][][]float64{[][]float64{[]float64{}, []float64{}, []float64{}}, [][]float64{[]float64{1.04841, 1.99159}, []float64{}, []float64{}}}, [][][]float64{[][]float64{[]float64{3.2, 3.84}, []float64{}, []float64{}}, [][]float64{[]float64{-0.0993826, 3.13938}, []float64{-0.0980522, 2.1975}, []float64{-0.224434, 1.79144}}}}, [][][][]float64{[][][]float64{[][]float64{[]float64{}, []float64{}, []float64{}}, [][]float64{[]float64{0.91053, 2.60947}, []float64{0.928068, 1.76806}, []float64{}}}, [][][]float64{[][]float64{[]float64{2.51036, 5.00964}, []float64{}, []float64{}}, [][]float64{[]float64{-0.00680503, 3.52681}, []float64{-0.00613359, 2.70227}, []float64{-0.0114892, 2.6919}}}}}}

  var sx, sy, sr float64
  var rp, rd []float64
  for index[0], sx = range spherex {
    for index[1], sy = range spherey {
      for index[2], sr = range spherer {
        s := NewSphere([]float64{sx, sy}, sr)
        for index[3], rp = range raypos {
          for index[4], rd = range raydir {
            intersection := s.Intersection(rp, rd)
            sort.Float64s(intersection)
            expected := answers[index[0]][index[1]][index[2]][index[3]][index[4]]
            sort.Float64s(expected)
            if len(intersection) != len(expected) {
              t.Error("Test case ", index, " returned ", intersection, " expected ", expected)
            } else {
              var match bool = true

              for i := 0; i < len(intersection); i++ {
                match = match && test.CloseEnough(intersection[i], expected[i], .0001) 
              }

              if !match {
                t.Error("Test case ", index, " returned ", intersection, " expected ", expected)
              }
            }
          }
        }
      }
    }
  }
}

func TestSphereIntersection(t *testing.T) {
  for i := 0; i < 4; i ++ {
    point := []float64{test.RandFloat(-2, 2), test.RandFloat(-2, 2)}
    r := test.RandFloat(1, 2)
    sphere := NewSphere(point, r)

    for j := 0; j < 4; j ++ {
      var p1, p2 []float64

      for {
        p2 = []float64{test.RandFloat(point[0] - r, point[0] + r), test.RandFloat(point[1] - r, point[1] + r)}

        if (p2[0] - point[0])*(p2[0] - point[0]) + (p2[1] - point[1])*(p2[1] - point[1]) < r*r {
          break;
        }
      }

      for {
        p1 = []float64{test.RandFloat(-10, 10), test.RandFloat(-10, 10)}

        if (p1[0] - point[0])*(p1[0] - point[0]) + (p1[1] - point[1])*(p1[1] - point[1]) > r*r {
          break;
        }
      }

      v := make([]float64, len(p1))
      for i := 0; i < len(p1); i ++ {
        v[i] = p2[i] - p1[i]
      }

      u := sphere.Intersection(p1, v)
      u_test := testIntersection(sphere, p1, v, 100)

      if len(u) == 0 {
        t.Error("No intersection point found for ", sphere.String(), " p1 = ", p1, "; v = ", v)
        return 
      } 

      if len(u_test) == 0 {
        t.Error("No test intersection point found for ", sphere.String(), " p1 = ", p1, "; v = ", v)
        return
      }

      intersection_point := make([]float64, 2)

      close_enough_test := false
      close_enough_F := true
      f := make([]float64, len(u))
      for q, uu := range u {
        for i := 0; i < 2; i++ {
          intersection_point[i] = p1[i] + uu * v[i]
        }

        f[q] = sphere.F(intersection_point)

        if !test.CloseEnough(f[q], 0.0, err_bs) {
          close_enough_F = false
        }

        if test.CloseEnough(u_test[0], uu, err_bs) {
          close_enough_test = true
        }
      }

      if !close_enough_F {
        t.Error("sphere intersection error for ", sphere.String(), ", p1 = ",
          p1, "; v = ", v, "; u = ", u, "; F = ", f)
      }

      if !close_enough_test {
        t.Error("test point and intersection point do not agree: u = ", u, "; u_test = ", u_test)
      }
    }
  }
}

func TestNewEllipsoid(t *testing.T) {
  if nil != NewElipsoidByCenterBasis(nil, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 1")
  }
  if nil != NewElipsoidByCenterBasis([]float64{0,0}, nil, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 2")
  }
  if nil != NewElipsoidByCenterBasis([]float64{0,0}, [][]float64{nil, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 3")
  }
  if nil != NewElipsoidByCenterBasis([]float64{0,0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 4")
  }
  if nil != NewElipsoidByCenterBasis([]float64{0,0}, [][]float64{[]float64{1,0,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 5")
  }
  if nil != NewElipsoidByCenterBasis([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1,1}) {
    t.Error("New ellipsoid surface error 6")
  }
  if nil != NewElipsoidByCenterBasis([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 7")
  }

  if nil == NewElipsoidByCenterBasis([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 8")
  }
}

//No need to test other functions because we know from testing the
//general polynomial surfaces. 
func TestEllipsoidF(t *testing.T) {
  for i := 0; i < 5; i ++ {
    
  }
}

func TestNewInfiniteCylinder(t *testing.T) {
  if nil != NewInfiniteCylinder(nil, [][]float64{[]float64{1,0}, []float64{0, 1}}) {
    t.Error("New infinite cylinder error 1")
  }
  if nil != NewInfiniteCylinder([]float64{0,0}, [][]float64{[]float64{1,0}, nil}) {
    t.Error("New infinite cylinder error 2")
  }
  if nil != NewInfiniteCylinder([]float64{0,0}, nil) {
    t.Error("New infinite cylinder error 3")
  }

  if nil == NewInfiniteCylinder([]float64{0,0}, [][]float64{[]float64{1, 0}, []float64{0, 1}}) {
    t.Error("New infinite cylinder error 4")
  }
  if nil == NewInfiniteCylinder([]float64{0,0}, [][]float64{[]float64{1,0}}) {
    t.Error("New infinite cylinder error 5")
  }
}

func TestInfiniteCylinderF(t *testing.T) {
  
}

func TestNewInfiniteCone(t *testing.T) {
  
}

func TestInfiniteConeF(t *testing.T) {
  
}

func TestNewInfiniteParaboloid(t *testing.T) {
  
}

func TestInfiniteParaboloidF(t *testing.T) {
  
}

func TestNewInfiniteHyperboloid(t *testing.T) {
  
}

func TestInfiniteHyperboloidF(t *testing.T) {
  
}

func TestNewCylinder(t *testing.T) {
  
}

func TestCylinderF(t *testing.T) {
  
}
