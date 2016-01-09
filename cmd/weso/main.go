package main

import (
	"fmt"
	"os"

	"github.com/morikuni/weso"
)

func main() {
	conf, ok := weso.NewConfig(os.Args[1:], os.Stderr)
	if !ok {
		os.Exit(1)
	}

	cli, err := weso.NewCLI(conf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.Close()
}
