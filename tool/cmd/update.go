package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(updCmd)

	updCmd.AddCommand(updFolderCmd)
	updCmd.AddCommand(updWebCmd)
	updFolderCmd.Flags().BoolVarP(&tool.Recurse, "recurse", "r", true, "toggles recursing through folders")
	updWebCmd.Flags().BoolVarP(&tool.Recurse, "recurse", "r", true, "toggles recursing through folders")
}

var (
	updCmd = &cobra.Command{
		Use:   "update { folder | web } folder",
		Short: "Updates various fotoDen resources",
	}
	updFolderCmd = &cobra.Command{
		Use:   "folder [-r] folder_name",
		Args:  cobra.ExactArgs(1),
		Short: "Updates fotoDen folder subdirectories",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tool.Recurse {
				err := tool.RecursiveVisit(args[0], tool.UpdateFolderSubdirectories)
				if err != nil {
					return err
				}
			}

			err := tool.UpdateFolderSubdirectories(args[0])
			if err != nil {
				return err
			}

			return nil
		},
	}
	updWebCmd = &cobra.Command{
		Use:   "web [-r] folder_name",
		Args:  cobra.ExactArgs(1),
		Short: "Updates fotoDen folder webpages",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tool.Recurse {
				err := tool.RecursiveVisit(args[0], tool.UpdateWeb)
				if err != nil {
					return err
				}
			}

			err := tool.UpdateWeb(args[0])
			if err != nil {
				return err
			}

			return nil
		},
	}
)
