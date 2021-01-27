package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"path"
)

/*
// CopyWeb
//
// Copies over the HTML pages needed for fotoDen to function.
// These are stored in CurrentConfig.WebSourceLocation. If they don't exist,
// instructions are printed for a 'default' theme to download,
// and how to add themes to fotoDen.
//
// The required theme files are:
// - photo-template.html
// - album-template.html
// - folder-template.html
//
// Takes a string indicating the mode, and a destination folder. Returns an error if any occur.
func CopyWeb(mode string, dest string) error {
	var err error
	inst := "One or more theme files may be missing. You'll need to download a valid fotoDen theme in order to use fotoDen. Once you have it downloaded, insert the files into $CONFIG/fotoDen/web, and try regenerating the folder."
	switch mode {
	case "album":
		err = generator.CopyFile(path.Join(generator.CurrentConfig.WebSourceLocation, "photo.html"), "photo.html", dest)
		err = generator.CopyFile(path.Join(generator.CurrentConfig.WebSourceLocation, "album.html"), "index.html", dest)
	case "folder":
		err = generator.CopyFile(path.Join(generator.CurrentConfig.WebSourceLocation, "folder.html"), "index.html", dest)
	default:
		return fmt.Errorf("A mode was not passed to CopyWeb. Aborting.")
	}

	if checkError(err) {
		return fmt.Errorf(inst)
	}
	return nil
}
*/

func GenerateWeb(m string, dest string, f *generator.Folder, opt GeneratorOptions) error {
	var err error
	var pageOptions *generator.StaticWebVars

	if opt.static || f.IsStatic {
		pageOptions, err = generator.NewStaticWebVars(dest)
		if checkError(err) {
			return err
		}
	} else {
		pageOptions = new(generator.StaticWebVars)
		pageOptions.IsStatic = false
	}

	switch m {
	case "album":
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "album-template.html"), path.Join(dest, "index.html"), pageOptions)
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "photo-template.html"), path.Join(dest, "photo.html"), pageOptions)
	case "folder":
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "folder-template.html"), path.Join(dest, "index.html"), pageOptions)
	default:
		return fmt.Errorf("A mode was not passed to GenerateWeb. Aborting.")
	}

	if checkError(err) {
		return err
	}

	return nil
}

func UpdateWeb(folder string) error {
	f := new(generator.Folder)

	err := f.ReadFolderInfo(path.Join(folder, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	err = GenerateWeb(f.FolderType, folder, f, genoptions)
	if checkError(err) {
		return err
	}

	return nil
}
