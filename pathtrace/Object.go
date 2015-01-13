package pathtrace

import "../surface"

type Object interface {
  surface.Surface
  UpdateVector(l LightRay, u float64)
}
