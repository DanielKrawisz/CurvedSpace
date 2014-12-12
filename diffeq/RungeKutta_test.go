package diffeq

//Tests for the Runge Kutta methods. 

import "testing"
import "math"
import "fmt"

//The Euler Method is here for testing purposes.
type eulerMethod struct {
}

func (*eulerMethod) Step(st State, f Derivative) error {
  newpos := st.newPosition()
  pos := st.position()
  vel := st.velocity()
  dt := st.Ds()
  for i := 0; i < len(pos); i++ {
    newpos[i] = pos[i] + dt * vel[i]
  }

  f.DxDs(newpos, st.newVelocity())

  return nil
}

func getRKMethods(n int, errscale float64) (rkmethods []*rungeKuttaStepSizer) {
  return []*rungeKuttaStepSizer{
    NewRungeKuttaSolverMethodMidpoint(n, errscale),
    NewRungeKuttaSolverMethodKutta(n, errscale),
    NewRungeKuttaSolverMethodBogackiSchampine(n, errscale),
    NewRungeKuttaSolverMethodFehlberg(n, errscale),
    NewRungeKuttaSolverMethodCashCarp(n, errscale),
    NewRungeKuttaSolverMethodDormandPrince(n, errscale)}
}

var rknames = []string{"Midpoint", "Kutta", "BogackiSchampine", "Fehlberg", "CashCarp", "DormandPrince"}

//Test to make sure that the runge kutta objects load properly. 
func TestRungeKuttaObjects(t *testing.T) {

  rkmethods := getRKMethods(2, 1)

  if rkmethods == nil {t.Error("Why isn't this working??")}

  for i := 0; i < len(rkmethods); i++ {
    if rkmethods[i] == nil {
      t.Error("RK method ", rknames[i + 1], " is nil. ")
    } else {
      if rkmethods[i].Name() != rknames[i] {t.Error("Wrong RK method loaded")}
    }
  }
}

//Test that the eulerMethod step function works the same as the runge kutta version.
func TestEulerAgainstEuler(t *testing.T) {
  rkEuler := NewRungeKuttaSolverMethodEuler(8, 1)
  Euler := eulerMethod{}

  var step int = 10
  Klist := [][]float64{[]float64{1}, []float64{1, .2}}

  position := [][]float64{[]float64{-1, 0}, []float64{1, 0}}
  velocity := [][]float64{[]float64{0, 1}, []float64{0, -.5}}

  for i, K := range Klist {

    //Two particles in a simple orbit around one another. 
    particle1 := NewNewtonianParticle(2, 2, 0,
      math.Sqrt(float64(2.)/float64(3.)) * 2 * math.Pi / float64(step), 
      position, velocity)
    particle2 := NewNewtonianParticle(2, 2, 0,
      math.Sqrt(float64(2.)/float64(3.)) * 2 * math.Pi / float64(step), 
      position, velocity)

    spring1 := NewHarmonicOscillator(particle1, []float64{1, 2}, [][][]float64{nil, [][]float64{K}})
    spring2 := NewHarmonicOscillator(particle2, []float64{1, 2}, [][][]float64{nil, [][]float64{K}})

    spring1.DxDs(particle1.position(), particle1.velocity())
    spring2.DxDs(particle2.position(), particle2.velocity())
    //After the number of step increments, the system should be in its original state. 
    for j := 0; j < step; j++ {    
      rkEuler.StepNoResize(particle1, spring1)
      Euler.Step(particle2, spring2)
      particle1.postStep()
      particle2.postStep()

      //The two results should always be very close to one another. 
      var err float64 = 0;
      for k := 0; k < 2; k++ {
        for l := 0; l < 2; l++ {
          err += math.Abs(particle1.x[k][l] - particle2.x[k][l])
        }
      }
      if err > .0000001 {
        t.Error("Euler methods disagree on trial ", i, ", step ", j)
      }
    }
  }
}

//We test by checking that the simulation becomes more accurate as the step size decreases
//And as the order of the RK method increases. 
//TODO: this test does not pass! 
func TestRungeKuttaHarmonicOscillatorNoStepSizing(t *testing.T) {
  rkmethods := getRKMethods(8, 1)
  stepIncrements := []int{10, 30, 100, 300, 1000, 3000, 10000}

  position := [][]float64{[]float64{-1, 0}, []float64{1, 0}}
  velocity := [][]float64{[]float64{0, 1}, []float64{0, -.5}}
  //We test by checking that the methods approach the correct value as the error scale is decreased.
  for i := 0; i < len(rkmethods); i++ {
    last_quad_sum := math.Inf(1)
    var particle *newtonianParticle

    for _, stepIncrement := range stepIncrements {

      fmt.Println("**** trial: ", i, "; step increment: ", stepIncrement)

      //Two particles in a simple orbit around one another. 
      particle = NewNewtonianParticle(2, 2, 0,
        math.Sqrt(float64(2.)/float64(3.)) * 2 * math.Pi / float64(stepIncrement), 
        position, velocity)

      spring := NewHarmonicOscillator(particle, []float64{1, 2}, [][][]float64{nil, [][]float64{[]float64{1}}})

      //fmt.Println(particle.ExportInstant())
      spring.DxDs(particle.position(), particle.velocity())
      //After the number of step increments, the system should be in its original state. 
      for j := 0; j < stepIncrement; j++ {    
        rkmethods[i].StepNoResize(particle, spring)
        particle.postStep()
        //fmt.Println(particle.ExportInstant())
      }

      fmt.Println("**** completed the steps. ")

      //We now check that the system is closer to its original position than in the previous trial. 
      var quad_sum float64 = 0
      dist := make([]float64, 2)
      distanceVector(particle.x[0], position[0], dist)
      quad_sum += quadrance(dist)
      distanceVector(particle.x[1], position[1], dist)
      quad_sum += quadrance(dist)

      if quad_sum >= last_quad_sum {
        t.Error("trial: ", i, "; step increment: ", stepIncrement, "; Last distance from expected location ",
        last_quad_sum, ", should be larger than current, which is ", quad_sum, ".")
      }

      last_quad_sum = quad_sum
    }

    finalPosition := particle.x
    var err float64 = 0
    for j := 0; j < 2; j++ {
      for k := 0; k < 2; k++ {
        err += math.Abs(finalPosition[j][k] - position[j][k])
      }
    }

    if err > .000001 {
      t.Error("trial: ", i, "; the total error is too great.")
    }
  }
}

//In this test, we add a nonlinear term and we test by ensuring that all the results end up 
//near to one another. 
func TestRungeKuttaHarmonicOscillatorNonlinearNoStepSizing(t *testing.T) {
  rkmethods := getRKMethods(8, 1)
  steps := 1000
  finalPositions := make([][][]float64, len(rkmethods))

  position := [][]float64{[]float64{-1, 0}, []float64{1, 0}}
  velocity := [][]float64{[]float64{0, 1.}, []float64{0, -.5}}
  //We test by checking that the methods approach the correct value as the error scale is decreased.
  for i := 0; i < len(rkmethods); i++ {

    //Two particles in a nonlinear orbit around one another. 
    particle := NewNewtonianParticle(2, 2, 0,
      math.Sqrt(float64(2.)/float64(3.)) * 2 * math.Pi / float64(steps), 
      position, velocity)

    spring := NewHarmonicOscillator(particle, []float64{1, 1}, [][][]float64{nil, [][]float64{[]float64{1, .2}}})

    spring.DxDs(particle.position(), particle.velocity())
    //After the number of step increments, the system should be in its original state. 
    for j := 0; j < steps; j++ {    
      rkmethods[i].StepNoResize(particle, spring)
      particle.postStep()
      //fmt.Println(particle.ExportInstant())
    }

    fmt.Println("**** completed the steps. ")

    //We now check that the system is closer to its original position than in the previous trial. 
    var quad_sum float64 = 0
    dist := make([]float64, 2)
    distanceVector(particle.x[0], position[0], dist)
    quad_sum += quadrance(dist)
    distanceVector(particle.x[1], position[1], dist)
    quad_sum += quadrance(dist)

    finalPositions[i] = particle.x
  }

  //Now check that the final positions are nearby to one another. 
  for i := 0; i < len(finalPositions) - 1; i++ {
    for j := i + 1; j < i + 2 && j < len(finalPositions); j++ {
      var err float64 = 0

      for k := 0; k < 2; k ++ {
        for l := 0; l < 2; l ++ {
          err += math.Abs(finalPositions[i][k][l] - finalPositions[j][k][l]) 
        }
      }

      if err > .001 {
        t.Error("Trial ", i, " is too far away from trial ", j, ".")
      }
    }
  }
}

//We do not need to test that the step sizing function produces an accurate result because 
//that will happen if the rkstep function is correct. We need to test, rather, that the step
//size changes in response to the proper conditions. 
func TestRungeKuttaStepSizing(t *testing.T) {
  //TODO
}

//Test that the function returns the correct errors when the differential equation enters
//singular states. 
func TestRungeKuttaErrorConditions(t *testing.T) {
  //TODO
}
