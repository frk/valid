package main

import (
	"fmt"
	"os"

	"github.com/frk/valid/cmd/internal/command"
	"github.com/frk/valid/cmd/internal/config"
)

func main() {
	cfg := config.DefaultConfig()
	cmd, err := command.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "valid: failed to initialize the command ...\n - %v\n", err)
		os.Exit(2)
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "valid: an error occurred ...\n - %v\n", err)
		os.Exit(2)
	}
}
