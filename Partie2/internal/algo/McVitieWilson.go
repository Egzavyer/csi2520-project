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
		resident.MatchedProgram = ""
		return
	}

	// Increment the next program index for the next call
	resident.NextProgIdx++

	progId := resident.Rol[idx]
	progToTry := programs[progId]

	evaluate(resident, progToTry, residents, programs)
}

func evaluate(resident *csv.Resident, program *csv.Program, residents map[int]*csv.Resident, programs map[string]*csv.Program) {

}
