package algo

import (
	"math"
	"slices"

	csv "csi2520/partie2/internal/csv"
)

// Performs the McVitie-Wilson algorithm sequentially to match residents to programs.
func McVitieWilsonSequential(residents map[int]*csv.Resident, programs map[string]*csv.Program) {
	for r := range residents {
		offer(r, residents, programs)
	}
}

// Attempts to offer the resident with id rid a spot in the next available program on their list.
// Will call [evaluate] to determine if they can be matched to said program.
// If they are matched, this function updates the corresponding states and returns.
// If the resident is not matched, the call to [evaluate] will re-call [offer] and auto-attempt
// the next program in the resident's list.
// When this method returns, the resident has been tried against every program in their list
// and either has been matched, or the state will reflect they could not be matched.
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

// Evaluates if a resident can be matched to a given program, and if so, updates the state to reflect that.
// If the resident cannot be matched, this method will re-call [offer] to attempt the next match.
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

// Checks if a program p prefers a resident r over their least-preferred current match.
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
