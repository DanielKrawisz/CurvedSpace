package diffeq

//A model of particles interacting with spring forces. This mainly exists to test the
//Runge Kutta differential equation solving methods in a variety of controlled conditions. 

import "math"
//import "fmt"

//A simple implementation of State.
type newtonianParticle struct {
  state
  dim int //Dimensions of the space in which the paricles move.
  particles int //Number of particles in the system.
  x, v, a, newx, newv, newa [][]float64
}

//Swap the new points with the old one to begin the next step.
func (p *newtonianParticle) postStep() {
  var tmp2 [][]float64

  p.state.postStep()

  tmp2 = p.x
  p.x = p.newx
  p.newx = tmp2
  tmp2 = p.v
  p.v = p.newv
  p.newv = tmp2
  tmp2 = p.a
  p.a = p.newa
  p.newa = tmp2
}

//Returns a new newtonianParticle state. 
//Can return nil!
func NewNewtonianParticle(dimension, particles int, initS, initDs float64, initpos [][]float64, initvel [][]float64) *newtonianParticle {
  //Check for invalid parameters. 
  if dimension <= 0 {return nil}
  if particles <= 0 {return nil}
  if math.IsNaN(initS) || math.IsInf(initS, 0) {return nil}
  if math.IsNaN(initDs) || math.IsInf(initDs, 0) || initDs == 0.0 {return nil}
  if initpos == nil {return nil}
  if initvel == nil {return nil}
  if len(initpos) != particles {return nil}
  if len(initvel) != particles {return nil}

  //These arrays allow the state to act like some arbitarry differential equation.
  n := dimension * particles
  pos := make([]float64, 2 * n)
  vel := make([]float64, 2 * n)
  newpos := make([]float64, 2 * n)
  newvel := make([]float64, 2 * n)

  //These arrays allow the state to act like a set of particles. 
  x := make([][]float64, particles)
  v := make([][]float64, particles)
  a := make([][]float64, particles)
  newx := make([][]float64, particles)
  newv := make([][]float64, particles)
  newa := make([][]float64, particles)

  for i := 0; i < particles; i++ {
    //Check for invalid parameters
    if initpos[i] == nil {return nil}
    if initvel[i] == nil {return nil}
    if len(initpos[i]) != dimension {return nil}
    if len(initvel[i]) != dimension {return nil}

    x[i] = pos[dimension * i : dimension * (i + 1)]
    v[i] = pos[n + dimension * i : n + dimension * (i + 1)]
    a[i] = vel[n + dimension * i : n + dimension * (i + 1)]

    newx[i] = newpos[dimension * i : dimension * (i + 1)]
    newv[i] = newpos[n + dimension * i : n + dimension * (i + 1)]
    newa[i] = newvel[n + dimension * i : n + dimension * (i + 1)]

    for j := 0; j < dimension; j++ {
      if math.IsNaN(initpos[i][j]) || math.IsInf(initpos[i][j], 0) {return nil}
      if math.IsNaN(initvel[i][j]) || math.IsInf(initvel[i][j], 0) {return nil}
      x[i][j] = initpos[i][j]
      v[i][j] = initvel[i][j]
    }
  }

  return &newtonianParticle{
    state{2 * dimension * particles, 0, initS, initDs, initDs, pos, vel, newpos, newvel, make([]float64, 2 * n)},
    dimension, particles, x, v, a, newx, newv, newa}
}

//A simple derivative function
type harmonicOscillator struct {
  dimension, particles int
  mass []float64
  K [][][]float64
}

//Calculate the distance vector from p1 to p2.
func distanceVector(p1, p2, dist []float64) {
  for k := 0; k < len(p1); k++ {
    dist[k] = p2[k] - p1[k]
  }
}

//Calculate the quadrance between two particles.
func quadrance(dist []float64) float64 {
  var quad float64 = 0
  for k := 0; k < len(dist); k++ {
    quad += dist[k] * dist[k]
  }
  return quad
}

//Calculate the force between two particles. 
func oscillatorForce(dist []float64, K []float64) []float64 {

  var force []float64 = make([]float64, len(dist))

  //If there are no force constants, the force is zero.
  if(len(K) > 0) {

    //If there is just one, we don't need the quadrance. 
    for i := 0; i < len(dist); i++ {
      force[i] = dist[i] * K[0]
    }

    //If there are more force constants, then we do this. 
    if len(K) > 1 {
      var xpow float64 = 1
      var fc float64 = 0
      var quad float64 = quadrance(dist)

      for k := 1; k < len(K); k++ {
        xpow *= quad
        fc += xpow * K[k]
      }

      for i := 0; i < len(dist); i++ {
        force[i] += dist[i] * fc
      }
    }
  }

  return force
}

//Does not check that x and v are the right lengths. 
func (o *harmonicOscillator) DxDs(pos, vel []float64) {

  //First set all accelerations to zero and the new
  //velocities to the correct values. 
  n := o.particles * o.dimension 
  for i := 0; i < n; i++ {
    vel[i] = pos[n + i]
    vel[n + i] = 0.0
  }

  //Generate list of particle positions. 
  x := make([][]float64, o.particles) 
  for i := 0; i < o.particles; i++ {
    x[i] = pos[o.dimension * i : o.dimension * (i + 1) ]
  }

  //The distance vector between two particles. 
  var dist []float64 = make([]float64, o.dimension)

  //Iterate only over the lower half of the matrix
  for i := 1; i < len(o.K); i++ {
    for j := 0; j < i; j++ {

      distanceVector(x[i], x[j], dist)

      var force []float64 = oscillatorForce(dist, o.K[i][j])

      //Add the accelerations. 
      for k := 0; k < o.dimension; k++ {
        vel[o.particles * (o.dimension + i) + k] += force[k] / o.mass[i]
        vel[o.particles * (o.dimension + j) + k] += -force[k] / o.mass[j]
      }
    }
  }
}

//The masses of the particles and the set of spring constants between them.
//Can return nil! 
//  p    - the newtonianParticle object. 
//  mass - the masses of the particles. 
//  K    - the sets of spring constants between the particles. It is a square matrix in which only
//         the entries below the diagonal are used.
func NewHarmonicOscillator(p *newtonianParticle, mass []float64, K [][][]float64) *harmonicOscillator{
  if K == nil {return nil}
  if p == nil {return nil}
  if mass == nil {return nil}
  if len(mass) != p.particles {return nil}
  if p.particles != len(K) {return nil}

  for i := 0; i < p.particles; i++ {
    if math.IsNaN(mass[i]) || math.IsInf(mass[i], 0) || mass[i] <= 0 {return nil}
  }

  for i := 1; i < p.particles; i++ {
    if math.IsNaN(mass[i]) || math.IsInf(mass[i], 0) || mass[i] <= 0 {return nil}
    if K[i] == nil {return nil}
    if len(K[i]) < i {return nil}
    
    for j := 0; j < i; j++ {
      if K[i][j] == nil {return nil}
    }
  }

  return &harmonicOscillator{p.dim, p.particles, mass, K}
}

//Create a solver that is a system of springs. 
func NewSpringSystem(dimension, particles int, initpos, initvel [][]float64, mass []float64, K [][][]float64, ds, until float64, step Step) Solver {
  p := NewNewtonianParticle(dimension, particles, 0, ds, initpos, initvel)
  if p == nil {return nil}
  h := NewHarmonicOscillator(p, mass, K)
  if p == nil {return nil}
  u := NewUntilTime(p, until, .0000000001)
  if u == nil {return nil}
  if step == nil {return nil}
  if step.Dimension() != 2 * dimension * particles {return nil}

  return &solver{10000, 3000, true, p, h, step, []Monitor{u}}
}
