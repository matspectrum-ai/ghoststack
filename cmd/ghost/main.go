package main

import (
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/cli"
)

func main() {
	root := cli.NewRootCommand()
	root.Version = fmt.Sprintf("%s (%s) built %s", cli.Version, cli.Commit, cli.BuildTime)
	root.SetVersionTemplate("GhostStack {{.Version}}\n")

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
