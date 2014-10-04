package diffeq

import "math"

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

func (p *newtonianParticle) copyVel() {
  for i := 0; i < p.n; i++ {
    p.vel[i] = p.pos[p.n + i]
  }
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
    state{2 * dimension * particles, 0, initS, initDs, initDs, pos, vel, newpos, newvel, make([]float64, n)},
    dimension, particles, x, v, a, newx, newv, newa}
}

//A simple derivative function
type harmonicOscillator struct {
  p *newtonianParticle
  mass []float64
  K [][][]float64
}

//Calculate the quadrance between two particles.
func quadrance(p1, p2, distancev []float64) float64 {
  var quad float64 = 0
  for k := 0; k < len(p1); k++ {
    distancev[k] = p1[k] - p2[k]
    quad += distancev[k] * distancev[k]
  }
  return quad
}

//Calculate the force between two particles. 
func oscillatorForce(quad float64, K []float64) float64 {

  var force float64 = 0
  var xpow float64 = 1
  for k := 0; k < len(K); k++ {
    xpow *= quad
    force += xpow * K[k]
  }

  return force
}

//Does not check that x and v are the right lengths. 
func (o *harmonicOscillator) DxDs() {

  //First set all accelerations to zero and the new
  //velocities to the correct values. 
  o.p.copyVel();
  for i := o.p.n; i < 2*o.p.n; i++ {
    o.p.vel[i] = 0.0
  }

  //The distance vector between two particles. 
  var distancev []float64 = make([]float64, o.p.dim)

  for i := 0; i < len(o.K); i++ {
    for j := 1; j < i; j++ {
      //Calculate the quadrance between two particles.
      quad := quadrance(o.p.x[i], o.p.x[j], distancev)

      //If the distance is zero, then the force is zero.
      if quad == 0.0 {

        var distance float64 = math.Sqrt(quad)

        var force = oscillatorForce(quad, o.K[i][j])

        //Normalize the force over the distance. 
        //No possibility of divide by zero here. 
        force /= distance
        var avi float64 = force / o.mass[i]
        var avj float64 = force / o.mass[j]

        //Add the accelerations. 
        for k := 0; k < o.p.dim; k++ {
          o.p.a[i][k] += -avi / distancev[k]
          o.p.a[j][k] += avj / distancev[k]
        }
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

  return &harmonicOscillator{p, mass, K}
}

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
