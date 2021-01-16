package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"path"
)

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
