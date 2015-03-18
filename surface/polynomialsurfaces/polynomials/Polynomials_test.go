package polynomials

import "testing"
import "sort"
import "github.com/DanielKrawisz/CurvedSpace/test"

//The error parameter.
var e float64 = .000001

//The general strategy for testing all polynomials is to divide
//the tests into cases based on the number of real roots that
//the polynomial has. Then the set of roots is created based on
//some random parameters and the full polynomial is constructed
//from those. Thus, the roots are known a priori. The function
//is then run to test whether it finds the same roots as those
//which we started with. 
func TestQuadratic(t *testing.T) {
  //Two cases: zero real roots or two roots. 
  //case 1, zero real roots. In this case the complex
  //roots are of the form p +/- i q.
  for i := 0; i < 10; i++ {
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)

    ans := QuadraticFormula(p * p + q * q, -2 * p)

    //ans should be empty. 
    if len(ans) != 0 {
      t.Error("quadratic formula error: case 1, expected no real roots for p = ", p , ", q = ", q, "; ans = ", ans)
    }
  }

  //case 2: two real roots. The roots are p and q.
  for i := 0; i < 10; i++ {
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)

    ans := QuadraticFormula(p * q, - p - q)

    if len(ans) != 2 {
      t.Error("quadratic formula error: case 2, expected two real roots for p = ", p , ", q = ", q, "; ans = ", ans)
    } else {
      test_list := []float64{p, q}
      sort.Float64s(test_list)
      sort.Float64s(ans)

      var close_enough bool = true
      for i := 0; i < len(test_list); i++ {
        close_enough = close_enough && test.CloseEnough(test_list[i], ans[i], e);
      }

      if !close_enough {
        t.Error("quadratic formula error: case 2, roots do not match for p = ", p , ", q = ", q, "; ans = ", ans)
      }
    }
  }
}

//For the cubic and quartic, we test the functions for the
//simplified formulas in the most general way possible. The
//full formula depends on the simplified one, so we do not
//test that one for a general polynomial. We only go through
//one case, and if it works for that it will work for all
//cases.
func TestCubic(t *testing.T) {

  //In the simplified cubic formala, all roots sum to zero. 
  //Two cases: one root or three roots.
  //case 1: In this case, the three roots are p, (-p + i q) / 2, (-p - i q) / 2. 
  for i := 0; i < 10; i++ {
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)

    ans := simplifiedCubicFormula(-(p*p*p + p*q*q)/4., (-3.*p*p + q*q)/4.)

    if len(ans) != 1 {
      t.Error("simplified cubic formula error: case 1, expected one real root for p = ",
        p , ", q = ", q, "; ans = ", ans)
    } else if !test.CloseEnough(ans[0], p, e) {
      t.Error("simplified cubic formula error: case 1, expected real roots to match p = ",
        p , ", q = ", q, "; ans = ", ans)
    }
  }

  //case 2: In this case the roots are p, q, and -p - q. 
  for i := 0; i < 10; i++ {
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)

    ans := simplifiedCubicFormula(p*p*q+p*q*q, -p*p-p*q-q*q)

    if len(ans) != 3 {
      t.Error("simplified cubic formula error: case 2, expected three real roots for p = ",
        p , ", q = ", q, "; ans = ", ans)
    } else {
      test_list := []float64{p, q, -p - q}
      sort.Float64s(test_list)
      sort.Float64s(ans)

      var close_enough bool = true
      for i := 0; i < len(test_list); i++ {
        close_enough = close_enough && test.CloseEnough(test_list[i], ans[i], e);
      }

      if !close_enough {
        t.Error("simplified cubic formula error: case 2, expected real roots to match p = ",
          p , ", q = ", q, "; ans = ", ans)
      }
    }

    //Now add a constant to the roots and test them against the more general formula. 
    r := test.RandFloat(-100, 100)
    ans2 := CubicFormula(p*p*q+p*q*q+p*p*r+p*q*r+q*q*r-r*r*r, -p*p-p*q-q*q + 3.*r*r, -3.*r)

    if len(ans2) != 3 {
      t.Error("cubic formula error: case 2, expected one real root for p = ",
        p , ", q = ", q, ", r = ", r, "; ans = ", ans2)
    } else {
      //move the new roots back to being centered at zero; this should put them in the same position
      //as the roots in the simplified formula. 
      for i := 0; i < 3; i++ {
        ans2[i] -= r;
      }
      sort.Float64s(ans2);

      var close_enough bool = true;
      for i := 0; i < len(ans); i++ {
        close_enough = close_enough && test.CloseEnough(ans[i], ans2[i], e)
      }

      if !close_enough {
        t.Error("cubic formula error: case 2, expected real root to match p = ",
          p , ", q = ", q, ", r = ", r, "; ans = ", ans2)
      }
    }
  }
}

func TestSimplifiedQuartic(t *testing.T) {

  //Three cases: zero roots, two roots, or one root. 
  //case 1: The roots are p + i q, p - i q, -p + i r, -p - i r and b >= 0. 
  //Since b == 2 p (q^2 - r^2), this happens when p has the same sign as q^2 - r^2
  for i := 0; i < 10; i++ {
    //Simplified formula
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)
    r := test.RandFloat(-100, 100)

    a := p*p*p*p + p*p*q*q + p*p*r*r + q*q*r*r
    b := 2*p*q*q - 2*p*r*r
    c := -2*p*p + q*q + r*r

    ans := simplifiedQuarticFormula(a, b, c)

    if len(ans) != 0 {
      t.Error("simplified quartic formula error: case 1, expected zero real roots for p = ",
        p , ", q = ", q, ", r = ", r, "; ans = ", ans)
    } 
  }

  //case 2: The roots are p, q, (-p - q)/2 + i r, (-p - q)/2 - i r 
  // b == -(p - q)^2 (p + q)/4 - (p - q) r^2
  for i := 0; i < 10; i++ {
    //Simplified formula
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)
    r := test.RandFloat(-100, 100)

    a := p*p*p*q/4. + p*p*q*q/2. + p*q*q*q/4.+p*q*r*r
    b := -p*p*p/4.+p*p*q/4.+p*q*q/4.-q*q*q/4.-p*r*r-q*r*r
    c := -3.*p*p/4.-p*q/2.-3.*q*q/4.+r*r

    ans := simplifiedQuarticFormula(a, b, c)

    if len(ans) != 2 {
      t.Error("simplified quartic formula error: case 2, expected two real roots for p = ",
        p , ", q = ", q, ", r = ", r, "; ans = ", ans)
    } else {
      test_list := []float64{p, q}
      sort.Float64s(test_list)
      sort.Float64s(ans)

      var close_enough bool = true
      for i := 0; i < len(test_list); i++ {
        close_enough = close_enough && test.CloseEnough(test_list[i], ans[i], e);
      }

      if !close_enough {
        t.Error("simplified quartic formula error: case 2, expected two roots to match p = ",
          p , ", q = ", q, ", r = ", r, "; ans = ", ans)
      }
    }
  }

  //case 3: The roots are p, q, r, -p - q - r
  // b == (p + q)(p + r)(q + r)
  for i := 0; i < 10; i++ {
    //Simplified formula
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)
    r := test.RandFloat(-100, 100)

    a := -p*p*q*r - p*q*q*r - p*q*r*r
    b := p*p*q + p*p*r + p*q*q + p*r*r + q*q*r + q*r*r + 2*p*q*r
    c := -p*p - p*q - p*r - q*q - q*r - r*r

    ans := simplifiedQuarticFormula(a, b, c)

    if len(ans) != 4 {
      t.Error("simplified quartic formula error: case 3, expected four real roots for p = ",
        p , ", q = ", q, ", r = ", r, "; ans = ", ans)
    } else {
      test_list := []float64{p, q, r, -p-q-r}
      sort.Float64s(test_list)
      sort.Float64s(ans)

      var close_enough bool = true
      for i := 0; i < len(test_list); i++ {
        close_enough = close_enough && test.CloseEnough(test_list[i], ans[i], e);
      }

      if !close_enough {
        t.Error("simplified quartic formula error: case 3, expected four roots to match p = ",
          p , ", q = ", q, ", r = ", r, "; ans = ", ans)
      }
    }
  }
}

//Tests only the four root case because it shouldn't make a difference here. 
func TestQuartic(t *testing.T) {

  //The roots are p, q, r, -p - q - r
  for i := 0; i < 10; i++ {
    p := test.RandFloat(-100, 100)
    q := test.RandFloat(-100, 100)
    r := test.RandFloat(-100, 100)
    s := test.RandFloat(-100, 100)

    a := -p*q*r*(p+q+r)
    b := (p + q)*(p + r)*(q + r)
    c := -p*p - p*q - p*r - q*q - q*r - r*r

    af := -(p + q + r - s)*(p + s)*(q + s)*(r + s)
    bf := p*p*q + p*p*r + p*q*q + p*r*r + q*q*r + q*r*r + 2*(-2*s*s*s + p*p*s + p*q*r + p*q*s + p*r*s + q*q*s + q*r*s + r*r*s)
    cf := -p*p - p*q - p*r - q*q - q*r - r*r + 6*s*s
    df := -4*s

    ans := simplifiedQuarticFormula(a, b, c)
    ans_full := QuarticFormula(af, bf, cf, df)

    if len(ans) != len(ans_full) { 
      t.Error("quartic formula error: mismatched number of roots; ans = ", ans, ", ans_full = ", ans_full)
    } else {
      //move the new roots back to being centered at zero; this should put them in the same position
      //as the roots in the simplified formula. 
      for i := 0; i < len(ans); i++ {
        ans_full[i] -= s;
      }

      var close_enough bool = true;
      for i := 0; i < len(ans); i++ {
        close_enough = close_enough && test.CloseEnough(ans[i], ans_full[i], e)
      }

      if !close_enough {
        t.Error("quartic formula error: expected roots to match p = ",
          p , ", q = ", q, ", r = ", r, ", s = ", s,"; ans = ", ans_full)
      }
    }
  }
}
