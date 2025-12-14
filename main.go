package main

import (
	"os"

	"github.com/manattan/clove/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
