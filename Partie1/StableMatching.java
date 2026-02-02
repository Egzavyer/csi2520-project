import java.io.IOException;

class StableMatching {
    public static void main(String[] args) {
        if (args.length != 2) {
            System.err.println("Usage: java StableMatching <residentsFile> <programsFile>");
            return;
        }

        String residentsFile = args[0];
        String programsFile  = args[1];

        try {
            // Match residents to programs using the Gale-Shapley algorithm
            var matchResults = new GaleShapley(residentsFile, programsFile)
                .matchResidentsToPrograms();
            
            // Print CSV header
            System.out.println("lastname,firstname,residentID,programID,name");

            // Print matched residents
            for (var entry : matchResults.matches.entrySet()) {
                Resident resident = entry.getKey();
                Program program = entry.getValue();

                StringBuilder sb = new StringBuilder()
                    .append(resident.getLastName())
                    .append(",")
                    .append(resident.getFirstName())
                    .append(",")
                    .append(resident.getId())
                    .append(",")
                    .append(program.getId())
                    .append(",")
                    .append(program.getName());

                System.out.println(sb.toString());
            }

            // Print unmatched residents
            for (Resident resident : matchResults.unmatchedResidents) {
                StringBuilder sb = new StringBuilder()
                    .append(resident.getLastName())
                    .append(",")
                    .append(resident.getFirstName())
                    .append(",")
                    .append(resident.getId())
                    .append(",")
                    .append("XXX")           // No program ID
                    .append(",")
                    .append("NOT_MATCHED"); // No program name

                System.out.println(sb.toString());
            }

            // Print summary statistics
            System.out.println();
            System.out.println("Number of unmatched residents: " + matchResults.unmatchedResidents.size());
            System.out.println("Number of positions available: " + matchResults.programs.stream().mapToInt(p -> p.getQuota() - p.getMatchedResidents().size()).sum());
        } catch (IOException | NumberFormatException e) {
            System.err.println("Error: " + e.getMessage());
        }
    }
}
