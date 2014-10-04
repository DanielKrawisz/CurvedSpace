package diffeq

import "testing"

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

func TestRungeKuttaHarmonicOscillator(t *testing.T) {
  rkmethods := getRKMethods(2, 1)

  //We test by checking that the methods approach the correct value as the error scale is decreased.
  for i := 0; i < len(rkmethods); i++ {
    
  }
}
