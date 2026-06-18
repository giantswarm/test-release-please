package main

import (
	"fmt"
	"os"

	"github.com/giantswarm/test-release-please/pkg/project"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("version %s\n", project.Version())
		fmt.Printf("  commit: %s\n", project.GitSHA())
		fmt.Printf("  built:  %s\n", project.BuildTimestamp())
		return
	}
	fmt.Println("usage: app version")
}
