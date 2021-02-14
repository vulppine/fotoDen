package tool

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/vulppine/fotoDen/generator"
)

// InitializeWebTheme sets up web templates according to a given URL, and path containing templates.
// All templates should be labelled with [photo, album, folder]-template.html.
func InitializeWebTheme(u string, srcpath string) error {
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

var (
	// The MD5 checksum of fotoDen.js. Must be defined during build.
	JSSum string
	// The MD5 checksum of fotoDen.min.js. Must be defined during build.
	JSMinSum string
)

// InitializefotoDenjs copies over a valid fotoDen.js into the configuration's web directory.
// The JS must be validated during build, otherwise it will warn the user that the script cannot be validated.
func InitializefotoDenjs(fpath string) error {
	if JSSum != "" || JSMinSum != "" {
		f, err := ioutil.ReadFile(fpath)
		if checkError(err) {
			return err
		}
		s := fmt.Sprintf("%x", md5.Sum(f))
		if s != JSSum && s != JSMinSum {
			verbose("checksum: " + s)
			verbose("valid checksums:")
			verbose("fotoDen.js: " + JSSum)
			verbose("fotoDen.min.js: " + JSMinSum)
			if !ReadInputAsBool("Warning: the selected JS file may have been modified, or is not fotoDen.js. Continue?", "y") {
				return fmt.Errorf("js not valid")
			}
		}
	} else {
		log.Println("Warning: fotoDen tool was not compiled with valid JS checksums. Ensure that your scripts are safe, or are valid for fotoDen websites.")
	}

	err := generator.CopyFile(path.Join(fpath), "fotoDen.js", path.Join(generator.CurrentConfig.WebSourceLocation))
	if checkError(err) {
		return err
	}

	return nil
}

// InitializefotoDenRoot sets up the root directory for fotoDen, including a folderInfo.json file.
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
		os.Chdir(generator.CurrentConfig.WebSourceLocation)
	}

	if len(t.Other) != 0 {
		os.Chdir("etc")
		err = generator.BatchCopyFile(t.Stylesheets, path.Join(rootpath, "theme", "etc"))
		if checkError(err) {
			return err
		}
	}

	os.Chdir(wd)

	var webconfig *generator.WebConfig

	if WizardFlag == true {
		webconfig = setupWebConfig(URLFlag)
	} else {
		webconfig = generator.GenerateWebConfig(URLFlag)
		webconfig.Theme = true // we're generating this from fotoDen tool, so we're using a theme obviously
		webconfig.WebsiteTitle = name
		webconfig.ImageRootDir = generator.CurrentConfig.ImageRootDirectory
		webconfig.ThumbnailFrom = webconfig.ImageSizes[0].SizeName
		webconfig.DisplayImageFrom = webconfig.ImageSizes[len(webconfig.ImageSizes)-1].SizeName
		for _, v := range webconfig.ImageSizes {
			webconfig.DownloadSizes = append(webconfig.DownloadSizes, v.SizeName)
		}

		if URLFlag == "" {
			fmt.Printf("You will have to configure your photo storage provider in %v.", path.Join(rootpath, "config.json"))
		}
	}

	err = webconfig.WriteWebConfig(path.Join(rootpath, "config.json"))
	if checkError(err) {
		return err
	}

	folder, err := generator.GenerateFolderInfo(rootpath, webconfig.WebsiteTitle) // do it in rootpath since we're not trying to scan for images in the current folder
	folder.Type = "folder"
	checkError(err)
	err = folder.WriteFolderInfo(path.Join(rootpath, "folderInfo.json"))
	checkError(err)

	err = GenerateWeb("folder", rootpath, folder, Genoptions)
	checkError(err)

	return nil
}

// InitializefotoDenConfig initializes a fotoDen config folder.
//
// This should only be done once.
//
// Takes a single string to set WebBaseURL as.
func InitializefotoDenConfig(u string, dest string) error {
	fmt.Println("Initializing fotoDen config with base URL: ", u)

	var config generator.Config

	if WizardFlag == true {
		config = setupConfig()
	} else {
		config = generator.DefaultConfig
		config.WebBaseURL = u
	}

	if dest == "" {
		dest = generator.RootConfigDir
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

	err = generator.WriteConfig(config, path.Join(dest, "config.json"))
	if checkError(err) {
		return err
	}

	return nil
}
