package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"os"
	"path"
	"path/filepath"
)

// UpdateFolderSubdirectories
//
// A function to easily update a folder's subdirectories.
//
// Takes the path of the fotoDen folder.
func UpdateFolderSubdirectories(fpath string) error {
	verbose("Updating folder subdirectories in " + fpath)
	folder := new(generator.Folder)

	err := folder.ReadFolderInfo(path.Join(fpath, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	folder.UpdateSubdirectories(fpath)

	err = folder.WriteFolderInfo(path.Join(fpath, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	return nil
}

// GenerateFolder
//
// Generates an entire fotoDen-compatible folder from
// any images within the current directory,
// including thumbnails, as well as copying over the
// images to a new source folder based on generator.CurrentConfig.
//
// TODO: Change this to accept a mode
// that allows either direct upload of named files by an array,
// to some storage provider,
// or a direct copy to a source folder in the folder structure.
//
// If yes: This will allow fotoDen to become a more 'central' tool
// If no: This allows fotoDen to be a part of a toolset
//
// Takes the folder's name, as well as its path.
// Returns an error if any occur.
func GenerateFolder(name string, fpath string, options GeneratorOptions) error {
	err := os.Mkdir(fpath, 0755)
	if checkError(err) {
		panic(err) // can't continue!
	}

	if *ThumbSrc != "" {
		err = generator.MakeFolderThumbnail(*ThumbSrc, fpath)
		checkError(err)
	}

	folder, err := generator.GenerateFolderInfo(fpath, name)
	if checkError(err) {
		return err
	}

	if options.imagegen == true {
		verbose("Generating album...")
		fileAmount, err := GenerateItems(fpath, options)
		if fileAmount > 0 {
			err = GenerateWeb("album", fpath, folder, options)
			checkError(err)
			folder.FolderType = "album"
			folder.ItemAmount = fileAmount
		} else {
			return fmt.Errorf("Error: No images detected in source! Use -generate folder or a valid source!")
		}
	} else {
		verbose("Generating folder...")
		folder.FolderType = "folder"
		err = GenerateWeb("folder", fpath, folder, options)
		checkError(err)
	}

	err = folder.WriteFolderInfo(path.Join(fpath, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	fpath, _ = filepath.Abs(fpath)
	verbose(path.Dir(fpath))
	verbose(path.Join(path.Dir(fpath), "folderInfo.json"))
	verbose(fmt.Sprint(fileCheck(path.Join(path.Dir(fpath), "folderInfo.json"))))

	if fileCheck(path.Join(path.Dir(fpath), "folderInfo.json")) {
		err = UpdateFolderSubdirectories(path.Dir(fpath))
		checkError(err)
	}

	return nil
}
