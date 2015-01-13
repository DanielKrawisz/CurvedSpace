package pathtrace

import "testing"
import "math"
import "../test"
import "../distributions"
import "../vector"

var mat_err float64 = .00001

//Variables and functions used for mockingn the random distributions
var sphereSurfacePoint [3]float64
var normallyDistributedVector []float64

func mockRandomSphereSurfacePoint() *[3]float64 {
  return &sphereSurfacePoint
}

func mockRandomNormallyDistributedVector(n int, mean, sigma float64) []float64 {
  return normallyDistributedVector
}

func getRandomIncoming(norm []float64) (incoming []float64) {
  for {
    incoming = distributions.RandomUnitSphereSurfacePoint()[:]

    if vector.Dot(norm, incoming) < 0 {
      break
    }
  }

  return
}

//This only tests dimension 3 right now. 
func TestLambertianReflection(t *testing.T) {
  //Change function used to select random vector so that it can be mocked out. 
  randomUnitSphereSurfacePoint = mockRandomSphereSurfacePoint

  dim := 3
  l := LambertianReflection

  normal := make([]float64, dim)
  normal[0] = .5
    
  //Should be an odd number to avoid the case
  var pts int = 17

  for i := 0; i < pts; i ++ {
    theta := 2 * math.Pi * float64(i) / float64(pts)

    sphereSurfacePoint[0] = math.Cos(theta)/2.
    sphereSurfacePoint[1] = math.Sin(theta)/2.
    sphereSurfacePoint[2] = 0

    expected := make([]float64, 3)
    for j := 0; j < 3; j ++ {
      expected[j] = normal[j] + sphereSurfacePoint[j]
    }

    vector.Normalize(expected)

    output := vector.Normalize(l([]float64{1,0,0}, normal))
    if !test.VectorCloseEnough(expected, output, .000001) {
      t.Error("Lambertian error: sphere surface point ",
        sphereSurfacePoint, ", expected ", expected, " output ", output)
    }
  }

  //Chage function back to how it was. 
  randomUnitSphereSurfacePoint = distributions.RandomUnitSphereSurfacePoint
}

func TestMirrorReflection(t *testing.T) {
  var norm []float64 = distributions.RandomUnitSphereSurfacePoint()[:]
  incoming := getRandomIncoming(norm)

  rf := MirrorReflection
  outgoing := rf.Interact(incoming, norm)
  d := vector.Dot(incoming, norm)
  test_vector := make([]float64, len(norm))
  for i := 0; i < len(norm); i++ {
    test_vector[i] = incoming[i] - 2 * d * norm[i]
  }

  if !test.VectorCloseEnough(test_vector, outgoing, mat_err) {
    t.Error("Mirror reflection error.",
      "incoming = ", incoming, "; outgoing = ", outgoing, "; norm = ", norm, " test_vector = ", test_vector)
  }
}

//TODO
//We test the specular reflection by comparing it to the mirror reflection.
/*func TestSpecularReflection(t *testing.T) {
  //Change function used to select random vector so that it can be mocked out. 
  randomNormallyDistributedVector = mockRandomNormallyDistributedVector

  var norm []float64 = distributions.RandomUnitSphereSurfacePoint()[:]

  normallyDistributedVector = make([]float64, len(norm))

  for i := 0; i < len(norm); i ++ {
    normallyDistributedVector[i] = norm[i]
  }

  incoming := getRandomIncoming(norm)

  rf := NewMirrorReflection() 
  sf := NewSpecularReflection(1)

  outgoing_exp  := rf.Interact(incoming, norm)
  outgoing_test := sf.Interact(incoming, norm)

  if test.VectorCloseEnough(vector.Minus(outgoing_test, outgoing_exp), normallyDistributedVector, mat_err) {
    t.Error("Mirror reflection error.")
  }

  for i := 0; i < len(norm); i ++ {
    normallyDistributedVector[i] = -2 * norm[i]
  }

  outgoing_test = sf.Interact(incoming, norm)

  d := vector.Dot(outgoing_test, norm)

  for i := 0; i < len(norm); i ++ {
    outgoing_test[i] -= norm[i]
  }

  //Chage function back to how it was. 
  randomNormallyDistributedVector = distributions.RandomNormallyDistributedVector
}

func TestBasicRefractiveTransmission(t *testing.T) {
  
}*/
