package main

import (
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{
		Use: "example",
	}
	root.AddCommand(
		Server(),
		Client(),
	)
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
