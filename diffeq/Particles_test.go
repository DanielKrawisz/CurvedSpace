package diffeq

import "testing"
import "math"
import "math/rand"
import "fmt"
import "time"

var seed_set bool = false

func setSeed() {
  rand.Seed( time.Now().UTC().UnixNano())
  seed_set = true
}

func randInt(min int, max int) int {
  if !seed_set {setSeed()}
  return min + rand.Intn(max-min+1)
}

func randFloat(min, max float64) float64 {
  if !seed_set {setSeed()}
  return min + (max - min) * rand.Float64()
}

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
  dimensions = randInt(1, 5)
  particles = randInt(1, 5)
  position = make([][]float64, particles)
  velocity = make([][]float64, particles)

  for i := 0; i < particles; i++ {
    position[i] = make([]float64, dimensions)
    velocity[i] = make([]float64, dimensions)

    for j := 0; j < dimensions; j++ {
      position[i][j] = randFloat(-1000, 1000)
      velocity[i][j] = randFloat(-1, 1)
    }
  }
 
  np = NewNewtonianParticle(dimensions, particles, 0, 1, position, velocity)

  for i := 0; i < particles; i++ {
    for j := 0; j < dimensions; j++ {
      np.a[i][j] = randFloat(-1, 1)
      np.newx[i][j] = randFloat(-1, 1)
      np.newv[i][j] = randFloat(-1, 1)
      np.newa[i][j] = randFloat(-1, 1)
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
      var check_particle = randInt(0, particles - 1)
      var check_dimension = randInt(0, dimensions - 1)

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
    np, dimensions, particles, _, _ := randomParticleSystem()

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

    newds := randFloat(1, 23)
    np.setDs(newds)
    ds = np.Ds()
    if ds != newds {
      t.Error("setDs() error: expected ", newds, ", got ", ds, ".")
    }

    /* These tests are screwed up... need to get them right! 
    var check_p = randInt(0, particles - 1)
    var check_d = randInt(0, dimensions - 1)
    x := np.position()[particles * check_p + check_d]
    if x != position[check_p][check_d] {
      t.Error("position() error: expected ", position[check_p][check_d], ", got ", x, ".")
    }

    v := np.velocity()[particles * check_p + check_d]
    if v != velocity[check_p][check_d] {
      t.Error("velocity() error: expected ", velocity[check_p][check_d], ", got ", v, ".")
    }

    //This next part checks newposition(), newvelocity(), and postStep()
    newx := np.newPosition()[particles * check_p + check_d]
    newv := np.newVelocity()[particles * check_p + check_d]
    np.postStep()
    newx_test := np.position()[particles * check_p + check_d]
    newv_test := np.velocity()[particles * check_p + check_d]
    x_test := np.newPosition()[particles * check_p + check_d]
    v_test := np.newVelocity()[particles * check_p + check_d]

    if newx != newv_test {
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
    */
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

func TestHarmonicOscillatorForceCalculation(t *testing.T) {
  for j := 0; j < 10; j ++ {
    dimension := randInt(1, 5)
    p1 := make([]float64, dimension)
    p2 := make([]float64, dimension)

    for i := 0; i < dimension; i++ {
      p1[i] = randFloat(-10, 10)
      p2[i] = randFloat(10, 10)
    }

    Klen := randInt(0, 5) //The number of terms in the force function. 
    K := make([]float64, Klen)

    var d float64 = 1
    for i := 0; i < Klen; i++ {
      K[i] = randFloat(0, 3) / d //Each force constant gets smaller and smaller. 
      d *= 2
    }

    distv := make([]float64, dimension)
    quad := quadrance(p1, p2, distv)

    //check if distv was calculated correctly. 
    vector_check := randInt(0, dimension - 1)
    var expected float64 = p1[vector_check] - p2[vector_check]
    if distv[vector_check] != expected {
      t.Error("Quadrance error: distancev was ", distv[vector_check], "; expected ", expected, ".")
    }

    //What about quad? 
    var expected_quad float64 = 0
    for i:= 0; i < dimension; i++ {
      expected_quad += distv[i] * distv[i]
    }
    if quad != expected_quad {
      t.Error("Quadrance error: quad was ", quad, "; expected ", expected_quad)
    }

    //Finally, check the force. 
    force := oscillatorForce(quad, K) 
    var expected_force float64 = 0
    var temp float64 = quad
    for i := 0; i < len(K); i++ {
      expected_force += temp * K[i]
      temp *= quad
    }

    if force != expected_force {
      t.Error("oscillatorForce error: force was ", force, "; expected ", expected_force)
    }
  }
}

/*
func TestHarmonicOscillatorAcceleration(t *testing.T) {
  for i := 0; i < 10; i ++ {
    np, dimensions, particles, position, velocity [][]float64 := randomParticleSystem()

    mass = make([]float64, particles)

    for j := 0 j < particles; j++ {
      mass[j] = randFloat(1, 3);
    }

    K := make([][][]float64, particles)

    //Test a particle system with only quadradic force constants. 
    for j := 1; j < particles; j++ {
      K[j] = make([][]float64, j)
      for k := 0; k < j; k++ { 
        K[j][k] := make([]float64, 1)
        K[j][k][0] = randFloat(0, 3)
      } 
    }

    //Check that the acceleration of a random particle is correct.
    check = randInt(0, particles - 1)
    

    K := make([][][]float64, particles)

    //Test again with up to 4 force constants. 
    for j := 1; j < particles; j++ {
      K[j] := make([][]float64, j)
      for k := 0; k < j; k++ { 
        max := randInt(0, 4)
        var d float64 = 1
        K[j][k] := make([]float64, max)

        for m := 0; m < max; m++ {
          K[j][k][0] = randFloat(0, 3) / d
          d *= 2
        }
      } 
    }
  }
}*/
