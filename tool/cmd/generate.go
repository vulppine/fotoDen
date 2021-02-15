package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.AddCommand(genFolderCmd)
	genFolderCmd.Flags().StringVar(&folderMeta.Name, "name", "", "name for fotoDen folders/albums")
	genFolderCmd.Flags().StringVar(&folderMeta.Desc, "desc", "", "description for fotoDen folders/albums")
	genFolderCmd.Flags().StringVar(&tool.ThumbSrc, "thumb", "", "location of the thumbnail for the folder/album")
	genFolderCmd.Flags().BoolVar(&tool.Genoptions.Static, "static", false, "toggle more static generation of websites in fotoDen folders/albums")

	genCmd.AddCommand(genAlbumCmd)
	genAlbumCmd.Flags().StringVar(&tool.Genoptions.Source, "source", "", "source for fotoDen images")
	genAlbumCmd.Flags().StringVar(&tool.NameFlag, "name", "", "name for fotoDen folders/albums")
	genAlbumCmd.Flags().StringVar(&tool.ThumbSrc, "thumb", "", "location of the thumbnail for the folder/album")
	genAlbumCmd.Flags().BoolVar(&tool.Genoptions.Copy, "copy", false, "toggle copying of images from source to fotoDen albums")
	genAlbumCmd.Flags().BoolVar(&tool.Genoptions.Gensizes, "gensizes", true, "toggle generation of all image sizes from source to fotoDen albums")
	genAlbumCmd.Flags().BoolVar(&tool.Genoptions.Sort, "sort", true, "toggle sorting of all images in fotoDen albums by name")
	genAlbumCmd.Flags().BoolVar(&tool.Genoptions.Meta, "meta", true, "toggle generation of metadata templates in fotoDen albums")
	genAlbumCmd.Flags().BoolVar(&tool.Genoptions.Static, "static", false, "toggle more static generation of websites in fotoDen folders/albums")
}

var (
	folderMeta = tool.FolderMeta{}
	wd, _  = os.Getwd()
	genCmd = &cobra.Command{
		Use:   "generate { album | folder } destination",
		Short: "Generates fotoDen folders/albums",
	}
	genFolderCmd = &cobra.Command{
		Use:   "folder [--name string] [--thumb image] [--static] destination",
		Short: "Generates a fotoDen folder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if folderMeta.Name == "" {
				folderMeta.Name = path.Base(wd)
				err := tool.GenerateFolder(folderMeta, args[0], tool.Genoptions)
				if err != nil {
					return err
				}
			} else {
				err := tool.GenerateFolder(folderMeta, args[0], tool.Genoptions)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
	genAlbumCmd = &cobra.Command{
		Use:   "album [--name string] [--source folder] [--copy] [--sort] [--gensizes] [--meta] [--thumb image] [--static] destination",
		Short: "Generates a fotoDen album",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tool.Genoptions.ImageGen = true
			if folderMeta.Name == "" {
				folderMeta.Name = path.Base(wd)
				err := tool.GenerateFolder(folderMeta, args[0], tool.Genoptions)
				if err != nil {
					return err
				}
			} else {
				err := tool.GenerateFolder(folderMeta, args[0], tool.Genoptions)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
)
