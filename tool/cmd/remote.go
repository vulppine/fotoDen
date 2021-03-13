package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(remoteCmd)
}

var (
	remoteCmd = &cobra.Command{
		Use: "remote",
		Short: "Utilities for remotely uploading to a webhost",
	}
)
