package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pakket-project/builder/cmd"
	"github.com/pakket-project/builder/util"
	"github.com/pakket-project/pakket/internals/config"
	pakketUtil "github.com/pakket-project/pakket/util"
)

func main() {
	// Error if not running MacOS
	if runtime.GOOS != "darwin" {
		fmt.Println("You must be on MacOS to run pakket!")
		os.Exit(1)
	}

	if runtime.GOARCH == "arm64" {
		util.Arch = "silicon"
		pakketUtil.Arch = "silicon"
	} else if runtime.GOARCH == "amd64" {
		util.Arch = "intel"
		pakketUtil.Arch = "intel"
	} else {
		fmt.Println("Unsupported architecture! Pakket only runs on Intel and Apple Silicon based Macs.")
		os.Exit(1)
	}

	if os.Getgid() != 0 && os.Getuid() != 0 {
		fmt.Println("You must run pakket as root!")
		os.Exit(1)
	}

	config.GetConfig()   // Get config
	config.GetLockfile() // Get lockfile

	cmd.Execute()
}
