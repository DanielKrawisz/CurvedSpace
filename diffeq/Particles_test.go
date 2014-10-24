package diffeq

//Tests for the particle system. 

import "testing"
import "math"
import "fmt"
import "../test"

//NewNewtonianParticle should alwys return nil in this test.
func TestNewNewtonianParticleFailureConditions(t *testing.T) {
  fmt.Println()

  invalid_dimensions := []int{-1, 0}
  invalid_particles := []int{-1, 0}
  invalid_initSes := []float64{math.Inf(1), math.Inf(-1), math.NaN()}
  invalid_initDses := []float64{math.Inf(1), math.Inf(-1), math.NaN(), 0}
  invalid_initposes := [][][]float64{
    nil, 
    [][]float64{
      []float64{0, 0},
      []float64{1, math.Inf(-1)}}, 
    [][]float64{
      []float64{math.Inf(1), 0},
      []float64{1, 0}}, 
    [][]float64{
      []float64{0, math.NaN()},
      []float64{1, 0}}, 
    [][]float64{
      []float64{1, 0}}, 
    [][]float64{
      []float64{0, 0},
      []float64{1}}}
  invalid_initvels := [][][]float64{
    nil, 
    [][]float64{
      []float64{1, 0},
      []float64{0, math.Inf(-1)}}, 
    [][]float64{
      []float64{math.Inf(1), 0},
      []float64{0, 0}}, 
    [][]float64{
      []float64{1, math.NaN()},
      []float64{0, 0}}, 
    [][]float64{
      []float64{1, 0}}, 
    [][]float64{
      []float64{1, 0},
      []float64{0}}}

  valid_dimension := 2
  valid_particles := 2
  var valid_initS float64 = 0
  var valid_initDs float64 = 1
  valid_initpos := [][]float64{
    []float64{0, 0},
    []float64{1, 0}}
  valid_initvel := [][]float64{
    []float64{0, 1},
    []float64{0, 0}}

  //Test 1st argument.
  for _, invalid_dimension := range invalid_dimensions {
    if NewNewtonianParticle(invalid_dimension, valid_particles, 
      valid_initS, valid_initDs, valid_initpos, valid_initvel) != nil {
      t.Error("Invalid first argument ", invalid_dimensions, " failed to be rejected.")
    }
  }

  //Test 2nd argument.
  for _, invalid_particle := range invalid_particles {
    if NewNewtonianParticle(valid_dimension, invalid_particle, 
      valid_initS, valid_initDs, valid_initpos, valid_initvel) != nil {
      t.Error("Invalid second argument ", invalid_particle, " failed to be rejected.")
    }
  }

  //Test 3rd argument.
  for _, invalid_initS := range invalid_initSes {
    if NewNewtonianParticle(valid_dimension, valid_particles, 
      invalid_initS, valid_initDs, valid_initpos, valid_initvel) != nil {
      t.Error("Invalid fourth argument ", invalid_initS, " failed to be rejected.")
    }
  }

  //Test 4th argument.
  for _, invalid_initDs := range invalid_initDses {
    if NewNewtonianParticle(valid_dimension, valid_particles, 
      valid_initS, invalid_initDs, valid_initpos, valid_initvel) != nil {
      t.Error("Invalid fifth argument ", invalid_initDs, " failed to be rejected.")
    }
  }

  //Test 5th argument.
  for _, invalid_initpos := range invalid_initposes {
    if NewNewtonianParticle(valid_dimension, valid_particles, 
      valid_initS, valid_initDs, invalid_initpos, valid_initvel) != nil {
      t.Error("Invalid sixth argument ", invalid_initpos, " failed to be rejected.")
    }
  }

  //Test 6th argument.
  for _, invalid_initvel := range invalid_initvels {
    if NewNewtonianParticle(valid_dimension, valid_particles, 
      valid_initS, valid_initDs, valid_initpos, invalid_initvel) != nil {
      t.Error("Invalid seventh argument ", invalid_initvel, " failed to be rejected.")
    }
  }
}

//Set up a random particle system. 
func randomParticleSystem() (np *newtonianParticle, dimensions, particles int, position, velocity [][]float64) {
  dimensions = test.RandInt(1, 5)
  particles = test.RandInt(1, 5)
  position = make([][]float64, particles)
  velocity = make([][]float64, particles)

  for i := 0; i < particles; i++ {
    position[i] = make([]float64, dimensions)
    velocity[i] = make([]float64, dimensions)

    for j := 0; j < dimensions; j++ {
      position[i][j] = test.RandFloat(-1000, 1000)
      velocity[i][j] = test.RandFloat(-1, 1)
    }
  }
 
  np = NewNewtonianParticle(dimensions, particles, 0, 1, position, velocity)

  for i := 0; i < particles; i++ {
    for j := 0; j < dimensions; j++ {
      np.a[i][j] = test.RandFloat(-1, 1)
      np.newx[i][j] = test.RandFloat(-1, 1)
      np.newv[i][j] = test.RandFloat(-1, 1)
      np.newa[i][j] = test.RandFloat(-1, 1)
    }
  }

  return
}

//Test that NewNewtonianParticle returns a consistent object. 
func TestNewNewtonianParticleConsistency(t *testing.T) {
  for i := 0; i < 10; i ++ {
    np, dimensions, particles, position, velocity := randomParticleSystem()
    n := particles * dimensions

    //This should never return nil. 
    if np == nil {
      t.Error("newtonianParticle creation failed with inputs ",
        "dimension: ", dimensions, "; particles: ", particles, "; position: ", position, "; velocity: ", velocity)
    }

    //Now generate some random array elements to test. 
    for j := 0; j < 3; j++ {
      var check_particle = test.RandInt(0, particles - 1)
      var check_dimension = test.RandInt(0, dimensions - 1)

      for k := 0; k < 2; k++ {

        if np.x[check_particle][check_dimension] != np.pos[check_particle * dimensions + check_dimension] {
          t.Error("Inconsistent particle system at x[", check_particle, "][", check_dimension, "]") 
        } 

        if np.v[check_particle][check_dimension] != np.pos[n + check_particle * dimensions + check_dimension] {
          t.Error("Inconsistent particle system at v[", check_particle, "][", check_dimension, "]") 
        }

        if np.a[check_particle][check_dimension] != np.vel[n + check_particle * dimensions + check_dimension] {
          t.Error("Inconsistent particle system at a[", check_particle, "][", check_dimension, "]") 
        }

        if np.newx[check_particle][check_dimension] != np.newpos[check_particle * dimensions + check_dimension] {
          t.Error("Inconsistent particle system at newx[", check_particle, "][", check_dimension, "]") 
        } 

        if np.newv[check_particle][check_dimension] != np.newpos[n + check_particle * dimensions + check_dimension] {
          t.Error("Inconsistent particle system at newv[", check_particle, "][", check_dimension, "]") 
        }

        if np.newa[check_particle][check_dimension] != np.newvel[n + check_particle * dimensions + check_dimension] {
          t.Error("Inconsistent particle system at newa[", check_particle, "][", check_dimension, "]") 
        }

        np.postStep()
      }
    }
  }
}

//Test all the member functions of a newtonianParticle
func TestNewtonianParticleFunctions(t *testing.T) {
  for i := 0; i < 10; i ++ {
    np, dimensions, particles, position, velocity := randomParticleSystem()

    //These next few tests are just getters and setters. 
    n := 2 * dimensions * particles
    z := np.Length()
    if z != n {
      t.Error("Length() error: expected ", n, ", got ", z, ".")
    }

    s := np.S()
    if s != 0.0 {
      t.Error("S() error: expected ", 0.0, ", got ", s, ".")
    }

    ds := np.Ds()
    if ds != 1.0 {
      t.Error("Ds() error: expected ", 1.0, ", got ", ds, ".")
    }

    newds := test.RandFloat(1, 23)
    np.setDs(newds)
    ds = np.Ds()
    if ds != newds {
      t.Error("setDs() error: expected ", newds, ", got ", ds, ".")
    }

    //Test velocity() and position()
    var check_p = test.RandInt(0, particles - 1)
    var check_d = test.RandInt(0, dimensions - 1)
    x := np.position()[dimensions * check_p + check_d]
    if x != position[check_p][check_d] {
      t.Error("position() error: expected ", position[check_p][check_d], ", got ", x, ".")
    }

    v := np.position()[n/2 + dimensions * check_p + check_d]
    if v != velocity[check_p][check_d] {
      t.Error("velocity() error: expected ", velocity[check_p][check_d], ", got ", v, ".")
    }

    //This next part checks newposition(), newvelocity(), and postStep()
    newx := test.RandFloat(-10, 10)
    newv := test.RandFloat(-10, 10)
    np.newPosition()[dimensions * check_p + check_d] = newx
    np.newPosition()[n / 2 + dimensions * check_p + check_d] = newv
    np.postStep()
    newx_test := np.position()[dimensions * check_p + check_d]
    newv_test := np.position()[n / 2 + dimensions * check_p + check_d]
    x_test := np.newPosition()[dimensions * check_p + check_d]
    v_test := np.newPosition()[n / 2 + dimensions * check_p + check_d]
 
    if newx != newx_test {
      t.Error("newPosition() error: expected ", newx, ", got ", newx_test, ".")
    }

    if newv != newv_test {
      t.Error("newVelocity() error: expected ", newv, ", got ", newv_test, ".")
    }

    if x != x_test {
      t.Error("postStep() error (1): expected ", x, ", got ", x_test, ".")
    }

    if v != v_test {
      t.Error("postStep() error (2): expected ", v, ", got ", v_test, ".")
    }
  }
}

//Test failure conditions. 
func TestNewHarmonicOscillatorFailureConditions(t *testing.T) {
  valid_particle := NewNewtonianParticle(2, 3, 0, 1,
    [][]float64{
      []float64{0, 0}, 
      []float64{1, 0},
      []float64{1, 1}},
    [][]float64{
      []float64{0, -1}, 
      []float64{0, 0},
      []float64{1, 1}})
  valid_mass := []float64{1, 1.5, 2}
  valid_K := [][][]float64 { //A valid K
    [][]float64{nil, nil, nil},
    [][]float64{[]float64{1}, nil, nil},
    [][]float64{[]float64{0}, []float64{1, 2}, nil}}

  invalid_particles := []*newtonianParticle{
    nil, 
    NewNewtonianParticle(2, 2, 0, 1, //Wrong number of particles.
    [][]float64{
      []float64{0, 0}, 
      []float64{1, 0}},
    [][]float64{
      []float64{0, -1}, 
      []float64{0, 0}})}
  invalid_masses := [][]float64{
    nil, 
    []float64{0, 1, 1}, 
    []float64{math.Inf(1), 1, 1}, 
    []float64{-1, 1, 1}, 
    []float64{math.NaN(), 1, 1}, 
    []float64{1}}
  invalid_Ks := [][][][]float64 {
    nil, //A nil K
    [][][]float64 { //K in which an entry below the diagonal is nil. 
      [][]float64{nil, nil, nil},
      [][]float64{[]float64{1}, nil, nil},
      [][]float64{nil, []float64{1, 2}, nil}}, 
    [][][]float64 { //not enough rows. 
      [][]float64{nil, nil, nil},
      [][]float64{[]float64{1}, nil, nil}}, 
    [][][]float64 { //not enough columns
      [][]float64{nil,},
      [][]float64{},
      [][]float64{[]float64{0}, []float64{1, 2}, nil}}}

  for _, invalid_particle := range invalid_particles {
    if NewHarmonicOscillator(invalid_particle, valid_mass, valid_K) != nil {
      t.Error("Invalid particles ", invalid_particle, " failed to be rejected.")
    }
  }

  for _, invalid_mass := range invalid_masses {
    if NewHarmonicOscillator(valid_particle, invalid_mass, valid_K) != nil {
      t.Error("Invalid mass ", invalid_mass, " failed to be rejected.")
    }
  }

  for _, invalid_K := range invalid_Ks {
    if NewHarmonicOscillator(valid_particle, valid_mass, invalid_K) != nil {
      t.Error("Invalid K ", invalid_K, " failed to be rejected.")
    }
  }
}

//Error propagation formula for the distance vector formula. 
func distanceVectorErrorPropagation(a, b, sa, sb float64) float64 {
  da := sa * a
  db := sa * b
  return math.Sqrt(da*da + db*db)
}

//Error propagation formula for the quadrance. 
func quadranceErrorPropagation(dist []float64, sd float64) float64 {
  var err float64 = 0
  for i := 0; i < len(dist); i++ {
    d := dist[i] * sd
    err += d * d
  }
  return 2 * math.Sqrt(err)
}

//Error propagation formula for the oscillator force. 
func oscillatorForceErrorPropagation(quad float64, K []float64, dx, dk float64) float64 {
  if len(K) == 0 {return .0000000000001} //The error should not be able to return zero.

  var err float64 = 0
  var errx float64 = 0

  var xpow float64 = 1
  var odd float64 = 1
  var qpow float64 = quad
  for i := 0; i < len(K); i++ {
    errx += K[i] * odd * xpow
    odd += 2
    xpow *= quad

    err += qpow 
    qpow *= (quad * quad)
  }
  return math.Sqrt(errx * errx * dx * dx + err * dk * dk)
}

func TestHarmonicOscillatorForceCalculation(t *testing.T) {
  var tolerance float64 = .0000001 //The accepted error for the inputs of most float calculations. 

  //First we do a few set examples, and then some random ones. 
  //Here are the set examples. 
  dimension := 2
  p1 := []float64{-1, 0}
  p2 := []float64{1, 0}

  //The set of force constants to be used in each calculation.
  Klist := [][]float64{[]float64{}, []float64{3}, []float64{2, 1}, []float64{3, 2, 1}}

  //The expected results. (The force is only along the x-axis, so only those numbers are included.)
  expected_forces := []float64{0, 6, 12, 54}

  for i, K := range Klist {
    distv := make([]float64, dimension)
    distanceVector(p1, p2, distv)
    force := oscillatorForce(distv, K)

    if force[0] != expected_forces[i] {
      t.Error("Oscillator force calculation error on trial ", i, ": expected ", expected_forces[i], ", got ", force[0])
    }
  }

  //Here are the random calculations. 
  for j := 0; j < 10; j ++ {
    dimension = test.RandInt(1, 5)
    p1 = make([]float64, dimension)
    p2 = make([]float64, dimension)

    for i := 0; i < dimension; i++ {
      p1[i] = test.RandFloat(-10, 10)
      p2[i] = test.RandFloat(10, 10)
    }

    Klen := test.RandInt(0, 5) //The number of terms in the force function. 
    K := make([]float64, Klen)

    var d float64 = 1
    for i := 0; i < Klen; i++ {
      K[i] = test.RandFloat(0, 3) / d //Each force constant gets smaller and smaller. 
      d *= 2
    }

    distv := make([]float64, dimension)
    distanceVector(p1, p2, distv)
    quad := quadrance(distv)

    //check if distv was calculated correctly. 
    vector_check := test.RandInt(0, dimension - 1)
    var expected float64 = p2[vector_check] - p1[vector_check]
    if !test.CloseEnough(distv[vector_check], expected,
      distanceVectorErrorPropagation(p2[vector_check], p1[vector_check], tolerance, tolerance)) {
      t.Error("Distance vector error: distancev was ", distv[vector_check], "; expected ", expected, ".")
    }

    //What about quad? 
    var expected_quad float64 = 0
    for i:= 0; i < dimension; i++ {
      expected_quad += distv[i] * distv[i]
    }
    if !test.CloseEnough(quad, expected_quad, quadranceErrorPropagation(distv, tolerance)) {
      t.Error("Quadrance error: quad was ", quad, "; expected ", expected_quad)
    }

    //Finally, check the force. 
    force := oscillatorForce(distv, K) 

    var expected_force []float64 = make([]float64, len(force))
    var xpow float64 = 1
    var fc float64 = 0

    for k := 0; k < len(K); k++ {
      fc += xpow * K[k]
      xpow *= quad
    }

    for i := 0; i < len(distv); i++ {
      expected_force[i] = distv[i] * fc
    }

    if !test.CloseEnough(force[vector_check], expected_force[vector_check],
      oscillatorForceErrorPropagation(quad, K, tolerance, tolerance)) {
      t.Error("oscillatorForce error: force was ", force[vector_check], "; expected ", expected_force[vector_check])
    }
  }
}

//From the previous test, we know that the force is calculated correctly, so we just have a few worked examples
//to make sure the rest of the function works. 
func TestHarmonicOscillatorAcceleration(t *testing.T) {
  velocity := [][]float64{[]float64{0, 1}, []float64{0, -1}}
  p := NewNewtonianParticle(2, 2, 0, 1,
    [][]float64{[]float64{-1, 0}, []float64{1, 0}},
    velocity)

  //The set of force constants to be used in each calculation.
  Klist := [][]float64{[]float64{}, []float64{3}, []float64{2, 1}, []float64{3, 2, 1}}
  //The masses to be used in each calculation.
  mlist := [][]float64{[]float64{1, 1}, []float64{1, 2}, []float64{1, 3}, []float64{1, 6}}

  //The expected accelerations along the x axis for each particle.
  expected_accel := [][]float64{[]float64{0, 0}, []float64{6, -3}, []float64{12, -4}, []float64{54, -9}}

  for i, K := range Klist {
    Kmtrx := [][][]float64{nil, [][]float64{K}}
    ha := NewHarmonicOscillator(p, mlist[i], Kmtrx)
    ha.DxDs(p.position(), p.velocity())

    if p.vel[1] != velocity[0][1] {
      t.Error("Oscillator velocity calculation error on trial ", i, ": expected ", velocity[0][1], ", got ", p.vel[1])
    }

    if p.vel[3] != velocity[1][1] {
      t.Error("Oscillator velocity calculation error on trial ", i, ": expected ", velocity[1][1], ", got ", p.vel[3])
    }

    if p.a[0][0] != expected_accel[i][0] {
      t.Error("Oscillator acceleration calculation error on trial ", i, ": expected ", expected_accel[i][0], ", got ", p.a[0][0])
    }

    if p.a[1][0] != expected_accel[i][1] {
      t.Error("Oscillator acceleration calculation error on trial ", i, ": expected ", expected_accel[i][1], ", got ", p.a[1][0])
    }
  }
}
