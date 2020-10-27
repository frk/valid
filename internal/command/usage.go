package command

import (
	"fmt"
	"os"
)

func printUsage() {
	fmt.Fprint(os.Stderr, usage)
}

const usage = `usage: isvalid ...
` //`
