package algo

import csv "csi2520/partie2/internal/csv"

func McVitieWilsonSequential(residents map[int]*csv.Resident, programs map[string]*csv.Program) {
	for r := range residents {
		offer(r, residents, programs)
	}
}

func offer(rid int, residents map[int]*csv.Resident, programs map[string]*csv.Program) {
	// TODO: find next available program in residentROL. What does available mean? No matches? Quota not reached yet? Not evaluated yet?

}
func evaluate(rid int, pid string, residents map[int]*csv.Resident, programs map[string]*csv.Program) {

}
