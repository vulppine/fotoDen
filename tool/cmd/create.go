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
	genFolderCmd.Flags().BoolVar(&opts.Static, "static", false, "toggle more static generation of websites in fotoDen folders/albums")

	genCmd.AddCommand(genAlbumCmd)
	genAlbumCmd.Flags().StringVar(&opts.Source, "source", "", "source for fotoDen images")
	genAlbumCmd.Flags().StringVar(&folderMeta.Name, "name", "", "name for fotoDen folders/albums")
	genAlbumCmd.Flags().StringVar(&folderMeta.Desc, "desc", "", "description for fotoDen folders/albums")
	genAlbumCmd.Flags().StringVar(&tool.ThumbSrc, "thumb", "", "location of the thumbnail for the folder/album")
	genAlbumCmd.Flags().BoolVar(&opts.Copy, "copy", false, "toggle copying of images from source to fotoDen albums")
	genAlbumCmd.Flags().BoolVar(&opts.Gensizes, "gensizes", true, "toggle generation of all image sizes from source to fotoDen albums")
	genAlbumCmd.Flags().BoolVar(&opts.Sort, "sort", true, "toggle sorting of all images in fotoDen albums by name")
	genAlbumCmd.Flags().BoolVar(&opts.Meta, "meta", true, "toggle generation of metadata templates in fotoDen albums")
	genAlbumCmd.Flags().BoolVar(&opts.Static, "static", false, "toggle more static generation of websites in fotoDen folders/albums")

	genCmd.AddCommand(genPageCmd)
	genPageCmd.Flags().StringVar(&t, "name", "", "the name of the webpage (used as title)")
}

// TODO:
// - Reorganize this into a 'create' command, which handles:
//   - Folders
//   - Albums
//   - Pages
//   - Sites
//
// ^ Maybe - splitting creation and site initialization might be
// a good thing, so that it's more distinct that you
// 'initialize' a site and 'create' things inside of it.

var (
	t          string
	opts       tool.GeneratorOptions
	folderMeta = tool.FolderMeta{}
	wd, _      = os.Getwd()
	genCmd     = &cobra.Command{
		Use:   "create { album | folder | page }",
		Short: "Creates fotoDen folders/albums/pages",
	}
	genFolderCmd = &cobra.Command{
		Use:   "folder [--name string] [--thumb image] [--static] destination",
		Short: "Creates a fotoDen folder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if folderMeta.Name == "" {
				folderMeta.Name = path.Base(wd)
				err := tool.GenerateFolder(folderMeta, args[0], opts)
				if err != nil {
					return err
				}
			} else {
				err := tool.GenerateFolder(folderMeta, args[0], opts)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
	genAlbumCmd = &cobra.Command{
		Use:   "album [--name string] [--copy] [--sort] [--gensizes] [--meta] [--thumb image] [--static] source destination",
		Short: "Creates a fotoDen album",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ImageGen = true
			opts.Source = args[0]

			if folderMeta.Name == "" {
				folderMeta.Name = path.Base(wd)
				err := tool.GenerateFolder(folderMeta, args[1], opts)
				if err != nil {
					return err
				}
			} else {
				err := tool.GenerateFolder(folderMeta, args[1], opts)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
	genPageCmd = &cobra.Command{
		Use:   "page [--name string] source",
		Short: "Creates a webpage using a fotoDen template and Markdown (incomplete)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.GeneratePage(args[0], t)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
