package main

import (
	"oneclick/inspect"
	"oneclick/process"
)

func main() {
	lts := inspect.NewCheck()
	config, bl := lts.ShowConfig()
	if !bl {
		process.GetInput()
	}
	process.Verification(config)
}
