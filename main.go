package main

import (
	"fmt"
	"os"

	"github.com/manattan/clove/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stdout, "エラー: %v\n", err)
		os.Exit(1)
	}
}
