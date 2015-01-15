package pathtrace

import "testing"

//TODO
//Just testing some particular conditions here. 

func testCameraInverse() {
  
}

func TestIsometricCamera(t *testing.T) {
  pos  := []float64{.1, .2, .3}
  mtrx := [][]float64{[]float64{1.3, 0, 0}, []float64{0, 1.2, 0}, []float64{0, 0, 1.1}}

  if IsometricCamera(nil, mtrx, 3, 3, 1.6, 1.7) != nil { t.Error("isometric camera error 1") }
  if IsometricCamera(pos, nil, 3, 3, 1.6, 1.7) != nil { t.Error("isometric camera error 2") }

  cam := IsometricCamera(pos, mtrx, 3, 3, 1.6, 1.7)

  if cam == nil { t.Error("isometric camera error 3") }
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
}

func TestToroidialCamera(t *testing.T) {
  
}

/*func TestTrapezoidalCamera(t *testing.T) {
  
}*/
