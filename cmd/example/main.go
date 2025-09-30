package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{
		Use: "example",
	}
	root.AddCommand(
		Server(),
		Client(),
		Migrate(),
	)
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		os.Exit(1)
	}
}
