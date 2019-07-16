package util

import (
	"os"
	"path/filepath"
)

// Dirname gets the absolute path to the running application.
// Returns a string with the absolute path to the application.
func Dirname() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return filepath.Dir(os.Args[0])
	}
	return path
}

// DirnameJoin join multiple path segments with the application dirname.
// `segments` is a list of path segments.
// Returns the full absolute path to solve folder.
func DirnameJoin(segments ...string) string {
	segments = append([]string{Dirname()}, segments...)
	joined := filepath.Join(segments...)
	return joined
}
