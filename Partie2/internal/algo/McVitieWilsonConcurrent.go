package algo

import (
	csv "csi2520/partie2/internal/csv"
	"slices"
	"sync"
)

// Performs the McVitie-Wilson algorithm concurrently to match residents to programs.
// A goroutine is spawned for every resident to concurrently optimize resident matching.
func McVitieWilsonConcurrent(
	residents map[int]*csv.Resident,
	programs map[string]*csv.Program,
) {

	var wg sync.WaitGroup

	for id := range residents {
		wg.Add(1)
		go offerConcurrent(id, residents, programs, &wg)
	}

	wg.Wait()
}

// Attempts to offer the resident with id rid a spot in the next available program on their list.
// Will call [evaluateConcurrent] to determine if they can be matched to said program.
// If they are matched, this function updates the corresponding states and returns.
// If the resident is not matched, the call to [evaluateConcurrent] will re-call [offerConcurrent] and auto-attempt
// the next program in the resident's list.
// When this method returns, the resident has been tried against every program in their list
// and either has been matched, or the state will reflect they could not be matched.
func offerConcurrent(
	rid int,
	residents map[int]*csv.Resident,
	programs map[string]*csv.Program,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	resident := residents[rid]

	resident.Lock()

	// Check that the resident still has programs to evaluate
	idx := resident.NextProgIdx
	if idx >= len(resident.Rol) {
		// We have exhausted the list of programs the resident wants
		resident.MatchedProgram = "XXX"
		resident.Unlock()
		return
	}

	// Increment the next program index for the next call
	resident.NextProgIdx++
	progID := resident.Rol[idx]
	resident.Unlock()

	program := programs[progID]

	wg.Add(1)
	go evaluateConcurrent(resident, program, residents, programs, wg)
}

// Evaluates if a resident can be matched to a given program, and if so, updates the state to reflect that.
// If the resident cannot be matched, this method will re-call [offerConcurrent] to attempt the next match.
func evaluateConcurrent(
	resident *csv.Resident,
	program *csv.Program,
	residents map[int]*csv.Resident,
	programs map[string]*csv.Program,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	program.Lock()

	if !slices.Contains(program.Rol, resident.ResidentID) {
		// Program will not take resident so go back and try to match them to their next program
		program.Unlock()

		wg.Add(1)
		go offerConcurrent(resident.ResidentID, residents, programs, wg)
		return
	}

	if program.SelectedResidents.Len() < program.NPositions {
		// Get the rank for the resident
		rank := slices.Index(program.Rol, resident.ResidentID)

		// Add the resident to the program's selected residents
		// Priority is the resident's rank (high indicates lower priority)
		program.SelectedResidents.Push(rank, resident)

		program.Unlock()

		// Set the resident's matched program
		resident.Lock()
		resident.MatchedProgram = program.ProgramID
		resident.Unlock()

		return
	}

	// Program full → check if it prefers the new resident over its least preferred current match
	rRank := slices.Index(program.Rol, resident.ResidentID)
	rPrimeRank, rPrime := program.SelectedResidents.Peek()

	if rRank < rPrimeRank {
		// Remove the least preferred resident and replace them with the new one
		program.SelectedResidents.Pop()
		program.SelectedResidents.Push(rRank, resident)
		program.Unlock()

		// Set the resident's matched program
		resident.Lock()
		resident.MatchedProgram = program.ProgramID
		resident.Unlock()

		// Re-match the displaced resident
		wg.Add(1)
		go offerConcurrent(rPrime.ResidentID, residents, programs, wg)
		return
	}

	program.Unlock()

	// Go back and try to match the resident to their next program
	wg.Add(1)
	go offerConcurrent(resident.ResidentID, residents, programs, wg)
}
