package main

import (
	"github.com/stewproject/builder/cmd"
	"github.com/stewproject/stew/internals/config"
)

func main() {
	config.GetConfig()
	cmd.Execute()
}
