package main

import (
	"fmt"

	"github.com/sunggun-yu/envp/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// set version
	cmd.Version = Version()
	cmd.Execute()
}

// Version returns version and build information. it's injected from ldflags
func Version() string {
	return fmt.Sprintf("prw %s, commit %s, built at %s", version, commit, date)
}
