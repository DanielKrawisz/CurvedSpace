package pathtrace

//This might be updated to be more of an interface or whatever. 
type LightRay struct {
  //The number of steps the ray has taken.
  depth int
  //The ray. 
  position, direction []float64
  //The wavelengths of light that this ray tracks.
  receptor []float64
  //The intensities of those wavelengths. 
  color []float64
  //As the ray bounces along, if it encounters glowing objects,
  //these numbers keep track of how the final color should be
  //adjusted to take these earlier interactions into account.
  emission []float64
  redirected float64 
}

func (r *LightRay) Trace(u float64) {
  for i := 0; i < 3; i ++ {
    r.position[i] = r.position[i] + u * r.direction[i]
  }
}

//From the parameters inside the ray object, derive the color of the ray.
func (r *LightRay) DeriveColor() (c []float64) {
  c = make([]float64, len(r.receptor))
  for i := 0; i < len(r.receptor); i ++ {
    c[i] = r.color[i] * r.redirected + r.emission[i]
  }
  return 
}

//A function to make an object glow.
func (r *LightRay) Glow(c []float64) {
  for i := 0; i < 3; i ++ {
    r.emission[i] = r.color[i] * c[i] * r.redirected + r.emission[i]
  }
  r.redirected = 0
}

//A function to make an object absorb light.
func (r *LightRay) Absorb(c []float64) {
  for l := 0; l < 3; l ++ {
    r.color[l] *= c[l]
  }
}

//A function for both.
func (r *LightRay) GlowAbsorbAverage(glow_color, transmit_color []float64, absorb float64) {
  for i := 0; i < 3; i ++ {
    r.emission[i] += r.redirected * absorb * glow_color[i]
    r.color[i] *= transmit_color[i]
  }
  r.redirected *= (1 - absorb)
}
