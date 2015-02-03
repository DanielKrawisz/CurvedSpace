package pathtrace

import "testing"
import "github.com/DanielKrawisz/CurvedSpace/test"
import "github.com/DanielKrawisz/CurvedSpace/vector"

//Just testing some particular conditions here. 

var cam_err float64 = .000001

type camTestCase struct {
  i, j int
  pos, dir []float64
}

//A mocking function to make testing easier. 
func MockCameraStochastic() float64 {
  return 0
}

func cameraTest(kind string, test_cases []*camTestCase, cam, cam_inv GenerateRay, t *testing.T) {
  camJitter = MockCameraStochastic
  zero := []float64{0, 0, 0}

  for i, test_case := range test_cases {
    pos, dir := cam(test_case.i, test_case.j) 
    pos_inv, dir_inv := cam_inv(test_case.i, test_case.j)

    if !test.VectorCloseEnough(test_case.pos, pos, cam_err) ||
      !test.VectorCloseEnough(test_case.dir, dir, cam_err) {
      t.Error(kind, " error! case ", i,":", test_case.i, test_case.j,
        "\n\tgot ", pos, dir, "\n\texpected ", test_case.pos, test_case.dir)
    }

    if !test.VectorCloseEnough(test_case.pos, vector.Plus(pos_inv, dir_inv), cam_err) ||
      !test.VectorCloseEnough(test_case.dir, vector.Minus(zero, dir_inv), cam_err) {
      t.Error("inverse ", kind, " cam error! case ", i,":", test_case.i, test_case.j,
        "\n\tgot ", pos_inv, dir_inv, "\n\texpected ", test_case.pos, test_case.dir)
    }
  }

  camJitter = CameraStochastic
}

func TestCameraCoordinates(t *testing.T) {
  var ou, ov float64
  camJitter = MockCameraStochastic

  ou, ov = CameraCoordinates(1, 1, 3, 3, 1, 1)
  if !(test.CloseEnough(ou, 0, cam_err) && test.CloseEnough(ov, 0, cam_err)) {
    t.Error("camera coordinates error, case 1, got ", ou, ov)
  }
  ou, ov = CameraCoordinates(0, 1, 3, 3, 1, 1)
  if !(test.CloseEnough(ou, -1, cam_err) && test.CloseEnough(ov, 0, cam_err)) {
    t.Error("camera coordinates error, case 2, got ", ou, ov)
  }
  ou, ov = CameraCoordinates(1, 0, 3, 3, 1, 1)
  if !(test.CloseEnough(ou, 0, cam_err) && test.CloseEnough(ov, 1, cam_err)) {
    t.Error("camera coordinates error, case 3, got ", ou, ov)
  }
  ou, ov = CameraCoordinates(2, 2, 3, 3, 1.3, 1.7)
  if !(test.CloseEnough(ou, 1.3, cam_err) && test.CloseEnough(ov, -1.7, cam_err)) {
    t.Error("camera coordinates error, case 3, got ", ou, ov)
  }

  camJitter = CameraStochastic
}

func TestIsometricCamera(t *testing.T) {
  camJitter = MockCameraStochastic

  pos  := []float64{.1, .2, .3}
  mtrx := [][]float64{[]float64{1.3, 0, 0}, []float64{0, 1.2, 0}, []float64{0, 0, 1.1}}

  if IsometricCamera(nil, mtrx, 3, 3, 1.6, 1.7) != nil { t.Error("isometric camera error 1") }
  if IsometricCamera(pos, nil, 3, 3, 1.6, 1.7) != nil { t.Error("isometric camera error 2") }

  cam := IsometricCamera(pos, mtrx, 3, 3, 1.6, 1.7)

  if cam == nil { t.Error("isometric camera error 3") }

  test_cases := []*camTestCase{
    &camTestCase{0, 0, []float64{0.1, 2.24, -1.46}, []float64{1.3, 2.04, -1.76}}, 
    &camTestCase{0, 1, []float64{0.1, 0.2, -1.46}, []float64{1.3, 0., -1.76}}, 
    &camTestCase{0, 2, []float64{0.1, -1.84, -1.46}, []float64{1.3, -2.04, -1.76}}, 
    &camTestCase{1, 0, []float64{0.1, 2.24, 0.3}, []float64{1.3, 2.04, 0.}}, 
    &camTestCase{1, 1, []float64{0.1, 0.2, 0.3}, []float64{1.3, 0., 0.}}, 
    &camTestCase{1, 2, []float64{0.1, -1.84, 0.3}, []float64{1.3, -2.04, 0.}}, 
    &camTestCase{2, 0, []float64{0.1, 2.24, 2.06}, []float64{1.3, 2.04, 1.76}}, 
    &camTestCase{2, 1, []float64{0.1, 0.2, 2.06}, []float64{1.3, 0., 1.76}}, 
    &camTestCase{2, 2, []float64{0.1, -1.84, 2.06}, []float64{1.3, -2.04, 1.76}}}

  for i, test_case := range test_cases {
    pos, dir := cam(test_case.i, test_case.j) 

    if !test.VectorCloseEnough(test_case.pos, pos, cam_err) ||
      !test.VectorCloseEnough(test_case.dir, dir, cam_err) {
      t.Error("cam error! case ", i,":", test_case.i, test_case.j,
        "\n\tgot ", pos, dir, "\n\texpected ", test_case.pos, test_case.dir)
    }
  }

  camJitter = CameraStochastic
}

func TestFlatCamera(t *testing.T) {
  pos  := []float64{.1, .2, .3}
  mtrx := [][]float64{[]float64{1.3, 0, 0}, []float64{0, 1.2, 0}, []float64{0, 0, 1.1}}

  if FlatCamera(nil, mtrx, 3, 3, 1.6, 1.7) != nil { t.Error("Flat camera error 1") }
  if FlatCamera(pos, nil, 3, 3, 1.6, 1.7) != nil { t.Error("Flat camera error 2") }

  cam := FlatCamera(pos, mtrx, 3, 3, 1.6, 1.7)
  cam_inv := InverseFlatCamera(pos, mtrx, 3, 3, 1.6, 1.7)

  if cam == nil { t.Error("Flat camera error 3") }
  if cam_inv == nil { t.Error("Flat camera error 4") }

  test_cases := []*camTestCase{
    &camTestCase{0, 0, []float64{0.1, 0.2, 0.3}, []float64{1.3, 2.04, -1.76}}, 
    &camTestCase{0, 1, []float64{0.1, 0.2, 0.3}, []float64{1.3, 0., -1.76}}, 
    &camTestCase{0, 2, []float64{0.1, 0.2, 0.3}, []float64{1.3, -2.04, -1.76}}, 
    &camTestCase{1, 0, []float64{0.1, 0.2, 0.3}, []float64{1.3, 2.04, 0.}}, 
    &camTestCase{1, 1, []float64{0.1, 0.2, 0.3}, []float64{1.3, 0., 0.}}, 
    &camTestCase{1, 2, []float64{0.1, 0.2, 0.3}, []float64{1.3, -2.04, 0.}}, 
    &camTestCase{2, 0, []float64{0.1, 0.2, 0.3}, []float64{1.3, 2.04, 1.76}}, 
    &camTestCase{2, 1, []float64{0.1, 0.2, 0.3}, []float64{1.3, 0., 1.76}}, 
    &camTestCase{2, 2, []float64{0.1, 0.2, 0.3}, []float64{1.3, -2.04, 1.76}}}

  cameraTest("flat camera", test_cases, cam, cam_inv, t)
}

func TestCylindricalCamera(t *testing.T) {

  pos  := []float64{.1, .2, .3}
  mtrx := [][]float64{[]float64{1.3, 0, 0}, []float64{0, 1.2, 0}, []float64{0, 0, 1.1}}

  if CylindricalCamera(nil, mtrx, 3, 3, 1.6, 1.7) != nil { t.Error("Cylindrical camera error 1") }
  if CylindricalCamera(pos, nil, 3, 3, 1.6, 1.7) != nil { t.Error("Cylindrical camera error 2") }

  cam := CylindricalCamera(pos, mtrx, 3, 3, 1.6, 1.7)
  cam_inv := InverseCylindricalCamera(pos, mtrx, 3, 3, 1.6, 1.7)

  if cam == nil { t.Error("Cylindrical camera error 3") }
  if cam_inv == nil { t.Error("Cylindrical camera error 4") }

  test_cases := []*camTestCase{
    &camTestCase{0, 0, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, 2.04, -1.09953}}, 
    &camTestCase{0, 1, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, 0., -1.09953}}, 
    &camTestCase{0, 2, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, -2.04, -1.09953}}, 
    &camTestCase{1, 0, []float64{0.1, 0.2, 0.3}, []float64{1.3, 2.04, 0.}}, 
    &camTestCase{1, 1, []float64{0.1, 0.2, 0.3}, []float64{1.3, 0., 0.}}, 
    &camTestCase{1, 2, []float64{0.1, 0.2, 0.3}, []float64{1.3, -2.04, 0.}}, 
    &camTestCase{2, 0, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, 2.04, 1.09953}}, 
    &camTestCase{2, 1, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, 0., 1.09953}}, 
    &camTestCase{2, 2, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, -2.04, 1.09953}}}

  cameraTest("cylindrical camera", test_cases, cam, cam_inv, t)
}

func TestIsometricCylindricalCamera(t *testing.T) {

  pos  := []float64{.1, .2, .3}
  mtrx := [][]float64{[]float64{1.3, 0, 0}, []float64{0, 1.2, 0}, []float64{0, 0, 1.1}}

  if IsometricCylindricalCamera(nil, mtrx, 3, 3, 1.6, 1.7) != nil {
    t.Error("Isometric Cylindrical camera error 1") }
  if IsometricCylindricalCamera(pos, nil, 3, 3, 1.6, 1.7) != nil {
    t.Error("Isometric Cylindrical camera error 2") }

  cam := IsometricCylindricalCamera(pos, mtrx, 3, 3, 1.6, 1.7)
  cam_inv := InverseIsometricCylindricalCamera(pos, mtrx, 3, 3, 1.6, 1.7)

  if cam == nil { t.Error("Isometric Cylindrical camera error 3") }
  if cam_inv == nil { t.Error("Isometric Cylindrical camera error 4") }

  test_cases := []*camTestCase{
    &camTestCase{0, 0, []float64{0.1, 2.24, 0.3}, []float64{-0.0379594, 0, -1.09953}}, 
    &camTestCase{0, 1, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, 0, -1.09953}}, 
    &camTestCase{0, 2, []float64{0.1, -1.84, 0.3}, []float64{-0.0379594, 0, -1.09953}}, 
    &camTestCase{1, 0, []float64{0.1, 2.24, 0.3}, []float64{1.3, 0, 0.}}, 
    &camTestCase{1, 1, []float64{0.1, 0.2, 0.3}, []float64{1.3, 0, 0.}}, 
    &camTestCase{1, 2, []float64{0.1, -1.84, 0.3}, []float64{1.3, 0, 0.}}, 
    &camTestCase{2, 0, []float64{0.1, 2.24, 0.3}, []float64{-0.0379594, 0, 1.09953}}, 
    &camTestCase{2, 1, []float64{0.1, 0.2, 0.3}, []float64{-0.0379594, 0, 1.09953}}, 
    &camTestCase{2, 2, []float64{0.1, -1.84, 0.3}, []float64{-0.0379594, 0, 1.09953}}}

  cameraTest("isometric cylindrical camera", test_cases, cam, cam_inv, t)
}

func TestPolarSphericalCamera(t *testing.T) {

  pos  := []float64{.1, .2, .3}
  mtrx := [][]float64{[]float64{1.3, 0, 0}, []float64{0, 1.2, 0}, []float64{0, 0, 1.1}}

  if PolarSphericalCamera(nil, mtrx, 3, 3, 1.6, 1.7) != nil { t.Error("Polar Spherical camera error 1") }
  if PolarSphericalCamera(pos, nil, 3, 3, 1.6, 1.7) != nil { t.Error("Polar Spherical camera error 2") }

  cam := PolarSphericalCamera(pos, mtrx, 3, 3, 1.6, 1.7)
  cam_inv := InversePolarSphericalCamera(pos, mtrx, 3, 3, 1.6, 1.7)

  if cam == nil { t.Error("Polar Spherical camera error 3") }
  if cam_inv == nil { t.Error("Polar Spherical camera error 4") }

  test_cases := []*camTestCase{
    &camTestCase{0, 0, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000},
      []float64{0.00489085698995457, 1.18999777254296, 0.141668510934542}}, 
    &camTestCase{0, 1, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-0.0379593789916753, 0, -1.09953096334566}}, 
    &camTestCase{0, 2, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{0.00489085698995457, -1.18999777254296, 0.141668510934542}}, 
    &camTestCase{1, 0, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-0.167497842584182, 1.18999777254296, 0}}, 
    &camTestCase{1, 1, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{1.30000000000000, 0, 0}}, 
    &camTestCase{1, 2, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-0.167497842584182, -1.18999777254296, 0}}, 
    &camTestCase{2, 0, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{0.00489085698995457, 1.18999777254296, -0.141668510934542}}, 
    &camTestCase{2, 1, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-0.0379593789916753, 0, 1.09953096334566}}, 
    &camTestCase{2, 2, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{0.00489085698995457, -1.18999777254296, -0.141668510934542}}}

  cameraTest("polar spherical camera", test_cases, cam, cam_inv, t)
}

func TestSphericalCamera(t *testing.T) {

  pos  := []float64{.1, .2, .3}
  mtrx := [][]float64{[]float64{1.3, 0, 0}, []float64{0, 1.2, 0}, []float64{0, 0, 1.1}}

  if SphericalCamera(nil, mtrx, 3, 3, 1.6, 1.7) != nil { t.Error("Spherical camera error 1") }
  if SphericalCamera(pos, nil, 3, 3, 1.6, 1.7) != nil { t.Error("Spherical camera error 2") }

  cam := SphericalCamera(pos, mtrx, 3, 3, 1.6, 1.7)
  cam_inv := InverseSphericalCamera(pos, mtrx, 3, 3, 1.6, 1.7)

  if cam == nil { t.Error("Spherical camera error 3") }
  if cam_inv == nil { t.Error("Spherical camera error 4") }

  test_cases := []*camTestCase{
    &camTestCase{0, 0, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{0.00489290866924453, -1.18998934375352, 0.141727939854774}}, 
    &camTestCase{0, 1, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-0.0379593789916753, 0, -1.09953096334566}}, 
    &camTestCase{0, 2, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{0.00489290866924453, 1.18998934375352, 0.141727939854774}}, 
    &camTestCase{1, 0, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-1.30000000000000, 0, 0}}, 
    &camTestCase{1, 1, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{1.30000000000000, 0, 0}}, 
    &camTestCase{1, 2, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-1.30000000000000, 0, 0}}, 
    &camTestCase{2, 0, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{0.00489290866924453, 1.18998934375352, -0.141727939854774}}, 
    &camTestCase{2, 1, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{-0.0379593789916753, 0, 1.09953096334566}}, 
    &camTestCase{2, 2, []float64{0.100000000000000, 0.200000000000000, 0.300000000000000}, 
      []float64{0.00489290866924453, -1.18998934375352, -0.141727939854774}}}

  cameraTest("spherical camera", test_cases, cam, cam_inv, t)
}

func TestToroidialCamera(t *testing.T) {

  pos := []float64{.1, .2, .3}
  R   := [][]float64{[]float64{3.3, 0, 0}, []float64{0, 3.7, 0}}
  r   := [][]float64{[]float64{1.1, 0, 0}, []float64{0, 0, 1.4}}

  if ToroidialCamera(nil, R, r, 6, 6) != nil { t.Error("Toroidial camera error 1") }
  if ToroidialCamera(pos, nil, r, 6, 6) != nil { t.Error("Toroidial camera error 2") }
  if ToroidialCamera(pos, R, nil, 6, 6) != nil { t.Error("Toroidial camera error 3") }

  cam := ToroidialCamera(pos, R, r, 6, 6)
  cam_inv := InverseToroidialCamera(pos, R, r, 6, 6)

  if cam == nil { t.Error("Spherical camera error 4") }
  if cam_inv == nil { t.Error("Spherical camera error 5") }

  test_cases := []*camTestCase{
    &camTestCase{0, 0, []float64{-3.2000000, 0.20000000, 0.30000000}, []float64{0.98108108, 0, 0}}, 
    &camTestCase{0, 1, []float64{-3.2000000, 0.20000000, 0.30000000}, []float64{0.49054054, 0, 1.2124356}}, 
    &camTestCase{0, 2, []float64{-3.2000000, 0.20000000, 0.30000000}, []float64{-0.49054054, 0, 1.2124356}}, 
    &camTestCase{0, 3, []float64{-3.2000000, 0.20000000, 0.30000000}, []float64{-0.98108108, 0, 0}}, 
    &camTestCase{0, 4, []float64{-3.2000000, 0.20000000, 0.30000000}, []float64{-0.49054054, 0, -1.2124356}}, 
    &camTestCase{0, 5, []float64{-3.2000000, 0.20000000, 0.30000000}, []float64{0.49054054, 0, -1.2124356}}, 
    &camTestCase{1, 0, []float64{-1.5500000, -3.0042940, 0.30000000}, []float64{0.49054054, 0.95262794, 0}}, 
    &camTestCase{1, 1, []float64{-1.5500000, -3.0042940, 0.30000000},
      []float64{0.24527027, 0.47631397, 1.2124356}}, 
    &camTestCase{1, 2, []float64{-1.5500000, -3.0042940, 0.30000000}, 
      []float64{-0.24527027, -0.47631397, 1.2124356}}, 
    &camTestCase{1, 3, []float64{-1.5500000, -3.0042940, 0.30000000}, []float64{-0.49054054, -0.95262794, 0}}, 
    &camTestCase{1, 4, []float64{-1.5500000, -3.0042940, 0.30000000}, 
      []float64{-0.24527027, -0.47631397, -1.2124356}}, 
    &camTestCase{1, 5, []float64{-1.5500000, -3.0042940, 0.30000000}, 
      []float64{0.24527027, 0.47631397, -1.2124356}}, 
    &camTestCase{2, 0, []float64{1.7500000, -3.0042940, 0.30000000}, 
      []float64{-0.49054054, 0.95262794, 0}}, 
    &camTestCase{2, 1, []float64{1.7500000, -3.0042940, 0.30000000}, 
      []float64{-0.24527027, 0.47631397, 1.2124356}}, 
    &camTestCase{2, 2, []float64{1.7500000, -3.0042940, 0.30000000}, 
      []float64{0.24527027, -0.47631397, 1.2124356}}, 
    &camTestCase{2, 3, []float64{1.7500000, -3.0042940, 0.30000000}, []float64{0.49054054, -0.95262794, 0}}, 
    &camTestCase{2, 4, []float64{1.7500000, -3.0042940, 0.30000000}, 
      []float64{0.24527027, -0.47631397, -1.2124356}}, 
    &camTestCase{2, 5, []float64{1.7500000, -3.0042940, 0.30000000}, 
      []float64{-0.24527027, 0.47631397, -1.2124356}}, 
    &camTestCase{3, 0, []float64{3.4000000, 0.20000000, 0.30000000}, []float64{-0.98108108, 0, 0}}, 
    &camTestCase{3, 1, []float64{3.4000000, 0.20000000, 0.30000000}, []float64{-0.49054054, 0, 1.2124356}}, 
    &camTestCase{3, 2, []float64{3.4000000, 0.20000000, 0.30000000}, []float64{0.49054054, 0, 1.2124356}}, 
    &camTestCase{3, 3, []float64{3.4000000, 0.20000000, 0.30000000}, []float64{0.98108108, 0, 0}}, 
    &camTestCase{3, 4, []float64{3.4000000, 0.20000000, 0.30000000}, []float64{0.49054054, 0, -1.2124356}}, 
    &camTestCase{3, 5, []float64{3.4000000, 0.20000000, 0.30000000}, []float64{-0.49054054, 0, -1.2124356}}, 
    &camTestCase{4, 0, []float64{1.7500000, 3.4042940, 0.30000000}, []float64{-0.49054054, -0.95262794, 0}}, 
    &camTestCase{4, 1, []float64{1.7500000, 3.4042940, 0.30000000}, 
      []float64{-0.24527027, -0.47631397, 1.2124356}}, 
    &camTestCase{4, 2, []float64{1.7500000, 3.4042940, 0.30000000}, 
      []float64{0.24527027, 0.47631397, 1.2124356}}, 
    &camTestCase{4, 3, []float64{1.7500000, 3.4042940, 0.30000000}, []float64{0.49054054, 0.95262794, 0}}, 
    &camTestCase{4, 4, []float64{1.7500000, 3.4042940, 0.30000000}, 
      []float64{0.24527027, 0.47631397, -1.2124356}}, 
    &camTestCase{4, 5, []float64{1.7500000, 3.4042940, 0.30000000}, 
      []float64{-0.24527027, -0.47631397, -1.2124356}}, 
    &camTestCase{5, 0, []float64{-1.5500000, 3.4042940, 0.30000000}, []float64{0.49054054, -0.95262794, 0}}, 
    &camTestCase{5, 1, []float64{-1.5500000, 3.4042940, 0.30000000}, 
      []float64{0.24527027, -0.47631397, 1.2124356}}, 
    &camTestCase{5, 2, []float64{-1.5500000, 3.4042940, 0.30000000}, 
      []float64{-0.24527027, 0.47631397, 1.2124356}}, 
    &camTestCase{5, 3, []float64{-1.5500000, 3.4042940, 0.30000000}, []float64{-0.49054054, 0.95262794, 0}}, 
    &camTestCase{5, 4, []float64{-1.5500000, 3.4042940, 0.30000000}, 
      []float64{-0.24527027, 0.47631397, -1.2124356}}, 
    &camTestCase{5, 5, []float64{-1.5500000, 3.4042940, 0.30000000}, 
      []float64{0.24527027, -0.47631397, -1.2124356}}}

  cameraTest("toroidial camera", test_cases, cam, cam_inv, t)
}

/*func TestTrapezoidalCamera(t *testing.T) {
  
}*/
