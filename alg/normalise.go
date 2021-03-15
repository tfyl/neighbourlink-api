package alg

import (
	"math"
	"time"
)

func NormaliseTime(t1 *time.Time,ctime *time.Time) float64 {
	return ctime.Sub(*t1).Minutes()
}

func NormalisePriority(p int) float64 {
	xcoeff := 0.2
	tanhCoeff := 0.5
	c := 0.5

	tanh := tanh(xcoeff * float64(p))

	return tanhCoeff*tanh + c
}

func tanh(x float64) float64 {
	e := math.Exp(x)
	en := 1/math.Exp(x)

	return (e-en)/(e+en)
}