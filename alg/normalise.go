package alg

import (
	"math"
	"time"
)

func NormaliseTime(t1 *time.Time,ctime *time.Time) float64 { // normalises the time (gets the delta between the most recent and oldest post)
	return ctime.Sub(*t1).Minutes()
}

func NormalisePriority(p int) float64 {
	xcoeff := 0.2 // x coefficient for the tanh function
	tanhCoeff := 0.5 // coefficient for tanh function
	c := 0.5 // + c value to translate the values

	tanh := tanh(xcoeff * float64(p)) // gets tanh value

	//           (  e^dx - x^-dx )
	// y = 0.5 x ( ------------- ) + 0.5   {0 < x}
	//           (  e^dx + x^-dx )

	// this is what the formula looks when written in cartesian form

	return tanhCoeff*tanh + c // return processed value
}

func tanh(x float64) float64 {
	// performs hyperbolic tanh function
	// e^x - x^-x
	// ----------
	// e^x + x^-x

	// this is derived from sinh(x)/cosh(x)

	//         e^x - x^-x
	// sinh = ------------
	//              2


	//         e^x + x^-x
	// cosh = ------------
	//              2

	e := math.Exp(x) // gets e^x
	en := 1/math.Exp(x)  // gets e^-x

	return (e-en)/(e+en)
}