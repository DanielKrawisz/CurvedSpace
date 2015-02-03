package functions

import "testing"

func TestConstantFunction(t *testing.T) {
  z := 2.616
  if z != ConstantFunction(z)([]float64{0,0,0}) {
    t.Error("constant function error")
  }
}

func TestChecksFunction(t *testing.T) {
  //TODO
}
