package platform

import "runtime"

type Platform string

const (
	PlatformLinux   Platform = "linux"
	PlatformWindows Platform = "windows"
	PlatformDarwin  Platform = "darwin"
	PlatformUnknown Platform = "unknown"
)

func Current() Platform {
	switch runtime.GOOS {
	case "linux":
		return PlatformLinux
	case "windows":
		return PlatformWindows
	case "darwin":
		return PlatformDarwin
	default:
		return PlatformUnknown
	}
}
