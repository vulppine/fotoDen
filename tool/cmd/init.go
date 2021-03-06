package cmd

import (
	"github.com/spf13/cobra"
	// "github.com/vulppine/fotoDen/generator"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(initCmd)

	// initCmd.AddCommand(initConfigCmd)
	// initConfigCmd.Flags().StringVar(&tool.URLFlag, "url", "", "what URL to initialize fotoDen with")
	initCmd.AddCommand(initSiteCmd)
	initSiteCmd.Flags().StringVar(&tool.URLFlag, "source", "", "where fotoDen should obtain its images")
	initSiteCmd.Flags().StringVar(&websiteInit.URL, "url", "", "what URL to initialize fotoDen with")
	initSiteCmd.Flags().StringVar(&websiteInit.Name, "name", "", "what name a site should have (with init site)")
	initSiteCmd.Flags().StringVar(&websiteInit.Theme, "theme", "", "what theme a site should use")
	// initCmd.AddCommand(initThemeCmd)
	// initThemeCmd.Flags().StringVar(&tool.URLFlag, "url", "", "what URL to initialize fotoDen with")
	// initCmd.AddCommand(initJSCmd)
}

var (
	initCmd = &cobra.Command{
		Use:   "init { site | theme } destination",
		Short: "Initializes various fotoDen resources",
		Long: `Initializes fotoDen resources. Takes two args: What to initialize, and where to put it, in that order.`,
	}
	/* Deprecated, see init.go in tool root
	initConfigCmd = &cobra.Command{
		Use:   "config [--url url] directory",
		Short: "Initializes a fotoDen configuration directory with the given name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.InitializefotoDenConfig(tool.URLFlag, args[0])
			return err
		},
	}
	*/
	websiteInit tool.WebsiteConfig
	initSiteCmd = &cobra.Command{
		Use:   "site [--name] destination",
		Short: "Initializes a fotoDen website in the given directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.InitializefotoDenRoot(args[0], websiteInit)
			return err
		},
	}
	/* Deprecated in favor of zipped/embed theme files
	initThemeCmd = &cobra.Command{
		Use:   "theme source",
		Short: "Initalizes a fotoDen theme into the configuration directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.CopyThemeToConfig(args[0])
			return err
		},
	}
	*/
	/*
		initJSCmd = &cobra.Command{
			Use:   "js source",
			Short: "Checks and copies over a fotoDen.js file into the configuration directory",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				err := tool.InitializefotoDenjs(args[0])
				if err != nil {
					return err
				}

				return nil
			},
		}
	*/
)
