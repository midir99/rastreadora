package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/midir99/rastreadora/cmd"
)

func main() {
	args, err := cmd.ParseArgs()
	if err != nil {
		fmt.Fprint(
			flag.CommandLine.Output(),
			"Error: ",
			err,
			"\nTry using the -h flag to get some help.\n",
		)
		os.Exit(1)
	}
	cmd.Execute(args)
}
