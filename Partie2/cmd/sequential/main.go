package main

import (
	"csi2520/partie2/internal/algo"
	csv "csi2520/partie2/internal/csv"
)

func main() {
	residentMap, err := csv.ReadResidentsCSV("../ResidentsPrograms/residentSmall.csv")
	if err != nil {
		panic(err)
	}

	programMap, err := csv.ReadProgramsCSV("../ResidentsPrograms/programSmall.csv")
	if err != nil {
		panic(err)
	}

	algo.McVitieWilsonSequential(residentMap, programMap)
}
