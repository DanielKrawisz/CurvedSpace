package polynomialsurfaces

import "github.com/DanielKrawisz/CurvedSpace/surface"
import "github.com/DanielKrawisz/CurvedSpace/surface/polynomialsurfaces/polynomials"
import "github.com/DanielKrawisz/CurvedSpace/surface/booleans"
import "fmt"
import "strings"
import "math"

//Use polynomial surfaces and solid constructive geometry to
//make some primitive shapes.

//TODO cone, paraboloid, hyperboloid

type sphere struct {
	dim int
	p   []float64
	r2  float64 //Must not be negative; may be zero or infinity.
	p2  float64 //A pre-computed number.
}

func (s *sphere) String() string {
	return strings.Join([]string{"sphere{", fmt.Sprint(s.p), ", ", fmt.Sprint(s.r2), "}"}, "")
}

func (s *sphere) R2() float64 {
	return s.r2
}

func (s *sphere) X() []float64 {
	return s.p
}

func (s *sphere) Dimension() int {
	return s.dim
}

func (s *sphere) F(x []float64) float64 {
	var d float64 = 0
	for i := 0; i < s.dim; i++ {
		z := x[i] - s.p[i]
		d += z * z
	}
	return s.r2 - d
}

func (s *sphere) Intersection(x, v []float64) []float64 {
	var px, xx, vv, vpx float64
	for i := 0; i < s.dim; i++ {
		vpx += v[i] * (x[i] - s.p[i])
		vv += v[i] * v[i]
		xx += x[i] * x[i]
		px += s.p[i] * x[i]
	}

	//v needs to have a length greater than zero, so don't need to take account of divide by zero.
	return polynomials.QuadraticFormula((s.p2-s.r2-2*px+xx)/vv, 2*vpx/vv)
}

//(x - p).(x - p) < r2
func (s *sphere) Interior(x []float64) bool {
	return s.F(x) > 0
}

// -2 (x - p)
func (s *sphere) Gradient(x []float64) []float64 {
	z := make([]float64, s.dim)
	for i := 0; i < s.dim; i++ {
		z[i] = -2 * (x[i] - s.p[i])
	}
	return z
}

func (s *sphere) Translate(x []float64) surface.Surface {
	for i := 0; i < s.dim; i++ {
		s.p[i] += x[i]
	}
	return s
}

func (s *sphere) CoordinateShift(x [][]float64) surface.Surface {
	// TODO
	return s
}

//May return nil
func NewSphere(p []float64, r float64) surface.Surface {
	if p == nil {
		return nil
	}
	var p2 float64 = 0.0
	for i := 0; i < len(p); i++ {
		p2 += p[i] * p[i]
	}
	return &sphere{len(p), p, r * r, p2}
}

//Given by the point at the center, n vectors forming an
//orthonormal set (though not required to be) and n parameters
//defining the axes of the ellipsoid. The result will only be
//an elipsoid if the parameters are all positive, but the
//function does not require them to be.
//May return nil
func NewEllipsoid(point []float64, vec [][]float64, param []float64) surface.Surface {
	if point == nil || vec == nil || param == nil {
		return nil
	}
	if len(point) != len(param) || len(point) != len(vec) {
		return nil
	}

	v := make([][]float64, len(point))

	for i := 0; i < len(point); i++ {
		if len(vec[i]) != len(point) {
			return nil
		}

		if param[i] < 0.0 {
			return nil
		}

		v[i] = make([]float64, len(point))

		for j := 0; j < len(point); j++ {
			v[i][j] = vec[i][j] / math.Sqrt(param[i])
		}
	}

	return NewQuadraticSurface(point, v, [][]float64{}, make([]float64, len(point)), 1)
}

//A surface that is infinite in some directions and finite in others.
//The vectors define the finite directions.
//May return nil.
func NewInfiniteCylinder(p []float64, vector [][]float64) surface.Surface {

	return NewQuadraticSurface(p, vector, [][]float64{}, make([]float64, len(p)), 1)
}

//Vectors are made to be orthonormal.
func NewInfiniteHyperboloid(p []float64, vp, vn [][]float64) surface.Surface {
	return NewQuadraticSurface(p, vp, vn, []float64{}, 1)
}

//Given by the point at the apex of the paraboloid,
//a set of vectors defining the symmetric tensor and a set
//defining the vector part of the quadratic surface.
//May return nil.
func NewInfiniteParaboloid(p []float64, vc [][]float64, vb []float64) surface.Surface {
	return NewQuadraticSurface(p, vc, [][]float64{}, vb, 0)
}

//The first set of vectors define what is inside the cone, the rest define
//what is outside.
func NewInfiniteCone(p []float64, vp [][]float64, vn [][]float64) surface.Surface {
	return NewQuadraticSurface(p, vp, vn, []float64{}, 0)
}

//The intersection of two conic sections into something which has a finite area
//if the vectors defining it are non zero.
func NewCompoundConic(p []float64, vp, vn [][]float64) surface.Surface {
	//TODO
	return nil
}

//The intersection of two infinite cylinder objects. In 3 dimensions, this
//just becomes a regular cylinder. Param gives the radii of the cylinder
//in each direction. (so you can have an elliptical cylinder too)
//May return nil
func NewCylinder(p []float64, vp, vn [][]float64, param []float64) surface.Surface {
	if p == nil || vp == nil || vn == nil || param == nil {
		return nil
	}

	dim := len(p)
	if len(vp)+len(vn) != dim || len(param) != dim {
		return nil
	}

	for i := 0; i < len(vp); i++ {
		for j := 0; j < dim; j++ {
			vp[i][j] *= param[i]
		}
	}

	for i := 0; i < len(vp); i++ {
		for j := 0; j < dim; j++ {
			vn[i][j] *= param[i+len(vp)]
		}
	}

	return booleans.NewIntersection(NewInfiniteCylinder(p, vp), NewInfiniteCylinder(p, vn))
}

func NewCone(p []float64, axis []float64, v [][]float64, param []float64) surface.Surface {
	if p == nil || param == nil || axis == nil || v == nil {
		return nil
	}

	dim := len(p)
	if len(p) != len(param)+1 || len(param) != len(v) || len(axis) != dim {
		return nil
	}

	//TODO
	return nil
}

func NewConoid(p []float64, vp, vn [][]float64, param []float64) surface.Surface {
	//TODO
	return nil
}

func NewParaboloid(p []float64, vp, vn [][]float64, param []float64) surface.Surface {
	//TODO
	return nil
}

func NewHyperboloid(p []float64, vp, vn [][]float64, param []float64) surface.Surface {
	//TODO
	return nil
}
