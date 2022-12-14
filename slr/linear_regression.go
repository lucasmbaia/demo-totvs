package slr

import (
	"errors"
	"math"
	"fmt"
)

type LinearRegression struct {
	VYR		  []float64 `json:"AxisY_Regression"`
	VYP		  []float64 `json:"AxisY_Prediction,omitempty"`
	B0		  float64   `json:"CoeficienteLinear"`
	B1		  float64   `json:"CoeficienteAngular"`
	AnalysisVariance  AnalysisVariance
}

type QuadraticLinearRegression struct {
	VX		  []float64 `json:"-"`
	VY		  []float64 `json:"-"`
	VYR		  []float64 `json:"AxisY_Regression"`
	VYP		  []float64 `json:"AxisY_Prediction,omitempty"`
	B0		  float64   `json:"CoeficienteLinear"`
	B1		  float64   `json:"CoeficienteAngular"`
	B2		  float64   `json:"CoeficienteQuadratico"`
	AnalysisVariance  AnalysisVariance
}

type AnalysisVariance struct {
	PearsonCorrelation	  float64
	GrauDeLiberadeTotal	  float64
	GrauDeLiberadeRegressao	  float64
	GrauDeLiberadeResiduo	  float64
	SomaDeQuadradosTotal	  float64
	SomaDeQuadradosRegressao  float64
	SomaDeQuadradosResiduo	  float64
	QuadradoMedioRegressao	  float64
	QuadradoMedioResiduo	  float64
	FCalc			  float64
	FTab			  float64
	B0Variance		  float64
	B1Variance		  float64
	TCalcB0			  float64
	TCalcB1			  float64
	TTab			  float64
	CoeficienteDeterminacao	  float64
	IsSignificantRegression	  bool
	IsSignificantInterception bool
}


func CalcSimpleLinearRegression(vx, vy, vxp []float64) (lr LinearRegression) {
	var (
		sumY	    float64
		sumPY	    float64
		averageX    float64
		averageY    float64
		b0	    float64
		b1	    float64
		squareSumY  float64
		squareSumPY float64
		y	    float64
		err	    error
	)

	_, sumY, averageX, averageY = CalcSumAndAverageValues(vx, vy)
	b1 = CalcRateOfChange(vx, vy, averageX, averageY)
	b0 = CalcRateOfIntersection(b1, averageX, averageY)

	for _, x := range vx {
		y = b0 + (b1 * x)

		lr.VYR = append(lr.VYR, y)
		sumPY += y
		squareSumPY += y * y
	}

	for _, y := range vy {
		squareSumY += y * y
	}

	for _, x := range vxp {
		lr.VYP = append(lr.VYP,  b0 + (b1 * x))
	}

	lr.B0 = b0
	lr.B1 = b1
	lr.AnalysisVariance.PearsonCorrelation = CalcCoefficientPerson(vx, vy)
	lr.AnalysisVariance.GrauDeLiberadeRegressao = 1
	lr.AnalysisVariance.GrauDeLiberadeTotal = float64(len(vx)) - lr.AnalysisVariance.GrauDeLiberadeRegressao
	lr.AnalysisVariance.GrauDeLiberadeResiduo = lr.AnalysisVariance.GrauDeLiberadeTotal - lr.AnalysisVariance.GrauDeLiberadeRegressao
	lr.AnalysisVariance.SomaDeQuadradosTotal = (squareSumY / 1) - ((sumY * sumY) / float64(len(vy)))
	lr.AnalysisVariance.SomaDeQuadradosRegressao = (squareSumPY / 1) - ((sumPY * sumPY) / float64(len(lr.VYR)))
	lr.AnalysisVariance.SomaDeQuadradosResiduo = lr.AnalysisVariance.SomaDeQuadradosTotal - lr.AnalysisVariance.SomaDeQuadradosRegressao
	lr.AnalysisVariance.QuadradoMedioRegressao = lr.AnalysisVariance.SomaDeQuadradosRegressao / lr.AnalysisVariance.GrauDeLiberadeRegressao
	lr.AnalysisVariance.QuadradoMedioResiduo = lr.AnalysisVariance.SomaDeQuadradosResiduo / lr.AnalysisVariance.GrauDeLiberadeResiduo
	lr.AnalysisVariance.FCalc = lr.AnalysisVariance.QuadradoMedioRegressao / lr.AnalysisVariance.QuadradoMedioResiduo

	if lr.AnalysisVariance.FTab, err = GetFTable(fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeRegressao)), fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeResiduo))); err == nil {
		if math.Abs(lr.AnalysisVariance.FCalc) > lr.AnalysisVariance.FTab {
			lr.AnalysisVariance.IsSignificantRegression = true
		}
	}

	lr.AnalysisVariance.CoeficienteDeterminacao = lr.AnalysisVariance.SomaDeQuadradosRegressao / lr.AnalysisVariance.SomaDeQuadradosTotal

	lr.AnalysisVariance.B0Variance, lr.AnalysisVariance.B1Variance = CalcVarianceB0B1(vx, lr.AnalysisVariance.QuadradoMedioResiduo)
	lr.AnalysisVariance.TCalcB0 = b0 / math.Sqrt(lr.AnalysisVariance.B0Variance)
	lr.AnalysisVariance.TCalcB1 = b1 / math.Sqrt(lr.AnalysisVariance.B1Variance)

	if lr.AnalysisVariance.TTab, err = GetTTable("5%", fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeResiduo))); err == nil {
		if math.Abs(lr.AnalysisVariance.TCalcB0) > lr.AnalysisVariance.TTab {
			lr.AnalysisVariance.IsSignificantInterception = true
		}
	}

	return
}

func CalcQuadraticLinearRegressionWithSeasonality(vx, vy, vxp []float64, p int) (lr QuadraticLinearRegression, err error) {
	var (
		residue	      []float64
		indiceSazonal = make([]float64, p)
		members	      float64
		iterator      int
		xp	      int
		sumY	      float64
		sumPY	      float64
		squareSumY    float64
		squareSumPY   float64
		differ	      float64
	)

	if lr, err = CalcQuadraticLinearRegression(vx, vy, vxp); err != nil {
		return
	}

	for idx, y := range vy {
		residue = append(residue, (y * 100) / lr.VYR[idx])
	}

	for i := 0; i < p; i++ {
		iterator = i
		xp = p
		members = 0

		for {
			indiceSazonal[i] += residue[iterator:iterator+1][0]
			members++

			if iterator + p >= len(residue) {
				break
			}

			iterator += xp
			if len(residue) - iterator < p {
				xp = len(residue) - iterator
			}
		}

		indiceSazonal[i] = indiceSazonal[i] / members
	}

	for _, y := range vy {
		sumY += y
		squareSumY += y * y
	}

	iterator = 0
	for idx, y := range lr.VYR {
		lr.VYR[idx] = (y * indiceSazonal[iterator]) / 100
		sumPY += lr.VYR[idx]
		squareSumPY += lr.VYR[idx] * lr.VYR[idx]

		iterator++

		if iterator == p {
			iterator = 0
		}
	}

	iterator = 0
	for idx, y := range lr.VYP {
		lr.VYP[idx] = (y * indiceSazonal[iterator]) / 100
		iterator++

		if iterator == p {
			iterator = 0
		}
	}

	lr.AnalysisVariance.SomaDeQuadradosResiduo = 0
	for idx, _ := range vy {
		differ = vy[idx] - lr.VYR[idx]
		lr.AnalysisVariance.SomaDeQuadradosResiduo += differ * differ
	}

	lr.AnalysisVariance.SomaDeQuadradosTotal = (squareSumY / 1) - ((sumY * sumY) / float64(len(vy)))
	lr.AnalysisVariance.SomaDeQuadradosRegressao = lr.AnalysisVariance.SomaDeQuadradosTotal - lr.AnalysisVariance.SomaDeQuadradosResiduo
	lr.AnalysisVariance.QuadradoMedioRegressao = lr.AnalysisVariance.SomaDeQuadradosRegressao / lr.AnalysisVariance.GrauDeLiberadeRegressao
	lr.AnalysisVariance.QuadradoMedioResiduo = lr.AnalysisVariance.SomaDeQuadradosResiduo / lr.AnalysisVariance.GrauDeLiberadeResiduo
	lr.AnalysisVariance.FCalc = lr.AnalysisVariance.QuadradoMedioRegressao / lr.AnalysisVariance.QuadradoMedioResiduo

	if lr.AnalysisVariance.FTab, _ = GetFTable(fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeRegressao)), fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeResiduo))); lr.AnalysisVariance.FTab > 0 {
		if math.Abs(lr.AnalysisVariance.FCalc) > lr.AnalysisVariance.FTab {
			lr.AnalysisVariance.IsSignificantRegression = true
		}
	}

	lr.AnalysisVariance.CoeficienteDeterminacao = lr.AnalysisVariance.SomaDeQuadradosRegressao / lr.AnalysisVariance.SomaDeQuadradosTotal

	lr.AnalysisVariance.B0Variance, lr.AnalysisVariance.B1Variance = CalcVarianceB0B1(vx, lr.AnalysisVariance.QuadradoMedioResiduo)
	lr.AnalysisVariance.TCalcB0 = lr.B0 / math.Sqrt(lr.AnalysisVariance.B0Variance)
	lr.AnalysisVariance.TCalcB1 = lr.B1 / math.Sqrt(lr.AnalysisVariance.B1Variance)

	if lr.AnalysisVariance.TTab, _ = GetTTable("5%", fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeResiduo))); lr.AnalysisVariance.TTab > 0 {
		if math.Abs(lr.AnalysisVariance.TCalcB0) > lr.AnalysisVariance.TTab {
			lr.AnalysisVariance.IsSignificantInterception = true
		}
	}

	return
}

func CalcQuadraticLinearRegression(vx, vy, vxp []float64) (lr QuadraticLinearRegression, err error) {
	var (
		values	    []float64
		y	    float64
		elements    = float64(len(vx))
		sumY	    float64
		sumPY	    float64
		squareSumY  float64
		squareSumPY float64
	)

	if values, err = SolveGaussian(vx, vy); err != nil {
		return
	}

	lr.VX = vx
	lr.VY = vy
	lr.B0 = values[0]
	lr.B1 = values[1]
	lr.B2 = values[2]

	for _, x := range vx {
		y = lr.B0 + (lr.B1 * x) + (lr.B2 * (x * x))

		lr.VYR = append(lr.VYR, y)
		sumPY += y
		squareSumPY += y * y
	}

	for _, y := range vy {
		sumY += y
		squareSumY += y * y
	}

	for _, x := range vxp {
		y = lr.B0 + (lr.B1 * x) + (lr.B2 * (x * x))
		lr.VYP = append(lr.VYP, y)
	}

	lr.AnalysisVariance.PearsonCorrelation = CalcCoefficientPerson(vx, vy)
	lr.AnalysisVariance.GrauDeLiberadeRegressao = 2
	lr.AnalysisVariance.GrauDeLiberadeTotal = elements - 1
	lr.AnalysisVariance.GrauDeLiberadeResiduo = lr.AnalysisVariance.GrauDeLiberadeTotal - lr.AnalysisVariance.GrauDeLiberadeRegressao
	lr.AnalysisVariance.SomaDeQuadradosTotal = (squareSumY / 1) - ((sumY * sumY) / float64(len(vy)))
	lr.AnalysisVariance.SomaDeQuadradosRegressao = (squareSumPY / 1) - ((sumPY * sumPY) / float64(len(lr.VYR)))
	lr.AnalysisVariance.SomaDeQuadradosResiduo = lr.AnalysisVariance.SomaDeQuadradosTotal - lr.AnalysisVariance.SomaDeQuadradosRegressao
	lr.AnalysisVariance.QuadradoMedioRegressao = lr.AnalysisVariance.SomaDeQuadradosRegressao / lr.AnalysisVariance.GrauDeLiberadeRegressao
	lr.AnalysisVariance.QuadradoMedioResiduo = lr.AnalysisVariance.SomaDeQuadradosResiduo / lr.AnalysisVariance.GrauDeLiberadeResiduo
	lr.AnalysisVariance.FCalc = lr.AnalysisVariance.QuadradoMedioRegressao / lr.AnalysisVariance.QuadradoMedioResiduo

	if lr.AnalysisVariance.FTab, _ = GetFTable(fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeRegressao)), fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeResiduo))); lr.AnalysisVariance.FTab > 0 {
		if math.Abs(lr.AnalysisVariance.FCalc) > lr.AnalysisVariance.FTab {
			lr.AnalysisVariance.IsSignificantRegression = true
		}
	}

	lr.AnalysisVariance.CoeficienteDeterminacao = lr.AnalysisVariance.SomaDeQuadradosRegressao / lr.AnalysisVariance.SomaDeQuadradosTotal

	lr.AnalysisVariance.B0Variance, lr.AnalysisVariance.B1Variance = CalcVarianceB0B1(vx, lr.AnalysisVariance.QuadradoMedioResiduo)
	lr.AnalysisVariance.TCalcB0 = lr.B0 / math.Sqrt(lr.AnalysisVariance.B0Variance)
	lr.AnalysisVariance.TCalcB1 = lr.B1 / math.Sqrt(lr.AnalysisVariance.B1Variance)

	if lr.AnalysisVariance.TTab, _ = GetTTable("5%", fmt.Sprintf("%d", int(lr.AnalysisVariance.GrauDeLiberadeResiduo))); lr.AnalysisVariance.TTab > 0 {
		if math.Abs(lr.AnalysisVariance.TCalcB0) > lr.AnalysisVariance.TTab {
			lr.AnalysisVariance.IsSignificantInterception = true
		}
	}

	return
}

func SolveGaussian(vx, vy []float64) (values []float64, err error) {
	var (
		sumX	  float64
		sumX2	  float64
		sumX3	  float64
		sumX4	  float64
		sumY	  float64
		sumXY	  float64
		sumX2Y	  float64
		elements  = float64(len(vx))
		relational  [][]float64
	)

	for idx, _ := range vx {
		sumX += vx[idx]
		sumX2 += vx[idx] * vx[idx]
		sumX3 += vx[idx] * vx[idx] * vx[idx]
		sumX4 += vx[idx] * vx[idx] * vx[idx] * vx[idx]
		sumY += vy[idx]
		sumXY += vx[idx] * vy[idx]
		sumX2Y += (vx[idx] * vx[idx]) * vy[idx]
	}

	relational = [][]float64{{elements, sumX, sumX2},{sumX, sumX2, sumX3},{sumX2, sumX3, sumX4}}

	values, err = GaussPartial(relational, []float64{sumY, sumXY, sumX2Y})

	return
}

func GaussPartial(a0 [][]float64, b0 []float64) ([]float64, error) {
	// make augmented matrix
	m := len(b0)
	a := make([][]float64, m)
	for i, ai := range a0 {
	        row := make([]float64, m+1)
	        copy(row, ai)
	        row[m] = b0[i]
	        a[i] = row
	}
	// WP algorithm from Gaussian elimination page
	// produces row-eschelon form
	for k := range a {
	        // Find pivot for column k:
	        iMax := k
	        max := math.Abs(a[k][k])
	        for i := k + 1; i < m; i++ {
	                if abs := math.Abs(a[i][k]); abs > max {
	                	iMax = i
	                	max = abs
	                }
	        }
	        if a[iMax][k] == 0 {
	                return nil, errors.New("singular")
	        }
	        // swap rows(k, i_max)
	        a[k], a[iMax] = a[iMax], a[k]
	        // Do for all rows below pivot:
	        for i := k + 1; i < m; i++ {
	                // Do for all remaining elements in current row:
	                for j := k + 1; j <= m; j++ {
	                	a[i][j] -= a[k][j] * (a[i][k] / a[k][k])
	                }
	                // Fill lower triangular matrix with zeros:
	                a[i][k] = 0
	        }
	}
	// end of WP algorithm.
	// now back substitute to get result.
	x := make([]float64, m)
	for i := m - 1; i >= 0; i-- {
	        x[i] = a[i][m]
	        for j := i + 1; j < m; j++ {
	                x[i] -= a[i][j] * x[j]
	        }
	        x[i] /= a[i][i]
	}
	return x, nil
}

func CalcVarianceB0B1(vx []float64, qmr float64) (varianceb0, varianceb1 float64) {
	var (
		sumX	    float64
		squareSumX  float64
		squareX	    float64
		elements    = float64(len(vx))
	)

	for _, x := range vx {
		sumX += x
		squareSumX += x * x
	}

	squareX = (sumX / elements) * (sumX / elements)
	varianceb0 = ((1 / elements) + (squareX / (squareSumX - ((sumX * sumX) / elements)))) * qmr
	varianceb1 = (1 / (squareSumX  - ((sumX * sumX) / elements))) * qmr

	return
}

func CalcRateOfChange(vx, vy []float64, ax, ay float64) float64 { //b1 = coeficiente angular = inclina????o da reta
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

func CalcRateOfIntersection(brate, ax, ay float64) float64 { //b0 = intercepto = coeficiente linear = tempo independente
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
