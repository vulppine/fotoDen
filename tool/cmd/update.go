package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(updCmd)

	updCmd.Flags().BoolP(&tool.Recurse, "recurse", "r", true, "toggles recursing through folders")
}

var (
	updCmd = &cobra.Command{
		Use: "update type",
		Short: "Updates various fotoDen resources",
		RunE: func (cmd *cobra.Command, args []string) error {
			err := tool.ParseUpdate(args[0], args[1])
			if err != nil {
				return err
			}

			return nil
		},
	}
)
