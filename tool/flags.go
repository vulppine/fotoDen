package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"os"
	"path"
)

// NameFlag sets the name for a folder/album. If this is not set, fotoDen will automatically use the folder's name.
var NameFlag string

// Recurse toggles recursive functions on directories. This is primarily used for the update command.
var Recurse bool

// URLFlag sets the URL for functions that require a URL. This is mostly used in initialization.
var URLFlag string

// ParseGen parses a mode, arguments, and generator options given to it, and redirects
// commands according to the given mode with its arg.
//
// Accepts two modes officially: album, and folder. Config should be deferred to initialization.
func ParseGen(mode string, arg string, options GeneratorOptions) error {
	wd, _ := os.Getwd()
	verbose("Starting from " + wd)
	if mode == "album" {
		Genoptions.imagegen = true
	}
	switch mode {
	case "album", "folder":
		switch {
		case arg == "":
			if NameFlag == "" {
				err := GenerateFolder(path.Base(wd), path.Base(wd), Genoptions)
				if checkError(err) {
					return err
				}
			} else {
				err := GenerateFolder(NameFlag, path.Base(wd), Genoptions)
				if checkError(err) {
					return err
				}
			}
		case arg != "":
			if NameFlag == "" {
				name := path.Base(wd)
				err := GenerateFolder(name, arg, Genoptions)
				if checkError(err) {
					return err
				}
			} else {
				err := GenerateFolder(NameFlag, arg, Genoptions)
				if checkError(err) {
					return err
				}
			}
		default:
			fmt.Println("Something really wrong happened. How the hell did you do that?")
			return fmt.Errorf("unexplainable error occurred")
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
			return fmt.Errorf("a config file already exists in the folder")
		}

		generator.WritefotoDenConfig(generator.DefaultConfig, path.Join(generator.FotoDenConfigDir, "config.json"))
	default:
		return fmt.Errorf("an invalid value was passed to -generate: " + mode)
	}

	return nil
}

// ParseUpdate parses a value and an arg, and directs continued execution depending on what value was given.
// This is a command line function - value is the equivalent to the second argument given (fotoDen update { value })
func ParseUpdate(value string, arg string) error {
	var err error
	switch {
	case Recurse:
		switch value {
		case "folder":
			err = RecursiveVisit(arg, UpdateFolderSubdirectories)
			if checkError(err) {
				return err
			}
		case "web":
			err = RecursiveVisit(arg, UpdateWeb)
			if checkError(err) {
				return err
			}
		}
	default:
		switch value {
		case "folder":
			err = UpdateFolderSubdirectories(arg)
			if checkError(err) {
				return err
			}
		case "web":
			err = UpdateWeb(arg)
			if checkError(err) {
				return err
			}
		}
	}

	return nil
}

/* This is here for reference
func ParseCmd() error {
	flag.Parse()
	arg := flag.Arg(0) // ignore the other flags silently

	if *verboseFlag || *verboseFlagShort {
		Verbose = true
		verbose(fmt.Sprint("Tool verbosity: ", Verbose))
		generator.Verbose = true
		verbose(fmt.Sprint("Generator verbosity: ", generator.Verbose))
	}

	genoptions := GeneratorOptions{
		Source:   *sourceFlag,
		Copy:     *copyFlag,
		Gensizes: *genSizeFlag,
		Meta:     *metaFlag,
		Static:   *staticFlag,
	}

	if *genFlag == "album" {
		genoptions.imagegen = true
	}
	verbose("Current generator options [source/copy/thumb]: " + fmt.Sprint(genoptions))

	flag.Visit(func(flag *flag.Flag) { verbose("Flag setting: " + fmt.Sprint(flag.Name, " ", flag.Value)) })

	if *configSrc == "" && *initFlag != "config" {
		err := generator.OpenfotoDenConfig(path.Join(generator.FotoDenConfigDir, "config.json"))
		if checkError(err) {
			return err
		}
	} else if *initFlag != "config" {
		err := generator.OpenfotoDenConfig(*configSrc)
		if checkError(err) {
			return err
		}
	}

	verbose("Generator config: " + fmt.Sprint(generator.CurrentConfig))

	if *initFlag != "" {
		verbose("Checking init flag..." + fmt.Sprint(*initFlag))
		switch *initFlag {
		case "config":
			err := InitializefotoDenConfig(*sourceFlag, arg)
			return err
		case "root":
			err := InitializefotoDenRoot(arg, *nameFlag)
			return err
		case "templates":
			err := InitializeWebTemplates(generator.CurrentConfig.WebBaseURL, arg)
			return err
		case "js":
			err := InitializefotoDenjs(generator.CurrentConfig.WebBaseURL, arg)
			return err
		default:
			return fmt.Errorf("invalid init flag set")
		}
	}

	switch {
	case *genFlag != "":
		err := ParseGen(*genFlag, arg, genoptions)
		if checkError(err) {
			return err
		}
	case *updFlag != "":
		err := parseUpdate(*updFlag, arg)
		if checkError(err) {
			return err
		}
	default:
		return fmt.Errorf("-init, -generate or -update must be defined")
	}

	return nil
}
*/
