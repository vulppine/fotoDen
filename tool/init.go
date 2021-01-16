package tool

import (
	"github.com/vulppine/fotoDen/generator"
	"path"
	"fmt"
	"os"
)

// Initialization
//
// This is where the fotoDen website initialization occurs.

// InitializeWebTemplates
//
// Sets up web templates according to a given URL, and path containing templates.
// All templates should be labelled with [photo, album, folder]-template.html.

func InitializeWebTemplates(u string, srcpath string) error {

	webvars, err := generator.NewWebVars(u)
	checkError(err)

	err = generator.ConfigureWebFile(path.Join(srcpath, "photo-template.html"), path.Join(generator.CurrentConfig.WebSourceLocation, "photo.html"), webvars)
	checkError(err)
	err = generator.ConfigureWebFile(path.Join(srcpath, "album-template.html"), path.Join(generator.CurrentConfig.WebSourceLocation, "album.html"), webvars)
	checkError(err)
	err = generator.ConfigureWebFile(path.Join(srcpath, "folder-template.html"), path.Join(generator.CurrentConfig.WebSourceLocation, "folder.html"), webvars)
	checkError(err)

	err = generator.CopyFile(path.Join(srcpath, "theme.css"), "theme.css", generator.CurrentConfig.WebSourceLocation)
	checkError(err)

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

	err := generator.GenerateWebRoot(rootpath)
	if checkError(err) {
		panic(err)
	}

	err = generator.CopyFile(path.Join(generator.CurrentConfig.WebSourceLocation, "fotoDen.js"), "fotoDen.js", path.Join(rootpath, "js"))
	checkError(err)
	err = generator.CopyFile(path.Join(generator.CurrentConfig.WebSourceLocation, "theme.css"), "theme.css", path.Join(rootpath, "css"))
	checkError(err)

	webconfig := generator.NewWebConfig(*SourceFlag)
	if *SourceFlag == "" {
		fmt.Println("You will have to configure your photo storage provider in %v.", path.Join(rootpath, "config.json"))
	}
	webconfig.WorkingDirectory = path.Base(rootpath)
	webconfig.ImageThumbDir = generator.CurrentConfig.ImageThumbDirectory
	webconfig.ImageSrcDir = generator.CurrentConfig.ImageSrcDirectory
	webconfig.ImageLargeDir = generator.CurrentConfig.ImageLargeDirectory
	webconfig.PhotoExtension = ".jpg"

	err = webconfig.GenerateWebConfig(path.Join(rootpath, "config.json"))
	if checkError(err) {
		return err
	}


	checkError(err)
	err = CopyWeb("folder", rootpath)
	checkError(err)

	folder, err := generator.GenerateFolderInfo(rootpath, name) // do it in rootpath since we're not trying to scan for images in the current folder
	folder.FolderType = "folder"
	checkError(err)
	err = folder.WriteFolderInfo(path.Join(rootpath, "folderInfo.json"))
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

func InitializefotoDenConfig(u string) error {
	fmt.Println("Initializing fotoDen config with base URL: ", u)
	err := os.MkdirAll(generator.DefaultConfig.WebSourceLocation, 0755)
	if checkError(err) { panic(err) }

	generator.CurrentConfig.WebBaseURL = u

	err = generator.WritefotoDenConfig(generator.DefaultConfig, path.Join(generator.FotoDenConfigDir, "config.json"))
	if checkError(err) { panic(err) }

	return nil
}
