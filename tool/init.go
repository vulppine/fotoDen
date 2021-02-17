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

func CopyThemeToConfig(srcpath string) error {
	wd, _ := os.Getwd()
	srcpath, err := filepath.Abs(srcpath)
	t, err := ReadThemeConfig(path.Join(srcpath, "theme.json"))
	tpath, err := filepath.Abs(path.Join(generator.RootConfigDir, "theme", t.ThemeName))
	err = os.MkdirAll(tpath, 0755)
	if checkError(err) {
		return err
	}

	err = os.Mkdir(filepath.Join(tpath, "html"), 0755)
	if checkError(err) {
		return err
	}

	err = generator.CopyFile(path.Join(srcpath, "theme.json"), "theme.json", tpath)
	err = generator.CopyFile(
		path.Join(srcpath, "html", "photo-template.html"),
		"photo-template.html",
		path.Join(tpath, "html"))
	err = generator.CopyFile(
		path.Join(srcpath, "html", "album-template.html"),
		"album-template.html",
		path.Join(tpath, "html"))
	err = generator.CopyFile(
		path.Join(srcpath, "html", "folder-template.html"),
		"folder-template.html",
		path.Join(tpath, "html"))
	if checkError(err) {
		return err
	}

	if len(t.Stylesheets) != 0 {
		err = os.Mkdir(
			path.Join(tpath, "css"), 0755)
		if checkError(err) {
			return err
		}

		os.Chdir(path.Join(srcpath, "css"))
		err = generator.BatchCopyFile(
			t.Stylesheets,
			path.Join(tpath, "css"))
		checkError(err)
		os.Chdir(wd)
	}

	if len(t.Scripts) != 0 {
		err = os.Mkdir(
			path.Join(tpath, "js"), 0755)
		if checkError(err) {
			return err
		}

		os.Chdir(path.Join(srcpath, "js"))
		err = generator.BatchCopyFile(
			t.Scripts, path.Join(tpath, "js"))
		checkError(err)
		os.Chdir(wd)
	}

	if len(t.Other) != 0 {
		err = os.Mkdir(
			path.Join(tpath, "other"), 0755)
		if checkError(err) {
			return err
		}

		os.Chdir(path.Join(srcpath, "etc"))
		err = generator.BatchCopyFile(
			t.Other, path.Join(tpath, "etc"))
		checkError(err)
		os.Chdir(wd)
	}

	if !fileCheck(path.Join(generator.RootConfigDir, "defaulttheme")) {
		f, err := os.Create(path.Join(generator.RootConfigDir, "defaulttheme"))
		if checkError(err) {
			return err
		}
		f.WriteString(t.ThemeName)
		f.Close()
	}

	return nil
}

// InitializeWebTheme sets up web templates according to a given URL, and path containing templates.
// All templates should be labelled with [photo, album, folder]-template.html.
func InitializeWebTheme(u string, srcpath string, dest string) error {
	wd, _ := os.Getwd()

	t, err := ReadThemeConfig(path.Join(srcpath, "theme.json"))
	webvars, err := generator.NewWebVars(u)
	checkError(err)

	tpath, err := filepath.Abs(path.Join(dest, "theme", t.ThemeName))
	err = os.MkdirAll(tpath, 0755)
	if checkError(err) {
		return err
	}

	err = generator.CopyFile(path.Join(srcpath, "theme.json"), "theme.json", tpath)

	err = os.Mkdir(
		path.Join(tpath, "html"), 0755)
	if checkError(err) {
		return err
	}

	err = generator.ConfigureWebFile(
		path.Join(srcpath, "html", "photo-template.html"),
		path.Join(
			tpath,
			"html",
			"photo-template.html"),
		webvars)
	checkError(err)

	err = generator.ConfigureWebFile(
		path.Join(srcpath, "html", "album-template.html"),
		path.Join(
			tpath,
			"html",
			"album-template.html"),
		webvars)
	checkError(err)

	err = generator.ConfigureWebFile(
		path.Join(srcpath, "html", "folder-template.html"),
		path.Join(
			tpath,
			"html",
			"folder-template.html"),
		webvars)
	checkError(err)

	if len(t.Stylesheets) != 0 {
		err = os.Mkdir(
			path.Join(tpath, "css"), 0755)
		if checkError(err) {
			return err
		}

		os.Chdir(path.Join(srcpath, "css"))
		err = generator.BatchCopyFile(
			t.Stylesheets,
			path.Join(tpath, "css"))
		checkError(err)
		os.Chdir(wd)
	}

	if len(t.Scripts) != 0 {
		err = os.Mkdir(
			path.Join(tpath, "js"), 0755)
		if checkError(err) {
			return err
		}

		os.Chdir(path.Join(srcpath, "js"))
		err = generator.BatchCopyFile(
			t.Scripts, path.Join(tpath, "js"))
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

	err := generator.CopyFile(path.Join(fpath), "fotoDen.js", generator.RootConfigDir)
	if checkError(err) {
		return err
	}

	return nil
}

// WebsiteConfig represents a struct that contains
// everything needed for the fotoDen generator to work.
type WebsiteConfig struct {
	Name            string
	RootLocation    string
	Theme           string
	URL             string
	GeneratorConfig generator.Config
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
//
// TODO: Break this apart into smaller chunks, holy fuck
func InitializefotoDenRoot(rootpath string, webconfig WebsiteConfig) error {
	rootpath, _ = filepath.Abs(rootpath)

	err := generator.GenerateWebRoot(rootpath)
	if checkError(err) {
		panic(err)
	}


	var w *generator.WebConfig

	if webconfig.Theme == "" {
		if fileCheck(path.Join(generator.RootConfigDir, "defaulttheme")) {
			f, err := os.Open(path.Join(generator.RootConfigDir, "defaulttheme"))
			d, err := ioutil.ReadAll(f)
			if checkError(err) {
				return err
			}

			webconfig.Theme = string(d)
		} else {
			return fmt.Errorf("theme does not exist, cannot continue generation")
		}
	}

	if WizardFlag == true {
		webconfig, w = setupWebsite(rootpath)
		generator.CurrentConfig = webconfig.GeneratorConfig
	} else {
		webconfig.RootLocation = rootpath
		webconfig.GeneratorConfig = generator.DefaultConfig
		webconfig.GeneratorConfig.WebBaseURL = webconfig.URL
		if webconfig.Name == "" {
			verbose("Name not specified, using base of given path.")
			webconfig.Name = filepath.Base(rootpath)
			verbose(webconfig.Name)
		}
		webconfig.GeneratorConfig.WebSourceLocation, err = filepath.Abs(
			path.Join(
				generator.RootConfigDir,
				"sites",
				webconfig.Name,
				"theme",
				webconfig.Theme,
			))
		generator.CurrentConfig = webconfig.GeneratorConfig
		w = generator.GenerateWebConfig(URLFlag)
		w.Theme = true // we're generating this from fotoDen tool, so we're using a theme obviously
		w.WebsiteTitle = webconfig.Name
		w.ImageRootDir = webconfig.GeneratorConfig.ImageRootDirectory
		w.ThumbnailFrom = w.ImageSizes[0].SizeName
		w.DisplayImageFrom = w.ImageSizes[len(w.ImageSizes)-1].SizeName
		for _, v := range w.ImageSizes {
			w.DownloadSizes = append(w.DownloadSizes, v.SizeName)
		}

		if URLFlag == "" {
			fmt.Printf("You will have to configure your photo storage provider in %v.", path.Join(rootpath, "config.json"))
		}
	}

	spath := path.Join(generator.RootConfigDir, "sites", webconfig.Name)
	verbose("Creating site config directory at: " + spath)
	err = os.MkdirAll(spath, 0755)
	if checkError(err) {
		return err
	}

	err = generator.WriteJSON(path.Join(spath, "config.json"), "multi", webconfig)
	if checkError(err) {
		return err
	}

	generator.CurrentConfig = webconfig.GeneratorConfig

	err = InitializeWebTheme(
		webconfig.URL,
		path.Join(generator.RootConfigDir, "theme", webconfig.Theme),
		path.Join(spath),
	)

	err = w.WriteWebConfig(path.Join(rootpath, "config.json"))
	if checkError(err) {
		return err
	}

	err = generator.CopyFile(path.Join(generator.RootConfigDir, "fotoDen.js"), "fotoDen.js", path.Join(rootpath, "js"))
	checkError(err)

	tpath := path.Join(spath, "theme", webconfig.Theme)
	t, err := ReadThemeConfig(path.Join(tpath, "theme.json"))
	if checkError(err) {
		return err
	}

	wd, _ := os.Getwd()

	os.Chdir(tpath)
	if len(t.Stylesheets) != 0 {
		os.Chdir("css")
		err = generator.BatchCopyFile(t.Stylesheets, path.Join(rootpath, "theme", "css"))
		if checkError(err) {
			return err
		}
		os.Chdir(tpath)
	}

	if len(t.Scripts) != 0 {
		os.Chdir("js")
		err = generator.BatchCopyFile(t.Scripts, path.Join(rootpath, "theme", "js"))
		if checkError(err) {
			return err
		}
		os.Chdir(tpath)
	}

	if len(t.Other) != 0 {
		os.Chdir("etc")
		err = generator.BatchCopyFile(t.Other, path.Join(rootpath, "theme", "etc"))
		if checkError(err) {
			return err
		}
	}

	os.Chdir(wd)

	folder, err := generator.GenerateFolderInfo(rootpath, w.WebsiteTitle) // do it in rootpath since we're not trying to scan for images in the current folder
	folder.Type = "folder"
	checkError(err)
	err = folder.WriteFolderInfo(path.Join(rootpath, "folderInfo.json"))
	checkError(err)

	err = GenerateWeb("folder", rootpath, folder, Genoptions)
	checkError(err)

	if !fileCheck(path.Join(generator.RootConfigDir, "defaultsite")) {
		f, err := os.Create(path.Join(generator.RootConfigDir, "defaultsite"))
		if checkError(err) {
			return err
		}
		f.WriteString(webconfig.Name)
		f.Close()
	}

	return nil
}

/* Deprecated since fotoDen configs are now included at a per-site basis
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
*/
