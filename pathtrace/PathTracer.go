package pathtrace

import "../diffeq"
import "math"

//TODO: create a solver for path tracing. 
//A solver for doing path tracing. 
type solverPathTracer struct {
  maxinteractions int //Maximum number that the ray can interact with objects.
  maxsteps int        //maximum number of differential equation steps. 
  exporthistory bool
  x diffeq.State
  step diffeq.Step
  objects []Object
  cd diffeq.CompoundDerivative
}

func (s *solverPathTracer) Run() {
  
  var min_u float64
  var interaction Object
  if s.cd.Length() > 0 {
    min_u = 1
  } else {
    min_u = math.Inf(1)
  }

  /*for _, object := range s.objects {
    intersection := object.Intersection()
    for _, u := range intersection {
      if u > 0 && u < min_u {
        min_u = u
        interaction = object
      }
    }
  }*/

  if interaction == nil {
    
  }
}

//TODO first need a LightRay, which implements state.
//Then need Material, then Object. 
//Then can finally write the algorithm. 


