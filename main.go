package main

import (
	"os"

	"sssh/internal/cli"
)

func main() {
	if err := cli.NewRoot().Execute(); err != nil {
		os.Exit(1)
	}
}
