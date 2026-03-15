package main

import (
	"os"

	"github.com/Dokkabei97/notebooklm-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
