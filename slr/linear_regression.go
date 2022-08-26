package slr

import (
	"math"
)

func CalcRateOfChange(vx, vy []float64, ax, ay float64) float64 { //b1
	var (
		sxy float64
		sxx float64
	)


	for i := 0; i < len(vx); i++ {
		sxy += (vx[i] -(ax)) * (vy[i] -(ay))
		sxx += (vx[i] -(ax)) * (vx[i] -(ax))
	}

	return sxy / sxx
}

func CalcRateOfIntersection(brate, ax, ay float64) float64 { //b0
	return ay - (brate * ax)
}

func CalcSumAndAverageValues(vx, vy []float64) (sx, sy, ax, ay float64) {
	for _, v := range vx {
		sx += v
	}

	for _, v := range vy {
		sy += v
	}

	ax = sx / float64(len(vx))
	ay = sy / float64(len(vy))

	return
}

func CalcCoefficientPerson(vx, vy []float64) float64 {
	var (
		sx    float64
		sy    float64
		sxx   float64
		syy   float64
		sxy   float64
		f1    float64
		f2    float64
		size  float64
	)

	sx, sy, _, _ = CalcSumAndAverageValues(vx, vy)
	size = float64(len(vx))

	for i := 0; i < len(vx); i++ {
		sxx += vx[i] * vx[i]
		syy += vy[i] * vy[i]
		sxy += vx[i] * vy[i]
	}

	f1 = (size * sxy) - (sx * sy)
	f2 = math.Sqrt((size * sxx - (sx * sx)) * (size * syy - (sy * sy)))

	return f1 / f2
}
