package semver

import (
	"sort"
)

<<<<<<< HEAD
type Versions []Version

=======
// Versions represents multiple versions.
type Versions []Version

// Len returns length of version collection
>>>>>>> 2273e7a... chore(vendor): update
func (s Versions) Len() int {
	return len(s)
}

<<<<<<< HEAD
=======
// Swap swaps two versions inside the collection by its indices
>>>>>>> 2273e7a... chore(vendor): update
func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

<<<<<<< HEAD
=======
// Less checks if version at index i is less than version at index j
>>>>>>> 2273e7a... chore(vendor): update
func (s Versions) Less(i, j int) bool {
	return s[i].LT(s[j])
}

// Sort sorts a slice of versions
func Sort(versions []Version) {
	sort.Sort(Versions(versions))
}
