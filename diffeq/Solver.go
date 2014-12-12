package diffeq

//The abstract components of a differential equation solver. 

import "container/list"
import "math"

//A single instant in a numerical solution to a differential equation. 
type Instant struct{
  Region int
  S float64
  X[] float64
}

//Represents an initial condition to a differential equation. 
type InitialCondition interface {
  Region() int
  State() []float64
  Time() float64
}

type Derivative interface {
  //Sets the derivative of a state to v in state x. 
  DxDs(x []float64, v []float64)
}

type Step interface {
  //Takes one step at a time in the numerical differential equation solver. 
  //Updates the state and estimated error based on the derivative and the error parameter. 
  Step(x State, f Derivative) error
  //The number of dimensions. Must be compatible with the state and the Monitors. 
  Dimension() int
}

//The monitor is used to manage discrete changes to the differential equation.
//For example, coordinate changes or solid objects. 
type Monitor interface {
  //updates the state of a differential equation in accordance with the given acceleration
  //and step type. The state of the diffeq be given when the monitor object is initialized
  //because there are different kinds of differential equation states which monitors must
  //be adapted to. 
  update(f Derivative, step Step) 
  end() bool //Reports whether it has detected that an ending condition has been reached. 
}

//A differential equation solver.  
type Solver interface {
  Run() (i *Instant, l *list.List)
}

//A solver that supports step monitors. 
type solverStepMonitor struct {
  maxsteps int
  maxarclength float64
  exporthistory bool
  x State
  f Derivative
  step Step
  monitor []Monitor
}

//Runs a system of differential equations. 
func (sol *solverStepMonitor) Run() (i *Instant, l *list.List) {
  var steps int = 0

  //A linked list to export the full history of the evolution.
  //Otherwise, just the last step is returned. 
  if(sol.exporthistory) {
    l = list.New()
  }

  //It is necessary to compute the first derivative because the 
  //repeating sequence of steps assumes that the derivative was
  //computed at the end of the last cycle. This is necessary in
  //order to take advantage of the first-same-as-last property
  //that some RK methods have. 
  sol.f.DxDs(sol.x.position(), sol.x.velocity())

  for {
    if(sol.exporthistory) {
      l.PushBack(sol.x.ExportInstant())
    }

    steps ++

    sol.step.Step(sol.x, sol.f)

    for _, m := range sol.monitor {
      m.update(sol.f, sol.step) 
      if(m.end()) {goto end}
    }

    if (steps >= sol.maxsteps) {break}
    if (sol.x.S() > sol.maxarclength) {break}

    sol.x.postStep()
  }
  end:

  sol.x.postStep()
  i = sol.x.ExportInstant()
  if(sol.exporthistory) {
    l.PushBack(sol.x.ExportInstant())
  }

  return i, l
}

//TODO: create a solver for path tracing. 
//A solver for doing path tracing. 
/*type solverPathTracer struct {
  maxinteractions int //Maximum number that the ray can interact with objects.
  maxsteps int        //maximum number of differential equation steps. 
  exporthistory bool
  x State
  step Step
}*/

type State interface {
  //The arc length parameter which defines how far the 
  //simulation has gone. In Newtonian physics this would
  //just be time, but more generally and in relativistic
  //physics it is an arc length.
  S() float64
  //The step size. 
  Ds() float64
  setDs(dt float64);
  //Set the next ds.
  nextDs(ds float64)
  //The number of elements in the arrays.
  Length() int
  //get the current instant. 
  ExportInstant() *Instant
  //The position data of the current state.
  position() []float64
  //The immediate change at this point.
  velocity() []float64
  //The velocity for the next state.
  newVelocity() []float64
  //The position for the next state.
  newPosition() []float64
  //A function to call after taking a step to prepare for the next one.
  postStep() 
  //The estimated error.
  errorEstimate() []float64
}

type state struct {
  n int // equal to dim * particles
  region int
  s, ds, newds float64
  pos, vel, newpos, newvel, err []float64
}

func (p *state) Length() int {
  return p.n;
}

func (p *state) S() float64 {
  return p.s;
}

func (p *state) Ds() float64 {
  return p.ds;
}

func (p *state) setDs(ds float64) {
  p.ds = ds;
}

func (p *state) nextDs(ds float64) {
  p.newds = ds;
}

func (p *state) position() []float64 {
  return p.pos;
}

func (p *state) velocity() []float64 {
  return p.vel;
}

func (p *state) newPosition() []float64 {
  return p.newpos;
}

func (p *state) newVelocity() []float64 {
  return p.newvel;
}

func (p *state) errorEstimate() []float64 {
  return p.err;
}

func (p *state) ExportInstant() *Instant {
  x := make([]float64, len(p.pos))
  for i := 0; i < len(x); i++ {
    x[i] = p.pos[i]
  }
  return &Instant{p.region, p.s, x}
}

//Swap the new points with the old one to begin the next step.
func (p *state) postStep() {
  var tmp1 []float64

  tmp1 = p.pos
  p.pos = p.newpos
  p.newpos = tmp1
  tmp1 = p.vel
  p.vel = p.newvel
  p.newvel = tmp1

  p.s += p.ds
  p.ds = p.newds
}

//A Monitor that stops the simulation at exactly a specific time. 
//Implements Monitor
type untilTime struct {
  endb bool
  endTime, err float64
  x State
}

func (m *untilTime) update(f Derivative, step Step) {
  if m.x.S() >= m.endTime - m.err {
    m.endb = true
    return 
  }
  if m.x.S() + m.x.Ds() >= m.endTime {
    m.x.setDs(m.endTime - m.x.S())
    
    step.Step(m.x, f)
    m.endb = true
  }
}

func (m *untilTime) end() bool {
  return m.endb
}

//A constructor for an UntilTime monitor. 
//Can return nil! 
func NewUntilTime(x State, endTime, err float64) *untilTime{
  if math.IsNaN(endTime) || math.IsInf(endTime, 0) {return nil}
  if math.IsNaN(err) || math.IsInf(err, 0) || err <= 0 {return nil}
  if x == nil {return nil}

  return &untilTime{false, endTime, err, x}
}
