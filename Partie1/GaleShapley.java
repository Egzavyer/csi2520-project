import java.io.*;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;

public class GaleShapley {
    /** The residents part of the matching process mapped by their IDs. */
    private HashMap<Integer, Resident> residents = new HashMap<>();
    /** The programs part of the matching process mapped by their IDs. */
    private HashMap<String, Program> programs = new HashMap<>();

    public GaleShapley(String residentsFilename, String programsFilename) throws IOException, NumberFormatException {
        loadResidents(residentsFilename);
        loadPrograms(programsFilename);
    }

    /**
     * The results of matching residents to programs.
     * Contains the matched residents and their programs, as well as the unmatched residents.
     * 
     * For convienience, also contains the collections of residents and programs used in the matching.
     * This has a negligible memory overhead as the collection references the existing memory.
     * Both residents and programs are immutable after loading, returned primarily for later statistics.
     */
    public class MatchResult {
        /** The map of matched residents to their assigned programs. */
        public HashMap<Resident, Program> matches;
        /** The list of unmatched residents. */
        public List<Resident> unmatchedResidents;

        /** The immutable collection of all residents part of the matching process. */
        public Collection<Resident> residents;
        /** The immutable collection of all programs part of the matching process. */
        public Collection<Program> programs;

        public MatchResult(HashMap<Resident, Program> matches, List<Resident> unmatched) {
            this.matches = matches;
            this.unmatchedResidents = unmatched;

            // Reference the HashMap values for low-overhead reference
            this.residents = GaleShapley.this.residents.values();
            this.programs = GaleShapley.this.programs.values();
        }
    }

    /**
     * Matches residents to programs using the stable Gale-Shapley algorithm.
     * @return A HashMap mapping each resident to their assigned program.
     */
    public MatchResult matchResidentsToPrograms() {
        // Residents that are not yet matched
        var residentsToMatch = new ArrayList<>(this.residents.values());
        // Residents that could not be matched
        List<Resident> unmatchedResidents = new ArrayList<>();

        while (!residentsToMatch.isEmpty()) {
            Resident resident = residentsToMatch.removeFirst();
            
            // Go through the resident's ROL which is already ordered by preference
            // Match if possible and break, otherwise keep going through the ROL
            for (String programID : resident.getRol()) {
                Program program = programs.get(programID);

                // Program does not accept this resident
                if (!program.member(resident.getId())) {
                    continue;
                }

                // Program has available quota
                else if (program.getMatchedResidents().size() < program.getQuota()) {
                    matchResidentToProgram(resident, program);
                    break;
                }

                // Program prefers this resident over another currently matched resident
                else if (program.prefers(resident.getId())) {
                    // Remove the least preferred resident
                    Resident removedResident = program.removeLeastPreferred();
                    removedResident.setMatchedProgram(null);
                    removedResident.setMatchedRank(-1);

                    // Removed resident becomes available again
                    residentsToMatch.add(removedResident);

                    // Match the new resident
                    matchResidentToProgram(resident, program);                
                    break;
                }
            }

            // Resident could not be matched to any program in their ROL
            if (resident.getMatchedProgram() == null) {
                unmatchedResidents.add(resident);
            }
        }

        // Build the results
        HashMap<Resident, Program> matches = new HashMap<>();
        residents.values().forEach(resident -> {
            if (resident.getMatchedProgram() != null)
                matches.put(resident, resident.getMatchedProgram());
        });

        return new MatchResult(matches, unmatchedResidents);
    }

    /**
     * Updates the states of both the resident and the program to reflect a match.
     * Sets the resident's matched program and rank, and adds the resident to the program's list of matched residents.
     * @param resident The resident being matched.
     * @param program The program the resident is being matched to.
     */
    private void matchResidentToProgram(Resident resident, Program program) {
        program.addResident(resident);
        resident.setMatchedProgram(program);
        resident.setMatchedRank(program.rank(resident.getId()));
    }

    /**
     * Loads residents from a CSV file with headers into the residents map.
     * @param filePath The path to the CSV file containing resident data.
     * @throws IOException When there is an issue reading the file.
     * @throws NumberFormatException When there is an issue parsing the resident's ID.
     */
    private void loadResidents(String filePath) throws IOException, NumberFormatException {
        try (BufferedReader br = new BufferedReader(new FileReader(filePath))) {
            br.readLine(); // Skip headers

            // Deserialize each line as a Resident and add them to the map
            String line;
            while ((line = br.readLine()) != null && line.length() > 0) {
                var split = line.split(",", 4);
                if (split.length != 4)
                    throw new IOException("Invalid line format: " + line);

                int residentID = Integer.parseInt(split[0]);
                String firstName = split[1];
                String lastName  = split[2];
                String[] rol = split[3].substring(2, split[3].length() - 2) // Ignore "[ and ]"
                                       .split(",");

                Resident resident = new Resident(residentID, firstName, lastName, rol);
                residents.put(residentID, resident);
            }
        }
    }

    /**
     * Loads programs from a CSV file with headers into the programs map.
     * @param filePath The path to the CSV file containing program data.
     * @throws IOException When there is an issue reading the file.
     * @throws NumberFormatException When there is an issue parsing the program's quota or ID's of residents in the ROL.
     */
    private void loadPrograms(String filePath) throws IOException, NumberFormatException {
        try (BufferedReader br = new BufferedReader(new FileReader(filePath))) {
            br.readLine(); // Skip headers

            // Deserialize each line as a Program and add them to the map
            String line;
            while ((line = br.readLine()) != null && line.length() > 0) {
                var split = line.split(",", 4);
                if (split.length != 4)
                    throw new IOException("Invalid line format: " + line);

                String programID = split[0];
                String name = split[1];
                int quota = Integer.parseInt(split[2]);
                int[] rol = Arrays.stream(split[3].substring(2, split[3].length() - 2) // Ignore "[ and ]"
                                                  .split(","))
                                  .mapToInt(Integer::parseInt)
                                  .toArray();

                Program program = new Program(programID, name, quota, rol);
                programs.put(programID, program);
            }
        }
    }
}
