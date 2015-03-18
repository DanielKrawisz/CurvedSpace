package polynomialsurfaces

import "testing"
import "github.com/DanielKrawisz/CurvedSpace/test"
import "github.com/DanielKrawisz/CurvedSpace/surface"
import "sort"
import "math"

var err_bs float64 = 0.00001

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
    p := test.RandFloatVector(-3, 3, 2)
    r := test.RandFloat(1, 5)

    s := NewSphere(p, r)

    var test_point []float64
    //Test 5 random points.
    for j := 0; j < 5; j++ {
      test_point = test.RandFloatVector(-5, 5, 2)

      f_exp := r*r - (p[0] - test_point[0]) * (p[0] - test_point[0]) - (p[1] - test_point[1]) * (p[1] - test_point[1])
      f := s.F(test_point)
      grad_exp := []float64{2 * (p[0] - test_point[0]), 2 * (p[1] - test_point[1])}
      grad := s.Gradient(test_point)

      if !test.CloseEnough(f_exp, f, err_bs) {
        t.Error("Sphere error p = ", p, ", r = ", r,
          ". F error at point ", test_point, ". Expected ", f_exp, "; got ", f)
      }

      if (f > 0) != (surface.SurfaceInterior(s, test_point)) {
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

    //four random intersection tests.
    for j := 0; j < 4; j ++ {
      var p1, p2 []float64

      //Generate some valid parameters for intersection lines. 
      for {
        p2 = []float64{test.RandFloat(point[0] - r, point[0] + r),
                 test.RandFloat(point[1] - r, point[1] + r)}

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

      test.IntersectionTester(sphere, p1, p2, t)
    }
  }
}

func TestNewEllipsoid(t *testing.T) {
  if nil != NewEllipsoid(nil, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 1")
  }
  if nil != NewEllipsoid([]float64{0,0}, nil, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 2")
  }
  if nil != NewEllipsoid([]float64{0,0}, [][]float64{nil, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 3")
  }
  if nil != NewEllipsoid([]float64{0,0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 4")
  }
  if nil != NewEllipsoid([]float64{0,0}, [][]float64{[]float64{1,0,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 5")
  }
  if nil != NewEllipsoid([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1,1}) {
    t.Error("New ellipsoid surface error 6")
  }
  if nil != NewEllipsoid([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 7")
  }

  if nil == NewEllipsoid([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}}, []float64{1, 1}) {
    t.Error("New ellipsoid surface error 8")
  }
}

//No need to test other functions because we know from testing the
//general polynomial surfaces. 
func TestEllipsoidF(t *testing.T) {

  basis := [][]float64{[]float64{1, 0, 0}, []float64{0, 1, 0}, []float64{0, 0, 1}}

  //First some tests centered at the origin.
  for i := 0; i < 5; i ++ {
    point := []float64{0,0,0}

    param := []float64{math.Exp(test.RandFloat(-1, 1)),
      math.Exp(test.RandFloat(-1, 1)), math.Exp(test.RandFloat(-1, 1))}

    ellipsoid := NewEllipsoid(point, basis, param)

    for j := 0; j < 5; j ++ {
      test_point := []float64{test.RandFloat(-5, 5), test.RandFloat(-5, 5), test.RandFloat(-5, 5)}

      var f_exp float64 = 1
      f_test := ellipsoid.F(test_point)

      for k := 0; k < 3; k ++ {
        f_exp -= (test_point[k]-point[k])*(test_point[k]-point[k])/param[k]
      }

      if !test.CloseEnough(f_test, f_exp, err_bs) {
        t.Error("ellipsoid error: ", ellipsoid.String(), "; for parameters ", param, " at point ", test_point, ": expected ", f_exp, ", got ", f_test)
      }
    }
  }

  for i := 0; i < 5; i ++ {
    point := []float64{test.RandFloat(-5, 5), test.RandFloat(-5, 5), test.RandFloat(-5, 5)}

    param := []float64{math.Exp(test.RandFloat(-1, 1)),
      math.Exp(test.RandFloat(-1, 1)), math.Exp(test.RandFloat(-1, 1))}

    ellipsoid := NewEllipsoid(point, basis, param)

    for j := 0; j < 5; j ++ {
      test_point := []float64{test.RandFloat(-5, 5), test.RandFloat(-5, 5), test.RandFloat(-5, 5)}

      var f_exp float64 = 1
      f_test := ellipsoid.F(test_point)

      for k := 0; k < 3; k ++ {
        f_exp -= (test_point[k]-point[k])*(test_point[k]-point[k])/param[k]
      }

      if !test.CloseEnough(f_test, f_exp, err_bs) {
        t.Error("ellipsoid error for ", ellipsoid.String(), "; for parameters ", param, ", central point ", point, ", at point ", test_point, ": expected ", f_exp, ", got ", f_test)
      }
    }
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
  dim := 4
  var v [][]float64
  var c []float64
  //Test with several different numbers of infinite dimensions.
  for i := 1; i < dim; i ++ {
    c = test.RandFloatVector(-1, 1, dim)

    v = make([][]float64, i)
    for j := 0; j < i; j ++ {
      v[j] = test.RandFloatVector(-4, 4, dim)
    }

    cylinder := NewInfiniteCylinder(c, v)

    //Five trials per test.
    for j := 0; j < 5; j ++ {
      p := test.RandFloatVector(-10, 10, dim) 
      test_point := make([]float64, dim)
      for k := 0; k < dim; k ++ {
        test_point[k] = p[k] - c[k]
      }

      test_f := cylinder.F(p)

      var exp_f float64 = 1
      for k := 0; k < i; k ++ {
        for l := 0; l < dim; l ++ {
          for m := 0; m < dim; m ++ {
            exp_f -= test_point[l] * v[k][l] * v[k][m] * test_point[m]
          }
        }
      }

      if !test.CloseEnough(exp_f, test_f, err_bs) {
        t.Error("error! ", cylinder.String(), "; param(", c, v, "), at point ",
          p, "; expected ", exp_f, " got ", test_f)
      }
    }
  }
}

func TestNewInfiniteCone(t *testing.T) {
  //TODO
}

func TestInfiniteConeF(t *testing.T) {
  //TODO
}

func TestNewInfiniteParaboloid(t *testing.T) {
  //TODO
}

func TestInfiniteParaboloidF(t *testing.T) {
  //TODO
}

func TestNewInfiniteHyperboloid(t *testing.T) {
  //TODO
}

func TestInfiniteHyperboloidF(t *testing.T) {
  //TODO
}

//Since the boolean objects and the components have been independently tested, 
//all that needs to be tested here is whether the object has the correct shape
//and whether the constructor fails and succeeds correctly. 
func TestNewCylinder(t *testing.T) {
  //TODO
}

func TestCylinderInterior(t *testing.T) {
  //TODO
}
