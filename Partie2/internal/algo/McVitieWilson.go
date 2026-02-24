package algo

import csv "csi2520/partie2/internal/csv"

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

	if !in(resident.ResidentID, program.Rol) {
		offer(resident.ResidentID, residents, programs)
	} else if program.NPositions != len(program.SelectedResidents) {
		resident.MatchedProgram = program.ProgramID
		program.SelectedResidents = append(program.SelectedResidents, resident)
	} else if rPrime, ok := prefers(program, resident); ok {
		program.SelectedResidents = program.SelectedResidents[len(program.SelectedResidents)-1:] // Remove r' from the selected residents
		program.SelectedResidents = append(program.SelectedResidents, resident)
		resident.MatchedProgram = program.ProgramID
		offer(rPrime.ResidentID, residents, programs)
	} else {
		offer(resident.ResidentID, residents, programs)
	}
}

func in(val int, list []int) bool {
	for v := range list {
		if val == v {
			return true
		}
	}
	return false
}

func prefers(p *csv.Program, r *csv.Resident) (*csv.Resident, bool) {
	rPrime := p.SelectedResidents[len(p.SelectedResidents)-1]
	var (
		rPrimeRank int
		rRank      int
	)
	for currentRes, idx := range p.Rol {
		if currentRes == rPrime.ResidentID {
			rPrimeRank = idx
		}
		if currentRes == r.ResidentID {
			rRank = idx
		}
	}

	if rPrimeRank < rRank {
		return rPrime, true
	} else {
		return nil, false
	}
}
