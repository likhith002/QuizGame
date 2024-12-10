package utility

import (
	"math/rand"
	"sort"
	"time"
)

var words = []string{
	"apple",
	"garden",
	"river",
	"cloud",
	"smile",
	"chair",
	"music",
	"light",
	"dream",
	"ocean",
	"table",
	"flower",
	"honey",
	"friend",
	"heart",
	"stone",
	"grass",
	"book",
	"moon",
	"star",
	"fire",
	"bird",
	"wind",
	"fruit",
	"wave",
	"laugh",
	"peace",
	"snow",
	"earth",
	"cake",
	"beach",
	"dance",
	"color",
	"bread",
	"night",
	"storm",
	"cloud",
	"path",
	"drum",
	"kite",
	"paint",
	"tree",
	"fish",
	"berry",
	"light",
	"road",
	"gem",
	"coat",
	"note"}

func GenerateUniqueRandomWords(generatedSets *map[string]struct{}) []string {
	rand.Seed(time.Now().UnixNano())

	for {
		// Step 1: Randomly select three words
		indices := rand.Perm(len(words))[:3] // Randomly select 3 unique indices
		selectedWords := []string{words[indices[0]], words[indices[1]], words[indices[2]]}

		// Step 2: Sort the words so that different orders don't count as different sets
		// This ensures the order of words doesn't affect the uniqueness check
		sort.Strings(selectedWords)

		// Step 3: Convert the selected words into a string key for the set
		key := selectedWords[0] + selectedWords[1] + selectedWords[2]

		// Step 4: Check if this combination has already been generated
		if _, exists := (*generatedSets)[key]; exists {
			// If a collision occurs, retry by selecting new words
			continue
		}

		// Step 5: Add the new combination to the set of generated strings
		(*generatedSets)[key] = struct{}{}

		// Step 6: Return the selected words
		return selectedWords
	}
}
