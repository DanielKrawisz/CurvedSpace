package diffeq

import "testing"
import "math"
import "fmt"

//The Euler Method is here for testing purposes.
/*type eulerMethod struct {
}

func (*eulerMethod) Step(st State, f Derivative) error {
  for i := 0; i < len(st.pos); i++ {
    st.newpos[i] = 
  }

  
}*/

func getRKMethods(n int, errscale float64) (rkmethods []*rungeKuttaStepSizer) {
  return []*rungeKuttaStepSizer{
    NewRungeKuttaSolverMethodEuler(n, errscale),
    NewRungeKuttaSolverMethodMidpoint(n, errscale),
    NewRungeKuttaSolverMethodKutta(n, errscale),
    NewRungeKuttaSolverMethodBogackiSchampine(n, errscale),
    NewRungeKuttaSolverMethodFehlberg(n, errscale),
    NewRungeKuttaSolverMethodCashCarp(n, errscale),
    NewRungeKuttaSolverMethodDormandPrince(n, errscale)}
}

var rknames = []string{"Euler", "Midpoint", "Kutta", "BogackiSchampine", "Fehlberg", "CashCarp", "DormandPrince"}

//Test to make sure that the runge kutta objects load properly. 
func TestRungeKuttaObjects(t *testing.T) {

  rkmethods := getRKMethods(2, 1)

  if rkmethods == nil {t.Error("Why isn't this working??")}

  for i := 0; i < len(rkmethods); i++ {
    if rkmethods[i] == nil {
      t.Error("RK method ", rknames[i], " is nil. ")
    } else {
      if rkmethods[i].Name() != rknames[i] {t.Error("Wrong RK method loaded")}
    }
  }
}

//We test by checking that the simulation becomes more accurate as the step size decreases
//And as the order of the RK method increases. 
//TODO: this test does not pass! 
func TestRungeKuttaHarmonicOscillatorNoStepSizing(t *testing.T) {
  rkmethods := getRKMethods(8, 1)
  stepIncrements := []int{10}//, 30, 100, 300, 1000}

  position := [][]float64{[]float64{-1, 0}, []float64{1, 0}}
  velocity := [][]float64{[]float64{0, .5}, []float64{0, -.5}}
  //We test by checking that the methods approach the correct value as the error scale is decreased.
  for i := 0; i < len(rkmethods); i++ {
    last_quad_sum := math.Inf(1)

    for _, stepIncrement := range stepIncrements {

      fmt.Println("**** trial: ", i, "; step increment: ", stepIncrement)

      //Two particles in a simple orbit around one another. 
      particle := NewNewtonianParticle(2, 2, 0,
        2 * math.Pi / float64(stepIncrement), 
        position, velocity)

      spring := NewHarmonicOscillator(particle, []float64{1, 1}, [][][]float64{nil, [][]float64{[]float64{1}}})

      fmt.Println("**** system created. Length: ", particle.Length())

      spring.DxDs(particle.position(), particle.velocity())
      //After the number of step increments, the system should be in its original state. 
      for j := 0; j < stepIncrement; j++ {    
        rkmethods[i].rkstep(particle, spring)
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
  }
}

//In this test, we add a nonlinear term and we test by ensuring that all the results end up 
//near to one another. 
func TestRungeKuttaHarmonicOscillatorNonlinearNoStepSizing(t *testing.T) {
}

//We do not need to test that the step sizing function produces an accurate result because 
//that will happen if the rkstep function is correct. We need to test, rather, that the step
//size changes in response to the proper conditions. 

//Test that the number of steps increases as the errscale increases. 

//Test that the 
