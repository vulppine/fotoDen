package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(albumCmd)
	rootCmd.AddCommand(folderCmd)

	albumCmd.AddCommand(albumAddCmd)
	albumAddCmd.Flags().BoolVarP(&sortf, "sort", "s", true, "sorts an album's images after adding")
	albumAddCmd.Flags().BoolVar(&tool.Genoptions.Copy, "copy", false, "toggle copying of images from source to fotoDen albums")
	albumAddCmd.Flags().BoolVar(&tool.Genoptions.Gensizes, "gensizes", true, "toggle generation of all image sizes from source to fotoDen albums")
	albumAddCmd.Flags().BoolVar(&tool.Genoptions.Meta, "meta", true, "toggle generation of metadata templates in fotoDen albums")

	albumCmd.AddCommand(albumDelCmd)

	albumCmd.AddCommand(updateCmd)
	folderCmd.AddCommand(updateCmd)

	updateCmd.AddCommand(updateInfoCmd)
	updateCmd.AddCommand(updateThumbCmd)

	updateCmd.Flags().StringVar(&nameFlag, "name", "", "set the name of a fotoDen folder/album")
	updateCmd.Flags().StringVar(&descFlag, "desc", "", "set the description of a fotoDen folder/album")
}

var (
	sortf     bool
	folderCmd = &cobra.Command{
		Use:   "folder",
		Short: "Works with fotoDen folders",
	}

	albumCmd = &cobra.Command{
		Use:   "album",
		Short: "Works with fotoDen albums",
	}
	albumAddCmd = &cobra.Command{
		Use:   "add album_name [options] images",
		Short: "Adds images to albums. Otherwise, updates the image if it exists.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if sortf {
				err := tool.InsertImage(args[0], "sort", tool.Genoptions, args[1:len(args)]...)
				return err
			}

			err := tool.InsertImage(args[0], "append", tool.Genoptions, args[1:len(args)]...)
			return err
		},
	}
	albumDelCmd = &cobra.Command{
		Use:   "delete album_name images",
		Short: "Deletes images from albums.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.DeleteImage(args[0], args[1:len(args)]...)
			return err
		},
	}
)

// update command for folders/albums

var (
	nameFlag string
	descFlag string

	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "updates fotoDen folder/album details",
	}
	updateInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "updates fotoDen album information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.UpdateFolder(args[0], nameFlag, descFlag)
			return err
		},
	}
	updateThumbCmd = &cobra.Command{
		Use:   "thumb",
		Short: "updates a fotoDen folder/album thumbnail",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.UpdateFolderThumbnail(args[0], args[1])
			return err
		},
	}
)
