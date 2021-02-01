package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(updCmd)

	updCmd.Flags().BoolVarP(&tool.Recurse, "recurse", "r", true, "toggles recursing through folders")
}

var (
	updCmd = &cobra.Command{
		Use:   "update { folder | web } folder",
		Args:  cobra.ExactArgs(2),
		Short: "Updates various fotoDen resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.ParseUpdate(args[0], args[1])
			if err != nil {
				return err
			}

			return nil
		},
	}
)
