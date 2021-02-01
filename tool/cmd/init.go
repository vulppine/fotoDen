package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/generator"
	"github.com/vulppine/fotoDen/tool"
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&tool.URLFlag, "url", "", "what URL to initialize fotoDen with")
	initCmd.Flags().StringVar(&name, "name", "", "what name a site should have (with init site)")
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
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "config":
				err := tool.InitializefotoDenConfig(tool.URLFlag, args[1])
				return err
			case "site":
				err := tool.InitializefotoDenRoot(args[1], name)
				return err
			case "theme":
				if tool.URLFlag != "" {
					err := tool.InitializeWebTheme(tool.URLFlag, args[1])
					return err
				}

				err := tool.InitializeWebTheme(generator.CurrentConfig.WebBaseURL, args[1])
				return err
			case "js":
				if tool.URLFlag != "" {
					err := tool.InitializefotoDenjs(tool.URLFlag, args[1])
					return err
				}

				err := tool.InitializefotoDenjs(generator.CurrentConfig.WebBaseURL, args[1])
				return err
			default:
				return fmt.Errorf("invalid init flag set")
			}
		},
	}
)
