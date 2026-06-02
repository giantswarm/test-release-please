package main

import (
	"fmt"
	"runtime/debug"
)

// version is overridden by goreleaser via -ldflags at build time.
var version = "dev"

func main() {
	v := version
	if v == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			v = info.Main.Version
		}
	}
	fmt.Printf("test-release-please %s\n", v)
}
