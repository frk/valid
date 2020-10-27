// isvalid is a tool for generating struct field validation.
package main

import (
	"fmt"
	"os"

	"github.com/frk/isvalid/internal/command"
)

func main() {
	conf := command.DefaultConfig
	conf.ParseFlags()
	if err := conf.ParseFile(); err != nil {
		fmt.Fprintf(os.Stderr, "isvalid: failed parsing config file ...\n - %v\n", err)
		os.Exit(2)
	}

	cmd, err := command.New(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "isvalid: failed to initialize the command ...\n - %v\n", err)
		os.Exit(2)
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "isvalid: an error occurred ...\n - %v\n", err)
		os.Exit(2)
	}
}
