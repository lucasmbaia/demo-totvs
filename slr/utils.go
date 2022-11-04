package slr

import (
	"strings"
	"strconv"
	"errors"
	"bufio"
	"fmt"
	"os"
)

func GetTTable(importance, grau string) (value float64, err error){
	var (
		file	    *os.File
		scanner	    *bufio.Scanner
		isFirstLine = true
		line	    []string
		column	    int
	)

	if file, err = os.Open("/opt/files/ttable.tsv"); err != nil {
		return
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line = strings.Split(scanner.Text(), "\t")

		if isFirstLine {
			for idx, str := range line {
				if str == importance {
					column = idx
				}
			}

			isFirstLine = false
			if column == 0 {
				err = errors.New(fmt.Sprintf("Importancia %s informada não foi encontrada", importance))
				return
			}
		} else {
			if line[0] == grau {
				if value, err = strconv.ParseFloat(line[column], 64); err != nil {
					return
				}

				break
			}
		}
	}

	if value == 0 {
		err = errors.New(fmt.Sprintf("Valor correspondente a importancia %s e grau %s informados não foi encontrado", importance, grau))
	}

	return
}

func GetFTable(grauRegressao, grauResiduo string) (value float64, err error){
	var (
		file	    *os.File
		scanner	    *bufio.Scanner
		isFirstLine = true
		line	    []string
		column	    int
	)

	if file, err = os.Open("/opt/files/ftable.tsv"); err != nil {
		return
	}
	defer file.Close()

	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line = strings.Split(scanner.Text(), "\t")

		if isFirstLine {
			for idx, str := range line {
				if str == grauRegressao {
					column = idx
				}
			}

			isFirstLine = false
			if column == 0 {
				err = errors.New(fmt.Sprintf("Grau de regressão %s informada não foi encontrado", grauRegressao))
				return
			}
		} else {
			if line[0] == grauResiduo {
				if value, err = strconv.ParseFloat(line[column], 64); err != nil {
					return
				}

				break
			}
		}
	}

	if value == 0 {
		err = errors.New(fmt.Sprintf("Valor correspondente ao grau de regressão %s e ao grau de resíduo %s informados não foi encontrado", grauRegressao, grauResiduo))
	}

	return
}
