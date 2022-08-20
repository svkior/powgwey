package quotes

import (
	"os"
)

// from:
// https://www.tutorialspoint.com/how-to-check-if-a-file-exists-in-golang
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
