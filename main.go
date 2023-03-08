package main

import (
	"context"
	"os"

	"github.com/pandodao/botastic-cli/cmd/root"
	"github.com/spf13/cobra"
)

var (
	version = "0.0.1"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "version", version)

	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	rootCmd := root.NewCmdRoot(version)

	rootCmd.SetArgs(expandedArgs)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		rootCmd.PrintErrln("execute failed:", err)
		os.Exit(1)
	}
}

func hasCommand(rootCmd *cobra.Command, args []string) bool {
	c, _, err := rootCmd.Traverse(args)
	return err == nil && c != rootCmd
}
