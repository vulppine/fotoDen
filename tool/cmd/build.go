package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var (
	buildCmd = &cobra.Command{
		Use:   "build buildfile destination",
		Short: "Generates a fotoDen folder/album from a valid YAML file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := new(tool.BuildFile)
			err := b.OpenBuildYAML(args[0])
			if err != nil {
				return err
			}

			err = b.Build(args[1])
			if err != nil {
				return err
			}

			return nil
		},
	}
)
