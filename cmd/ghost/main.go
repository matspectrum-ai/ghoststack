package main

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/cli"
)


var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	root := cli.NewRootCommand()
	root.Version = fmt.Sprintf("%s (%s) built %s", Version, Commit, BuildTime)
	root.SetVersionTemplate("GhostStack {{.Version}}\n")

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
