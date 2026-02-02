import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * Represents a program in the stable matching problem.
 */
public class Program {
    private String id;
    private String name;
    private int quota;

    /** List of the IDs of the residents this program prefers sorted from most to least preferred. */
    private int[] rol;
    /** List of residents matched to this program. */
    private List<Resident> matchedResidents = new ArrayList<>();

    public Program(String id, String name, int quota, int[] rol) {
        this.id = id;
        this.name = name;
        this.quota = quota;
        this.rol = rol;
    }

    public String getId() { return id; }
    public String getName() { return name; }
    public int getQuota() { return quota; }
    public int[] getRol() { return rol; }
    public List<Resident> getMatchedResidents() { return matchedResidents; }
    
    /**
     * Checks if a resident with the given ID is in the program's ROL.
     * @param residentID The ID of the resident to check.
     * @return True if the resident is in the ROL, else false.
     */
    public boolean member(int residentID) {
        return Arrays.stream(rol).anyMatch(id -> id == residentID);
    }
    
    /**
     * Gets the preference rank of a resident. A lower number indicates a higher preference.
     * @implSpec If a resident is not present in the ROL, returns -1.
     * @param residentID The ID of the resident to get the rank for.
     * @return The rank of the resident in the ROL or -1.
     */
    public int rank(int residentID) {
        for (int i = 0; i < rol.length; i++) {
            if (rol[i] == residentID) {
                return i;
            }
        }

        return -1;
    }

    /**
     * Determines if this program prefers the given resident over any of its currently matched residents.
     * Does not modify the state of the program.
     * @param residentID The ID of the resident to check preference for.
     * @return True if the program prefers the resident over at least one currently matched resident, else false.
     */
    public boolean prefers(int residentID) {
        // We'll use anyMatch here for early exit as soon as we find a resident that is less preferred.
        // Finding the rank of the least preferred matched resident is not necessary as this method does not modify state,
        // thus a call to removeLeastPreferred() is necessary by the caller later on.
        return matchedResidents.stream().anyMatch(r -> rank(r.getId()) > rank(residentID));
    }
    
    /**
     * Removes the least preferred resident that is currently matched to this program returning the removed resident.
     * @apiNote The requested method name by the assignment is "leastPreferred", however this name is somewhat misleading as the implementation modifies state.
     * @return The least preferred resident or null if no residents are matched.
     */
    public Resident removeLeastPreferred() {
        // NOTE: This could be optimized by keeping matchedResidents sorted by ROL order or using a priority queue or hash map.
        //       However, this is sufficient for the current implementation, thus no optimizations adding such complexity have been performed.
        //       Additionally, we could use r1.getMatchedRank() instead of calling rank(r1.getId()), but that couples state too tightly.

        return matchedResidents.stream()
            .max((r1, r2) -> Integer.compare(rank(r1.getId()), rank(r2.getId())))
            .map(resident -> {
                matchedResidents.remove(resident); // Remove the least preferred resident
                return resident;                   // Return the removed resident
            })
            .orElse(null);
    }
    
    /**
     * Adds a resident to the list of matched residents.
     * @param resident The resident to add.
     */
    public void addResident(Resident resident) {
        matchedResidents.add(resident);
    }
}
