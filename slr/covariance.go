package slr

import (
	"errors"
)

func CalcCov(values []float64, h int) (cov float64, err error) {
	var (
		n	    float64
		averageDev1 float64
		averageDev2 float64
		deviation1  []float64
		deviation2  []float64
		sumProduct  float64
	)

	if len(values) < h {
		err = errors.New("Range deve ser maior que o h informado")
		return
	}

	n = float64(len(values[h:]))

	for _, dev1 := range values[h:] {
		averageDev1 += dev1
	}

	for _, dev2 := range values[:len(values)-h] {
		averageDev2 += dev2
	}

	averageDev1 = averageDev1 / n
	averageDev2 = averageDev2 / n

	for _, dev1 := range values[h:] {
		deviation1 = append(deviation1, dev1 - averageDev1)
	}

	for _, dev2 := range values[:len(values)-h] {
		deviation2 = append(deviation2, dev2 - averageDev2)
	}

	for idx, _ := range deviation1 {
		sumProduct += deviation1[idx] * deviation2[idx]
	}

	cov = sumProduct / n

	return
}
