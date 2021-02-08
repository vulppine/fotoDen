package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/generator"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.AddCommand(initConfigCmd)
	initConfigCmd.Flags().StringVar(&tool.URLFlag, "url", "", "what URL to initialize fotoDen with")
	initCmd.AddCommand(initSiteCmd)
	initSiteCmd.Flags().StringVar(&tool.URLFlag, "url", "", "what URL to initialize fotoDen with")
	initSiteCmd.Flags().StringVar(&name, "name", "", "what name a site should have (with init site)")
	initCmd.AddCommand(initThemeCmd)
	initThemeCmd.Flags().StringVar(&tool.URLFlag, "url", "", "what URL to initialize fotoDen with")
	initCmd.AddCommand(initJSCmd)
}

var (
	name    string
	initCmd = &cobra.Command{
		Use:   "init { config | site | theme } destination",
		Short: "Initializes various fotoDen resources",
		Long: `Initializes fotoDen resources. Takes two args: What to initialize, and where to put it, in that order.
config creates a configuration directory in the given location.
If interactive mode is set, a wizard will appear assisting in creating the configuration,
otherwise the defaults are generated with the given URL as defined in --url.

root creates the root of a fotoDen website in the given directory.
If the name flag is not set, the name of the folder will automatically be used.

templates creates a set of templates into the current config folder (defined by config, or default).
If the url flag is not set, it will use the current configuration's base URL.

js is deprecated, and will be removed or replaced.`,
	}
	initConfigCmd = &cobra.Command{
		Use:   "config [--url url] directory",
		Short: "Initializes a fotoDen configuration directory with the given name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.InitializefotoDenConfig(tool.URLFlag, args[0])
			return err
		},
	}
	initSiteCmd = &cobra.Command{
		Use:   "site [--name] destination",
		Short: "Initializes a fotoDen website in the given directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := tool.InitializefotoDenRoot(args[0], name)
			return err
		},
	}
	initThemeCmd = &cobra.Command{
		Use:   "theme source",
		Short: "Initalizes a fotoDen theme into the configuration directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if tool.URLFlag != "" {
				err := tool.InitializeWebTheme(tool.URLFlag, args[0])
				return err
			}

			err := tool.InitializeWebTheme(generator.CurrentConfig.WebBaseURL, args[0])
			return err
		},
	}
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
)
