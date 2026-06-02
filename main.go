package main

import "fmt"

// Overridden by goreleaser via -ldflags at build time.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fmt.Printf("test-release-please %s (commit %s, built %s)\n", version, commit, date)
}
