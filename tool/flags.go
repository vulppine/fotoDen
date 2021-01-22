package tool

import (
	"flag"
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"os"
	"path"
)

// Flags
//
// The command line flags.
//
// fotoDen tool requires -generate or -update passed to it, otherwise it won't do anything.
// A single arg is required, and that should be a the name of the folder being processed.
// If this isn't detected when -generate is passed, it will use the name of the current folder.
// Equally, -generate can also take a 'name' flag - if this isn't detected, the name of
// the current folder will be used.
//
// -generate is currently a bool, but should be a string later for getting photos from a source
// -update is a bool, for now
//
// -generate will generate a fotoDen folder in the current folder with a given shortname.
//
// -update updates the folder. (If the given args are null, or if there is no info file, it will return an error.)

var GenFlag = flag.String("generate", "", "Generates a fotoDen structure in the current folder, or the default config in the configuration directory. Accepted modes: album, folder, config.")
var NameFlag = flag.String("name", "", "The name of the folder (not the path). If this is blank, or not called, the current name of the folder will be used in generation.")
var UpdFlag = flag.Bool("update", false, "Updates the fotoDen folder. (Currently only updates the folder subdirectories)")
var SourceFlag = flag.String("source", "", "The source used for fotoDen images. This is multi-context - calling this during -generate full will take images from the source directory as its base, and calling this during -init root will use this as the fotoDen image storage provider.")
var CopyFlag = flag.Bool("copy", false, "Copies files over to GeneratorConfig.ImageSrcDirectory. Useful if you're copying over to a remote directory.")
var ThumbSrc = flag.String("folthumb", "", "The name of the thumbnail in the source directory. This will be selected as the thumbnail of the folder, and is copied over to the root of the folder.")
var GenSizeFlag = flag.Bool("gensizes", true, "Tells the generator to generate all sizes in the config. This is automatically set to true.")
var ConfigSrc = flag.String("config", "", "The name of the config file to use. If this isn't set, the one is $CONFIG/fotoDen is used - otherwise, an error is returned. Call 'fotoDen -generate config' to create a config in either $CONFIG/fotoden, or in a relative folder if defined.")
var InitFlag = flag.String("init", "", "Initializes various aspects of fotoDen. Accepted values: config, root, templates, js. config should only be done if the config folder was removed, as it is automatically called at the first start of the program.")
var WizardFlag = flag.Bool("interactive", true, "Enables interactive mode. Interactive mode occurs when settings need to be configured in files.")
var VerboseFlag = flag.Bool("verbose", false, "Sets verbose mode.")
var VerboseFlagShort = flag.Bool("v", false, "Sets verbose mode.")

func ParseGen(mode string, arg string, options GeneratorOptions) error {
	wd, _ := os.Getwd()
	verbose("Starting from " + wd)
	switch mode {
	case "album", "folder":
		switch {
		case arg == "":
			if *NameFlag == "" {
				err := GenerateFolder(path.Base(wd), path.Base(wd), options)
				if checkError(err) {
					return err
				}
			} else {
				err := GenerateFolder(*NameFlag, path.Base(wd), options)
				if checkError(err) {
					return err
				}
			}
		case arg != "":
			if *NameFlag == "" {
				name := path.Base(wd)
				err := GenerateFolder(name, arg, options)
				if checkError(err) {
					return err
				}
			} else {
				err := GenerateFolder(*NameFlag, arg, options)
				if checkError(err) {
					return err
				}
			}
		default:
			fmt.Println("Something really wrong happened. How the hell did you do that?")
			return fmt.Errorf("Unexplainable error occurred")
		}
	case "config":
		var file *os.File
		var err error
		if arg != "" {
			file, err = os.OpenFile(path.Join(arg, "config.json"), os.O_RDWR|os.O_CREATE, 0755)
			checkError(err)
		} else {
			file, err = os.OpenFile(path.Join(generator.FotoDenConfigDir, "config.json"), os.O_RDWR|os.O_CREATE, 0755)
			checkError(err)
		}
		defer file.Close()
		stat, err := file.Stat()
		checkError(err)

		if stat.Size() > 0 {
			return fmt.Errorf("A config file already exists in the folder. Consider modifying it with a text editor.")
		} else {
			generator.WritefotoDenConfig(generator.DefaultConfig, path.Join(generator.FotoDenConfigDir, "config.json"))
		}
	default:
		return fmt.Errorf("An invalid value was passed to -generate. Aborting. mode: " + mode)
	}

	return nil
}

func ParseCmd() error {
	flag.Parse()
	arg := flag.Arg(0) // ignore the other flags silently

	if *VerboseFlag || *VerboseFlagShort {
		Verbose = true
		verbose(fmt.Sprint("Tool verbosity: ", Verbose))
		generator.Verbose = true
		verbose(fmt.Sprint("Generator verbosity: ", generator.Verbose))
	}

	genoptions := GeneratorOptions{
		source: *SourceFlag,
		copy: *CopyFlag,
		gensizes: *GenSizeFlag,
	}

	if *GenFlag == "folder" {
		genoptions.imagegen = false
	}
	verbose("Current generator options [source/copy/thumb]: " + fmt.Sprint(genoptions))

	flag.Visit(func(flag *flag.Flag) { verbose("Flag setting: " + fmt.Sprint(flag.Name, " ", flag.Value)) })

	if *ConfigSrc == "" && *InitFlag != "config" {
		err := generator.OpenfotoDenConfig(path.Join(generator.FotoDenConfigDir, "config.json"))
		verbose("Generator config: " + fmt.Sprint(generator.CurrentConfig))
		if checkError(err) {
			return err
		}
	}

	if *InitFlag != "" {
		verbose("Checking init flag..." + fmt.Sprint(*InitFlag))
		switch *InitFlag {
		case "config":
			err := InitializefotoDenConfig(*SourceFlag, arg)
			return err
		case "root":
			err := InitializefotoDenRoot(arg, *NameFlag)
			return err
		case "templates":
			err := InitializeWebTemplates(generator.CurrentConfig.WebBaseURL, arg)
			return err
		case "js":
			err := InitializefotoDenjs(generator.CurrentConfig.WebBaseURL, arg)
			return err
		default:
			return fmt.Errorf("An invalid init flag was set. Aborting.")
		}
	}

	switch {
	case *GenFlag != "":
		err := ParseGen(*GenFlag, arg, genoptions)
		if checkError(err) {
			return err
		}
	case *UpdFlag:
		err := UpdateFolderSubdirectories(arg)
		if checkError(err) {
			return err
		}
	default:
		return fmt.Errorf("-init, -generate or -update must be defined. Aborting.")
	}

	return nil
}
