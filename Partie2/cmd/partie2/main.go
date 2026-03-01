package main

import (
	"csi2520/partie2/internal/algo"
	csv "csi2520/partie2/internal/csv"
	"fmt"
	"os"
	"slices"
	"time"
)

func main() {
	// Parse cmdline args
	if len(os.Args) != 4 {
		fmt.Println("Usage: projet <sequential | concurrent> <residentsFile> <programsFile>")
		return
	}
	execution := os.Args[1]
	residentFile := os.Args[2]
	programFile := os.Args[3]

	residentMap, err := csv.ReadResidentsCSV(residentFile)
	if err != nil {
		panic(err)
	}

	programMap, err := csv.ReadProgramsCSV(programFile)
	if err != nil {
		panic(err)
	}

	var start time.Time
	var end time.Time

	switch execution {
	case "sequential":
		start = time.Now()

		algo.McVitieWilsonSequential(residentMap, programMap)

		end = time.Now()
	case "concurrent":
		start = time.Now()

		algo.McVitieWilsonConcurrent(residentMap, programMap)

		end = time.Now()
	default:
		fmt.Printf("Invalid program execution type: %s", execution)
		fmt.Println("Usage: projet <sequential | concurrent> <residentsFile> <programsFile>")
		return
	}

	// ================ OUTPUT RESULTS ================
	var unmatchedResidentCount int
	var numPositionsAvailable int

	// Print CSV header
	fmt.Println("lastname,firstname,residentID,programID,name")

	// Sort residents by lastname for deterministic output
	sortedResidents := make([]*csv.Resident, 0, len(residentMap))
	for _, r := range residentMap {
		sortedResidents = append(sortedResidents, r)
	}
	slices.SortFunc(sortedResidents, func(a, b *csv.Resident) int {
		if a.Lastname < b.Lastname {
			return -1
		} else if a.Lastname > b.Lastname {
			return 1
		}
		return 0
	})

	// Print match results
	for _, resident := range sortedResidents {
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

	fmt.Printf("\nExecution time: %s\n", end.Sub(start))
}
