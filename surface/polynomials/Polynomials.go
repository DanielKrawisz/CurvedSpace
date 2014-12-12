package polynomials

import "math"
import "math/cmplx"

//This file finds the real roots of polynomials of degrees 
// 2 through 4. Higher-degree polynomials can only be
// solved numerically in general. Higher degree polynomials
// can be solved if the Galois group is solvable. 

//Solutions for the equation
//
//  x^2 + a x + b == 0
//
//The parameters are named in a nonstandard way so
//that the formula generalizes more easily to qubics
//and quartics. It is assumed that all the coefficients
//are real, which means there are two cases. Either
//both solutions are complex or both are real. 
func QuadraticFormula(a, b float64) []float64 {
  desc := a*a - 4 * b
  if desc < 0 {
    return []float64{}
  } else {
    s := math.Sqrt(desc)
    return []float64{(a - s)/2., (a + s) / 2.}
  }
}

//Solutions for the equation
//
//  x^3 + b x + a == 0
//
//This is a cubic equation in which all roots sum
//to zero. There are two cases. Either there is one
//real solution or three. 
//
//If desc is negative, then there are three real roots
//and if it is positive, there is one. 
//TODO: make some constants for the numbers in here. 
func simplifiedCubicFormula(a, b float64) []float64 {
  desc := a * a / 4. + b * b * b / 27.
  //If the descriminant is negative, then there should be three real roots. 
  if desc < 0 {
    d := math.Sqrt(-desc)
    p := cmplx.Pow(complex(-a/2., d), 1.0/3.0)
    q := cmplx.Pow(complex(-a/2., -d), 1.0/3.0)
    cw := complex(-1./2., -math.Sqrt(3)/2.)
    ccw := complex(-1./2., math.Sqrt(3)/2.)
    return []float64{real(p) + real(q), real(cw * p) + real(ccw * q), real(ccw * p) + real(cw * q)}
  } else { //If it is positive, then there will only be one real root. 
    d := math.Sqrt(desc)
    r1 := -a/2. + d
    r2 := -a/2. - d
    var s1, s2 float64
    if r1 < 0 {
      s1 = -1;
    } else {
      s1 = 1;
    }
    if r2 < 0 {
      s2 = -1;
    } else {
      s2 = 1;
    }
    return []float64{s1*math.Pow(s1*r1, 1.0/3.0) + s2*math.Pow(s2*r2, 1.0/3.0)}
  }
}

//Solutions for the equation
//
//  x^3 + c x^2 + b x + a == 0
//
//TODO optimize the formula. 
func CubicFormula(a, b, c float64) []float64 {
  Z := simplifiedCubicFormula(a + 2*c*c*c/27. - (b*c)/3., b - c*c/3.)
  for i := 0; i < len(Z); i++ {
    Z[i] -= c/3.
  }
  return Z
}

//Solutions for the equation
//
//  x^4 + c x^2 + b x + a == 0
//
//Three possibilities: all solutions are complex, 
//two solutions are real, or four solutions are real. 
func simplifiedQuarticFormula(a, b, c float64) []float64 {
  //The cubic equation comes from attempting to find something that completes the square
  //so as to turn the quartic formula into a quadratic formula. It is always possible
  //because a cubic equation with real coefficients will always have at least one real root.
  //When there are zero or four real roots, then Z has three real roots, whereas
  // when there are only two real roots, Z has only one real root. This can be proven
  // by the fact that the descriminanent in the cubic formula comes out as a negative square
  // for the cases in which there are zero or four real roots, and comes out as a positive
  // square when there are two real roots. 
  //It doesn't matter which root is used, but it appears always to be possible to choose
  // a root which makes Ad and Bd positive. No proof for that, but the tests always pass.
  Zl := CubicFormula((c * c * c - a * c - b * b / 4.0)/2., 2. * c * c - a, (5./2.) * c)
  var Z, Bd float64

  for _, Z = range Zl{
    Bd = c + 2 * Z
    if Bd >= 0 { break }
  }

  oo := c + Z
  Ad := oo*oo - a

  //It must be the case that Ad and Bd be real and have the same sign. Proof: 
  //Since Z, c, and a are all real, it must be that Ad and Bd are real. 
  //It must be that b == (+/-) 2 A B, and since b is real, the product of A and B must be real. 

  A := math.Sqrt(Ad) //Ad and Bd now always positive (I think)
  B := math.Sqrt(Bd) 

  var bs float64
  if b < 0 {
    bs = -1
  } else {
    bs = 1
  } //Necessary to ensure that -b / 2 == A B. 

  //Ensure desc2 >= desc1. 
  var desc1, desc2, s0 float64
  if A >= 0 {
    desc1 = Bd/4. - A - oo
    desc2 = Bd/4. + A - oo
    s0 = 1
  } else {
    desc1 = Bd/4. + A - oo
    desc2 = Bd/4. - A - oo
    s0 = -1
  }

  if desc2 >= 0 {
    if desc1 >= 0 { //Four real solutions.
      d1 := math.Sqrt(desc1)
      d2 := math.Sqrt(desc2)
      return []float64{s0*bs*B/2. - d1, s0*bs*B/2. + d1, -s0*bs*B/2. + d2, -s0*bs*B/2. - d2}
    } else { //Two real solutions.
      d := math.Sqrt(desc2)
      return []float64{-s0*bs*B/2. + d, -s0*bs*B/2. - d} 
    }
  } else { //No real solutions.
    return []float64{}
  }
}

//Solutions for the equation
//
//  x^4 + d x^3 + c x^2 + b x + a == 0
//
//Three possibilities: all solutions are complex, 
//two solutions are real, or four solutions are real. 
func QuarticFormula(a, b, c, d float64) []float64 {
  Z := simplifiedQuarticFormula(a - b*d/4. + c*d*d/16. - 3*d*d*d*d/256., b - c*d/2. + d*d*d/8., c - 3*d*d/8)
  for i := 0; i < len(Z); i++ {
    Z[i] -= d/4
  }
  return Z
}
