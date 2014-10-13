package diffeq

//Runge Kutta is a broad class of numerical differential equation
//solvers. These are all variants that support error estimation,
//which allows for the step size to be shortened and lengthened. 

import "math"

type RungeKuttaStepSizer interface {
  //The order of the method. 
  Order() int
  //The name of the method. 
  Name() string
  //Take one rk step without estimating step size. 
  //Updates the newstate and estimated error. 
  rkstep(err []float64)
}

//Implements Step
type rungeKuttaStepSizer struct {
  n int
  name string //The type of Runge Kutta method. 
  steps int //The number of steps in this method.
  order int //The order of this method.
  fsal bool //first-same-as-last property
  a, b, c, db []float64 //The Runge Kutta parameters
  tmp []float64 //for storing temporary data as the calculation proceeds.
  K [][]float64 // 
  errscale float64//A parameter limiting the allowed estimated error per step. 
}

func (rk *rungeKuttaStepSizer) Dimension() int {
  return rk.n
}

func (rk *rungeKuttaStepSizer) Name() string {
  return rk.name
}

func (rk *rungeKuttaStepSizer) Order() int {
  return rk.order
}

//Take a single runge kutta step. 
func (rk *rungeKuttaStepSizer) rkstep(st State, f Derivative) {
  var i, j, k, apos int
  n := st.Length()
  ds := st.Ds()
  pos := st.position()
  vel := st.velocity()
  newpos := st.newPosition()
  err := st.errorEstimate()

  //Used to reference the right part of the a parameter.
  apos = 0; 

  //First the intermediate steps are calculated.
  for i = 0; i < rk.steps; i++ {
    for k = 0; k < n; k++ {
      rk.tmp[k] = 0
      for j = 1; j <= i; j++ {
        rk.tmp[k] += rk.a[apos + j]*rk.K[j-1][k]
      }
      rk.tmp[k] = pos[k] + ds*(rk.tmp[k] + rk.a[apos]*vel[k])
    }
    f.DxDs(rk.tmp, rk.K[i])
    apos += rk.steps
  }

  //Then the output and error are calculated. 
  for k = 0; k < n; k++ {
    newpos[k] = 0
    err[k] = 0
    for j = 1; j <= rk.steps; j++ {
      newpos[k]   += rk.b[j]*rk.K[j-1][k]
      err[k] += rk.db[j]*rk.K[j-1][k]
    }

    newpos[k] = pos[k] + ds*(newpos[k] + rk.b[0]*vel[k])
    err[k] = ds*(err[k] + rk.db[0]*vel[k])
  }
}

//calculate the derivative of the next step. 
//Assumes that K has already been calculated. 
func (rk *rungeKuttaStepSizer) nextVelocity(st State, f Derivative) {
  newvel := st.newVelocity()
  if rk.fsal {
    for i := 0; i < st.Length(); i++ {
      newvel[i] = rk.K[rk.steps-1][i]
    }
  } else {
    f.DxDs(st.newPosition(), newvel)
  }
}

//Do a runge kutta step without step sizing. For testing purposes.
func (rk *rungeKuttaStepSizer) StepNoResize(st State, f Derivative) {
  rk.rkstep(st, f)
  rk.nextVelocity(st, f)
}

type UnderflowError struct{
}

func (e *UnderflowError) Error() string {
  return "Underflow Error"
}

type OverflowError struct{
}

func (e *OverflowError) Error() string {
  return "Overflow Error"
}

const safety = 0.9
const pgrow = -0.2
const pshrink = -0.25
const errcon = 1.89

func (rk *rungeKuttaStepSizer) Step(st State, f Derivative) error{
  var i int
  var errmax, accelmax, dstemp float64

  n := st.Length()
  s := st.S()
  ds := st.Ds()
  vel := st.velocity()
  err := st.errorEstimate()

  //Repeat until the estimated error is within an acceptable limit. 
  for {
    rk.rkstep(st, f)
    errmax=0.0; accelmax=0.0

    for i = 0; i < n; i++ {
      accelmax = math.Max(accelmax, math.Abs(vel[i]))
    }

    for i = 0; i < n; i++ {
      errmax = math.Max(errmax, math.Abs(err[i]/(ds*rk.errscale*accelmax)))
    }

    //The max estimated error is not too great, so we can proceed. 
    if (errmax <= 1.0) {break}

    dstemp = safety*ds*math.Pow(errmax, pshrink)
    if ds >= 0.0 {
      ds = math.Max(dstemp, 0.1*ds)
    } else {
      ds = math.Min(dstemp, 0.1*ds)
    }

    if s + ds == s {
      return &UnderflowError{}
    } 
    if math.IsNaN(ds) || math.IsInf(ds, 0) {
      return &OverflowError{}
    }
  }

  rk.nextVelocity(st, f)

  //Whether to grow or shrink the step size. 
  if(errmax > errcon) {
    st.nextDs(safety*ds*math.Pow(errmax, pgrow))
  } else {
    st.nextDs( ds * 5.0 )
  }

  return nil
}

func (rk *rungeKuttaStepSizer) rkinitialize(n int, errscale float64) {
  rk.errscale = errscale
  //Set up some temporary variables for the intermediate parts
  //of the computation. 
  rk.tmp = make([]float64, n)
  rk.K = make([][]float64, rk.steps)
  for i := 0; i < rk.steps; i++ {
    rk.K[i] = make([]float64, n)
  }
}

func NewRungeKuttaSolverMethodEuler(n int, errscale float64) *rungeKuttaStepSizer{
  rk := new(rungeKuttaStepSizer)
  rk.name = "Euler"
  rk.steps = 0
  rk.order = 1
  rk.fsal = false
  rk.a = []float64{0.}
  rk.b = []float64{1.}
  rk.c = []float64{0.}
  d := []float64{  0}
  rk.db = []float64{rk.b[0] - d[0]}
  rk.rkinitialize(n, errscale)
  return rk;
}

func NewRungeKuttaSolverMethodMidpoint(n int, errscale float64) *rungeKuttaStepSizer{
  rk := new(rungeKuttaStepSizer)
  rk.name = "Midpoint"
  rk.steps = 1
  rk.order = 2
  rk.fsal = false
  rk.a = []float64{2./3.}
  rk.b = []float64{.25,  .75}
  rk.c = []float64{2./3.}
  d := []float64{  1.,   0.}
  rk.db = []float64{rk.b[0] - d[0], rk.b[1] - d[1]}
  rk.rkinitialize(n, errscale)
  return rk;
}

func NewRungeKuttaSolverMethodKutta(n int, errscale float64) *rungeKuttaStepSizer{
  rk := new(rungeKuttaStepSizer)
  rk.name = "Kutta"
  rk.steps = 2
  rk.order = 3
  rk.fsal = false
  rk.a = []float64{.5,   0.,
                   -1.,  2.}
  rk.b = []float64{1/6., 2/3., 1/6.}
  rk.c = []float64{1/2., 1.}
  d := []float64{  0.,   1.,   0.}
  rk.db = []float64{rk.b[0] - d[0], rk.b[1] - d[1], rk.b[2] - d[2]}
  rk.rkinitialize(n, errscale)
  return rk;
}

func NewRungeKuttaSolverMethodBogackiSchampine(n int, errscale float64) *rungeKuttaStepSizer{
  rk := new(rungeKuttaStepSizer)
  rk.name = "BogackiSchampine"
  rk.steps = 3
  rk.order = 4
  rk.fsal = true
  rk.a = []float64{1/2.,  0,     0,
                   0.,    3./4., 0,
                   2/9.,  1/3.,  4/9.}
  rk.b = []float64{2/9.,  1/3.,  4/9.,  0}
  rk.c = []float64{1/2.,  3./4., 1.}
  d := []float64{  7/24., 1/4.,  1/3.,  1/8.}
  rk.db = []float64{rk.b[0] - d[0], rk.b[1] - d[1], rk.b[2] - d[2], rk.b[3] - d[3]}
  rk.rkinitialize(n, errscale)
  return rk;
}

func NewRungeKuttaSolverMethodFehlberg(n int, errscale float64) *rungeKuttaStepSizer{
  rk := new(rungeKuttaStepSizer)
  rk.name = "Fehlberg"
  rk.steps = 5
  rk.order = 5
  rk.fsal = false
  rk.a = []float64{1/4.,         0,            0,            0,            0,
                   3/32.,        9/32.,        0,            0,            0,
                   1932/2197.,   -7200./2197,  7296./2197,   0,            0,
                   439./216,     -8.,          3680/513.,    -845/4104.,   0,
                   -8/27.,       2.,           -3544/2565.,  1859/4104.,   -11/40}
  rk.b = []float64{25/216.,      0.,           1408/2565.,   2197/4104.,   -1/5.,        0}
  rk.c = []float64{1/4.,         3/8.,         12/13.,       1,            1/2.}
  d := []float64{  16/135.,      0.,           6656/12825.,  28561/56430., -9/50.,       2/55.}
  rk.db = []float64{rk.b[0] - d[0], rk.b[1] - d[1], rk.b[2] - d[2], rk.b[3] - d[3], rk.b[4] - d[4], rk.b[5] - d[5]}
  rk.rkinitialize(n, errscale)
  return rk;
}

func NewRungeKuttaSolverMethodCashCarp(n int, errscale float64) *rungeKuttaStepSizer{
  rk := new(rungeKuttaStepSizer)
  rk.name = "CashCarp"
  rk.steps = 5
  rk.order = 5
  rk.fsal = false
  rk.a = []float64{1/5.,          0.,            0.,            0.,            0.,
                   3/40.,         9/40.,         0.,            0.,            0.,
                   3/10.,         -9/10.,        6/5.,          0.,            0.,
                  -11/54.,       5/2.,          -70./27,       35/27.,        0.,
                   1631/55296.,   175/512.,      575/13824.,    44275/110592., 253/4096.}
  rk.b = []float64{37/378.,       0.,            250/621.,      125./594,      0.,            512/1771.}
  rk.c = []float64{1/5.,          3/10.,         3/5.,          1.,            7/8.}
  d := []float64{  2825/27648.,   0.,            18575/48384.,  13525/55296.,  277/14336.,    1/4.}
  rk.db = []float64{rk.b[0] - d[0], rk.b[1] - d[1], rk.b[2] - d[2], rk.b[3] - d[3], rk.b[4] - d[4], rk.b[5] - d[5]}
  rk.rkinitialize(n, errscale)
  return rk;
}

func NewRungeKuttaSolverMethodDormandPrince(n int, errscale float64) *rungeKuttaStepSizer{
  rk := new(rungeKuttaStepSizer)
  rk.name = "DormandPrince"
  rk.steps = 6
  rk.order = 5
  rk.fsal = true
  rk.a = []float64{1/5.,          0.,            0.,            0.,            0.,            0.,
                   3/40.,         9/40.,         0.,            0.,            0.,            0.,
                   44/45.,        -56/15.,       32/9.,         0.,            0.,            0.,
                   19372/6561.,   -25360/2187.,  64448/6561.,   -212/729.,     0.,            0.,
                   9017/3168.,    -355/33.,      46732/5247.,   49/176.,      -5103/18656.,   0.,
                   35/384.,       0.,            500/1113.,     125/192.,     -2187/6784.,    11./84}
  rk.b = []float64{35/384.,       0.,            500/1113.,     125/192.,     -2187/6784.,    11./84,        0}
  rk.c = []float64{1/5.,          3/10.,         4/5.,          8/9.,         1.,             1.}
  d := []float64{  5179/57600.,   0.,            7571/16695.,   393/640.,     -92097/339200., 187/2100.,     1/40.}
  rk.db = []float64{rk.b[0] - d[0], rk.b[1] - d[1], rk.b[2] - d[2], rk.b[3] - d[3], rk.b[4] - d[4], rk.b[5] - d[5], rk.b[6] - d[6]}
  rk.rkinitialize(n, errscale)
  return rk;
}
