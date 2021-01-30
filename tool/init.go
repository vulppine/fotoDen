package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"os"
	"path"
	"path/filepath"
)

// Initialization
//
// This is where the fotoDen website initialization occurs.

// InitializeWebTemplates
//
// Sets up web templates according to a given URL, and path containing templates.
// All templates should be labelled with [photo, album, folder]-template.html.
func InitializeWebTemplates(u string, srcpath string) error {
	wd, _ := os.Getwd()

	t, err := ReadThemeConfig(path.Join(srcpath, "theme.json"))
	webvars, err := generator.NewWebVars(u)
	checkError(err)

	err = generator.CopyFile(path.Join(srcpath, "theme.json"), "theme.json", path.Join(generator.CurrentConfig.WebSourceLocation))

	err = os.Mkdir(path.Join(generator.CurrentConfig.WebSourceLocation, "html"), 0755)
	if checkError(err) {
		return err
	}

	err = generator.ConfigureWebFile(path.Join(srcpath, "html", "photo-template.html"), path.Join(generator.CurrentConfig.WebSourceLocation, "html", "photo-template.html"), webvars)
	checkError(err)
	err = generator.ConfigureWebFile(path.Join(srcpath, "html", "album-template.html"), path.Join(generator.CurrentConfig.WebSourceLocation, "html", "album-template.html"), webvars)
	checkError(err)
	err = generator.ConfigureWebFile(path.Join(srcpath, "html", "folder-template.html"), path.Join(generator.CurrentConfig.WebSourceLocation, "html", "folder-template.html"), webvars)
	checkError(err)

	if len(t.Stylesheets) != 0 {
		err = os.Mkdir(path.Join(generator.CurrentConfig.WebSourceLocation, "css"), 0755)
		if checkError(err) {
			return err
		}

		os.Chdir(path.Join(srcpath, "css"))
		err = generator.BatchCopyFile(t.Stylesheets, path.Join(generator.CurrentConfig.WebSourceLocation, "css"))
		checkError(err)
		os.Chdir(wd)
	}

	if len(t.Scripts) != 0 {
		err = os.Mkdir(path.Join(generator.CurrentConfig.WebSourceLocation, "js"), 0755)
		if checkError(err) {
			return err
		}

		os.Chdir(path.Join(srcpath, "js"))
		err = generator.BatchCopyFile(t.Scripts, path.Join(generator.CurrentConfig.WebSourceLocation, "js"))
		checkError(err)
		os.Chdir(wd)
	}

	return nil
}

// Initialize fotoDen.js
//
// Sets a single variable as needed in fotoDen.js, from a path where it is located.
// Copies it over to generator.CurrentConfig.WebSourceLocation afterwards.
func InitializefotoDenjs(u string, fpath string) error {

	webvars, err := generator.NewWebVars(u)

	checkError(err)

	err = generator.ConfigureWebFile(path.Join(fpath), path.Join(generator.CurrentConfig.WebSourceLocation, "fotoDen.js"), webvars)
	if checkError(err) {
		return err
	}

	return nil
}

// Initialize fotoDen root
//
// Sets up the root directory for fotoDen, including a folderInfo.json file.
// Creates the folder structure,
// copies over the folder page as well as the theme.css file,
// and copies over fotoDen.js,
// and then generates a folderInfo.json file according to the name given in the command line,
// otherwise generates with a blank name.
//
// By default, the root of a generated fotoDen website should be specifically a folder.
// A fotoDen page can be anything, as long as the required tags are there,
// but if it is being generated via this tool,
// it is, by default, a folder.
func InitializefotoDenRoot(rootpath string, name string) error {
	rootpath, _ = filepath.Abs(rootpath)

	err := generator.GenerateWebRoot(rootpath)
	if checkError(err) {
		panic(err)
	}

	err = generator.CopyFile(path.Join(generator.CurrentConfig.WebSourceLocation, "fotoDen.js"), "fotoDen.js", path.Join(rootpath, "js"))
	checkError(err)

	t, err := ReadThemeConfig(path.Join(generator.CurrentConfig.WebSourceLocation, "theme.json"))
	if checkError(err) {
		return err
	}

	wd, _ := os.Getwd()

	os.Chdir(generator.CurrentConfig.WebSourceLocation)
	if len(t.Stylesheets) != 0 {
		os.Chdir("css")
		err = generator.BatchCopyFile(t.Stylesheets, path.Join(rootpath, "theme", "css"))
		if checkError(err) {
			return err
		}
		os.Chdir(generator.CurrentConfig.WebSourceLocation)
	}

	if len(t.Scripts) != 0 {
		os.Chdir("js")
		err = generator.BatchCopyFile(t.Scripts, path.Join(rootpath, "theme", "js"))
		if checkError(err) {
			return err
		}
	}

	os.Chdir(wd)

	var webconfig *generator.WebConfig

	if *wizardFlag == true {
		webconfig = SetupWebConfig(*sourceFlag)
	} else {
		webconfig = generator.GenerateWebConfig(*sourceFlag)
		if *sourceFlag == "" {
			fmt.Printf("You will have to configure your photo storage provider in %v.", path.Join(rootpath, "config.json"))
		}
	}

	err = webconfig.WriteWebConfig(path.Join(rootpath, "config.json"))
	if checkError(err) {
		return err
	}

	folder, err := generator.GenerateFolderInfo(rootpath, name) // do it in rootpath since we're not trying to scan for images in the current folder
	folder.FolderType = "folder"
	checkError(err)
	err = folder.WriteFolderInfo(path.Join(rootpath, "folderInfo.json"))
	checkError(err)

	err = GenerateWeb("folder", rootpath, folder, genoptions)
	checkError(err)

	return nil
}

// Initialize fotoDen config folder
//
// Initializes the fotoDen config folder.
//
// This should only be done once.
//
// Takes a single string to set WebBaseURL as.
func InitializefotoDenConfig(u string, dest string) error {
	fmt.Println("Initializing fotoDen config with base URL: ", u)

	var config generator.GeneratorConfig

	if *wizardFlag == true {
		config = SetupConfig()
	} else {
		config = generator.DefaultConfig
		generator.CurrentConfig.WebBaseURL = u
	}

	if dest == "" {
		dest = generator.FotoDenConfigDir
	} else {
		dest, err := filepath.Abs(dest)
		if checkError(err) {
			return err
		}
		config.WebSourceLocation = path.Join(dest, "web")
	}

	err := os.MkdirAll(config.WebSourceLocation, 0755)
	if checkError(err) {
		return err
	}

	err = generator.WritefotoDenConfig(config, path.Join(dest, "config.json"))
	if checkError(err) {
		return err
	}

	return nil
}
