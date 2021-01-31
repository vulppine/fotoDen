package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVar(&tool.Genoptions.Source, "source", "", "source for fotoDen images")
	genCmd.Flags().StringVar(&tool.NameFlag, "name", "", "name for fotoDen folders/albums")
	genCmd.Flags().BoolVar(&tool.Genoptions.Copy, "copy", false, "toggle copying of images from source to fotoDen albums")
	genCmd.Flags().BoolVar(&tool.Genoptions.Gensizes, "gensizes", true, "toggle generation of all image sizes from source to fotoDen albums")
	genCmd.Flags().BoolVar(&tool.Genoptions.Sort, "sort", true, "toggle sorting of all images in fotoDen albums by name")
	genCmd.Flags().BoolVar(&tool.Genoptions.Meta, "meta", true, "toggle generation of metadata templates in fotoDen albums")
	genCmd.Flags().BoolVar(&tool.Genoptions.Static, "static", false, "toggle more static generation of websites in fotoDen folders/albums")
}

var (
	genCmd = &cobra.Command{
		Use: "generate",
		Short: "Generates fotoDen folders/albums",
		Args: cobra.ExactArgs(2),
		RunE: func (cmd *cobra.Command, args []string) error {
			err := tool.ParseGen(args[0], args[1], tool.Genoptions)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
