package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/tool"
	"github.com/vulppine/fotoDen/generator"
)

func Execute() error {
	return rootCmd.Execute()
}

func verbose(input string) {
	if *v {
		log.Println(input)
	}
}

func debug(input interface{}) {
	if *d {
		log.Println(input)
	}
}

var (
	d = rootCmd.PersistentFlags().Bool("debug", false, "Prints debug information to console.")
	v = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Prints verbose information made by fotoDen")
	config string
	configSrc generator.GeneratorConfig
	rootCmd = &cobra.Command{
		Use: "fotoDen",
		Short: "A static photo gallery generator",
	}
)

func setRootFlags() {
	if *v {
		tool.Verbose = true
		generator.Verbose = true
	}

	if config != "" {
		verbose(filepath.Join(config, "config.json"))
		err := generator.OpenfotoDenConfig(filepath.Join(config, "config.json"))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else {
		err := generator.OpenfotoDenConfig(filepath.Join(generator.FotoDenConfigDir, "config.json"))
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}

func init() {
	cobra.OnInitialize(setRootFlags)
	rootCmd.PersistentFlags().BoolVarP(&tool.WizardFlag, "interactive", "i", true, "Allows fotoDen to display interactive prompts")
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "The config directory to use for fotoDen")
}
