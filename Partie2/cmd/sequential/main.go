package main

import (
	"csi2520/partie2/internal/algo"
	csv "csi2520/partie2/internal/csv"
	"fmt"
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

	// ================ OUTPUT RESULTS ================
	var unmatchedResidentCount int
	var numPositionsAvailable int

	// Print CSV header
	fmt.Println("lastname,firstname,residentID,programID,name")

	// Print match results
	for _, resident := range residentMap {
		// Get the name of the matched program (defaults to NOT_MATCHED)
		var programName string
		if resident.MatchedProgram != "XXX" {
			programName = programMap[resident.MatchedProgram].Name
		} else {
			programName = "NOT_MATCHED"
			unmatchedResidentCount++
		}

		fmt.Printf("%s,%s,%d,%s,%s\n",
			resident.Lastname,
			resident.Firstname,
			resident.ResidentID,
			resident.MatchedProgram,
			programName)
	}

	// Get unfilled position count from programs
	for _, program := range programMap {
		unfilled := program.NPositions - program.SelectedResidents.Len()
		numPositionsAvailable += unfilled
	}

	// Print summary statistics
	fmt.Println()
	fmt.Println("Number of unmatched residents:", unmatchedResidentCount)
	fmt.Println("Number of positions available:", numPositionsAvailable)
}
