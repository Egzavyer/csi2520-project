/**
 * Represents a resident in the stable matching problem.
 */
public class Resident {
    private int id;
    private String firstName;
    private String lastName;

    /** List of programs this resident prefers sorted from most to least preferred. */
    private String[] rol;

    /** The program this resident is matched to. */
    private Program matchedProgram;
    
    /** The rank of this resident in their matched program. */
    private int matchedRank = -1;

    public Resident(int id, String firstName, String lastName, String[] rol) {
        this.id = id;
        this.firstName = firstName;
        this.lastName = lastName;
        this.rol = rol;
    }
    
    public int getId() { return id; }
    public String getFirstName() { return firstName; }
    public String getLastName() { return lastName; }
    public String[] getRol() { return rol; }

    public Program getMatchedProgram() { return matchedProgram; }
    public void setMatchedProgram(Program matchedProgram) { this.matchedProgram = matchedProgram; }
    
    public int getMatchedRank() { return matchedRank; }
    public void setMatchedRank(int rank) { this.matchedRank = rank; }
    
    //#region Overrides for HashMap keying
    
    @Override
    public int hashCode() {
        // Resident ID is unique thus sufficient for hash code
        return id;
    }

    @Override
    public boolean equals(Object obj) {
        return obj instanceof Resident other && this.id == other.id;
    }

    //#endregion
}
