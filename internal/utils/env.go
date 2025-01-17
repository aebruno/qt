package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func QT_VERSION() string {
	if version := os.Getenv("QT_VERSION"); version != "" {
		return version
	}

	return "5.8.0"
}

func QT_VERSION_MAJOR() string {
	return strings.Join(strings.Split(QT_VERSION(), ".")[:2], ".")
}

func QT_DIR() string {
	if dir := os.Getenv("QT_DIR"); dir != "" {
		return filepath.Clean(dir)
	}
	if runtime.GOOS == "windows" {
		return "C:\\Qt\\Qt" + QT_VERSION()
	}
	return filepath.Join(os.Getenv("HOME"), "Qt"+QT_VERSION())
}

func QT_STUB() bool {
	return strings.ToLower(os.Getenv("QT_STUB")) == "true" || UseMsys2()
}

func QT_DEBUG() bool {
	return strings.ToLower(os.Getenv("QT_DEBUG")) == "true"
}

func CheckBuildTarget(buildTarget string) {
	switch buildTarget {
	case "desktop", "android", "ios", "ios-simulator",
		"sailfish", "sailfish-emulator", "rpi1", "rpi2", "rpi3", "windows", "darwin", "linux",
		"linux-docker", "windows-docker", "android-docker":
		{
			var buildDocker = strings.HasSuffix(buildTarget, "-docker")
			switch buildTarget {
			case "windows":
				{
					if runtime.GOOS == "windows" && !buildDocker {
					} else if runtime.GOOS == "linux" || buildDocker {
					} else {
						Log.Fatalf("%v is currently not supported as a deploy target on %v", buildTarget, runtime.GOOS)
					}
				}

			case "darwin", "ios", "ios-simulator":
				{
					if runtime.GOOS == "darwin" && !buildDocker {
					} else {
						Log.Fatalf("%v is currently not supported as a deploy target on %v (not even with docker)", buildTarget, runtime.GOOS)
					}
				}

			case "linux":
				{
					if runtime.GOOS == "linux" && !buildDocker {
					} else if buildDocker {
					} else {
						Log.Fatalf("%v is currently not supported as a deploy target on %v", buildTarget, runtime.GOOS)
					}
				}
			}
		}

	default:
		{
			Log.Panicf("failed to recognize build target %v", buildTarget)
		}
	}

	if buildTarget == "android" || strings.HasPrefix(buildTarget, "ios") {
		switch {
		case UseMsys2():
			{
				Log.Fatalf("%v is not supported as a deploy target on %v with MSYS2 -> install the official Qt version instead and try again", buildTarget, runtime.GOOS)
			}

		case UseHomeBrew():
			{
				Log.Fatalf("%v is not supported as a deploy target on %v with HomeBrew -> install the official Qt version instead and try again", buildTarget, runtime.GOOS)
			}

		case UsePkgConfig():
			{
				Log.Fatalf("%v is not supported as a deploy target on %v with PkgConfig -> install the official Qt version instead and try again", buildTarget, runtime.GOOS)
			}
		}
	}
}

func IsCI() bool {
	return strings.ToLower(os.Getenv("CI")) == "true"
}
