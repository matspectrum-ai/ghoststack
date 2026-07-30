package main

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	fmt.Printf("GhostStack %s (%s) built %s on %s/%s\n", Version, Commit, BuildTime, runtime.GOOS, runtime.GOARCH)
}
