package algo

import (
	csv "csi2520/partie2/internal/csv"
	"slices"
	"sync"
)

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

func offerConcurrent(
	rid int,
	residents map[int]*csv.Resident,
	programs map[string]*csv.Program,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	resident := residents[rid]

	resident.Mu.Lock()

	// Check that the resident still has programs to evaluate
	idx := resident.NextProgIdx
	if idx >= len(resident.Rol) {
		// We have exhausted the list of programs the resident wants
		resident.MatchedProgram = "XXX"
		resident.Mu.Unlock()
		return
	}

	// Increment the next program index for the next call
	resident.NextProgIdx++
	progID := resident.Rol[idx]
	resident.Mu.Unlock()

	program := programs[progID]

	wg.Add(1)
	go evaluateConcurrent(resident, program, residents, programs, wg)
}

func evaluateConcurrent(
	resident *csv.Resident,
	program *csv.Program,
	residents map[int]*csv.Resident,
	programs map[string]*csv.Program,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	program.Mu.Lock()

	if !slices.Contains(program.Rol, resident.ResidentID) {
		// Program will not take resident so go back and try to match them to their next program
		program.Mu.Unlock()

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

		program.Mu.Unlock()

		// Set the resident's matched program
		resident.Mu.Lock()
		resident.MatchedProgram = program.ProgramID
		resident.Mu.Unlock()

		return
	}

	// Program full → check if it prefers the new resident over its least preferred current match
	rRank := slices.Index(program.Rol, resident.ResidentID)
	rPrimeRank, rPrime := program.SelectedResidents.Peek()

	if rRank < rPrimeRank {
		// Remove the least preferred resident and replace them with the new one
		program.SelectedResidents.Pop()
		program.SelectedResidents.Push(rRank, resident)
		program.Mu.Unlock()

		// Set the resident's matched program
		resident.Mu.Lock()
		resident.MatchedProgram = program.ProgramID
		resident.Mu.Unlock()

		// Re-match the displaced resident
		wg.Add(1)
		go offerConcurrent(rPrime.ResidentID, residents, programs, wg)
		return
	}

	program.Mu.Unlock()

	// Go back and try to match the resident to their next program
	wg.Add(1)
	go offerConcurrent(resident.ResidentID, residents, programs, wg)
}
