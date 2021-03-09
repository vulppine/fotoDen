package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vulppine/fotoDen/generator"
	"github.com/vulppine/fotoDen/tool"
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
	d         = rootCmd.PersistentFlags().Bool("debug", false, "Prints debug information to console.")
	v         = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Prints verbose information made by fotoDen")
	configDir string
	configSrc generator.Config
	site      string
	rootCmd   = &cobra.Command{
		Use:   "fotoDen { init | generate | update } args [--config string] [--verbose | -v] [--interactive | -i]",
		Short: "A static photo gallery generator",
	}
)

func setRootFlags() {
	if *v {
		tool.Verbose = true
		generator.Verbose = true
	}

	if configDir == "" {
		configDir = generator.RootConfigDir
	} else {
		generator.RootConfigDir, _ = filepath.Abs(configDir)
	}

	if site == "" {
		_, err := os.Stat(filepath.Join(generator.RootConfigDir, "defaultsite"))
		if os.IsNotExist(err) {
			verbose("WARNING: defaultsite does not exist")
			site = "___NOSITE"
		} else {
			f, err := os.Open(filepath.Join(generator.RootConfigDir, "defaultsite"))
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			d, err := ioutil.ReadAll(f)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			site = string(d)
		}
	}

	if site != "___NOSITE" {
		verbose(filepath.Join(configDir, "sites", site, "config.json"))
		s := new(tool.WebsiteConfig)
		err := generator.ReadJSON(filepath.Join(generator.RootConfigDir, "sites", site, "config.json"), s)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		generator.CurrentConfig = s.GeneratorConfig
	}
}

func init() {
	cobra.OnInitialize(setRootFlags)
	rootCmd.PersistentFlags().BoolVarP(&tool.WizardFlag, "interactive", "i", false, "Allows fotoDen to display interactive prompts")
	rootCmd.PersistentFlags().StringVar(&configDir, "config-dir", "", "The config directory to use for fotoDen")
	rootCmd.PersistentFlags().StringVar(&site, "site", "", "The website that fotoDen should focus on.")
}
