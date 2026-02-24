package algo

import (
	"math"
	"slices"

	csv "csi2520/partie2/internal/csv"
)

func McVitieWilsonSequential(residents map[int]*csv.Resident, programs map[string]*csv.Program) {
	for r := range residents {
		offer(r, residents, programs)
	}
}

func offer(rid int, residents map[int]*csv.Resident, programs map[string]*csv.Program) {
	resident := residents[rid]

	// Check that the resident still has programs to evaluate
	idx := resident.NextProgIdx
	if idx >= len(resident.Rol) {
		// We have exhausted the list of programs the resident wants
		resident.MatchedProgram = "XXX"
		return
	}

	// Increment the next program index for the next call
	resident.NextProgIdx++

	progId := resident.Rol[idx]
	progToTry := programs[progId]

	evaluate(resident, progToTry, residents, programs)
}

func evaluate(resident *csv.Resident, program *csv.Program, residents map[int]*csv.Resident, programs map[string]*csv.Program) {
	if !slices.Contains(program.Rol, resident.ResidentID) {
		// Program will not take resident so go back and try to match them to their next program
		offer(resident.ResidentID, residents, programs)
	} else if program.NPositions != program.SelectedResidents.Len() {
		// Get the rank for the resident (or max int size if not in Rol as high rank is lower priority)
		rank := slices.Index(program.Rol, resident.ResidentID)
		if rank == -1 {
			rank = math.MaxInt
		}

		// Add the resident to the program's selected residents
		// Priority is the resident's rank (high indicates lower priority)
		program.SelectedResidents.Push(rank, resident)
		// Set the resident's matched program
		resident.MatchedProgram = program.ProgramID
	} else if rPrime, ok := prefers(program, resident); ok {
		// Get the rank for the resident
		rRank := slices.Index(program.Rol, resident.ResidentID) // Will not be -1 as prefers returns false if not present and this won't run

		// Remove the least preferred resident and replace them with the new one
		program.SelectedResidents.Pop()
		program.SelectedResidents.Push(rRank, resident)
		// Set the resident's matched program
		resident.MatchedProgram = program.ProgramID

		// Re-match the displaced resident
		offer(rPrime.ResidentID, residents, programs)
	} else {
		// Go back and try to match the resident to their next program
		offer(resident.ResidentID, residents, programs)
	}
}

func prefers(p *csv.Program, r *csv.Resident) (*csv.Resident, bool) {
	rRank := slices.Index(p.Rol, r.ResidentID)
	if rRank == -1 {
		return nil, false
	}

	rPrimeRank, rPrime := p.SelectedResidents.Peek()
	if rRank < rPrimeRank {
		return rPrime, true
	}

	return nil, false
}
