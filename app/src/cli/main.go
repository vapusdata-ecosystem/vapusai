package main

import (
	"fmt"
	"os"

	cmd "github.com/vapusdata-ecosystem/vapusai-studio/cli/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
