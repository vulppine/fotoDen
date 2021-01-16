package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"os"
	"path"
	"path/filepath"
	"sync"
)

// GenerationOptions
//
// Some options for the generator.
// Includes:
// - source
// - copy
// - thumb
// - large
//
// from the flags.

type GeneratorOptions struct {
	source string
	copy   bool
	thumb  bool
	large  bool
}

// GenerateFolderStructure
//
// Generates a fotoDen-compatible folder structure,
// without copying over any images.
//
// It will be up to the end user to update the folder information file.
//
// Takes the folder's name, as well as its path.
// Returns an error if any occur.

func GenerateFolderStructure(name string, fpath string) error {

	verbose("Generating folder structure...")
	verbose("Making folder " + fpath)
	err := os.Mkdir(fpath, 0755)
	if checkError(err) {
		panic(err) // can't continue!
	}

	err = generator.MakeAlbumDirectoryStructure(fpath)
	if checkError(err) {
		panic(err)
	}

	newFolder, err := generator.GenerateFolderInfo(fpath, name)
	if checkError(err) {
		return err
	}

	err = CopyWeb("folder", fpath) // it's automatically a folder because of the 0 image amount
	checkError(err)

	err = newFolder.WriteFolderInfo(path.Join(fpath, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	return nil
}

// UpdateFolderSubdirectories
//
// A function to easily update a folder's subdirectories.
//
// Takes the path of the fotoDen folder.

func UpdateFolderSubdirectories(fpath string) error {
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
// Takes the folder's name, as well as its path.
// Returns an error if any occur.

func GenerateFolder(name string, fpath string, options *GeneratorOptions) error {

	err := os.Mkdir(fpath, 0755)
	if checkError(err) {
		panic(err) // can't continue!
	}

	err = generator.MakeAlbumDirectoryStructure(fpath)
	if checkError(err) {
		panic(err)
	}

	var waitgroup sync.WaitGroup

	if options.source != "" {
		verbose("Changing to directory: " + options.source)
		wd, err := os.Getwd()
		checkError(err)
		fpath, err = filepath.Abs(fpath)
		checkError(err)
		source, err := filepath.Abs(options.source)
		checkError(err)

		defer os.Chdir(wd)
		os.Chdir(source)
		verbose("Current directory: " + source)
	}

	dir, err := os.Open("./")
	defer dir.Close()
	if checkError(err) {
		panic(err)
	}

	verbose("Reading items in folder: " + options.source)
	dirContents, err := dir.Readdir(0)
	if checkError(err) {
		return err
	}
	fileArray := generator.IsolateImages(generator.GetArrayOfFiles(dirContents))

	if *ThumbSrc != "" {
		err = generator.MakeFolderThumbnail(*ThumbSrc, fpath)
		checkError(err)
	}

	folder, err := generator.GenerateFolderInfo(fpath, name)
	if checkError(err) {
		return err
	}

	if len(fileArray) > 0 {
		folder.FolderType = "album"
		err = CopyWeb("album", fpath)
		checkError(err)

		if options.copy == true {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				fmt.Println("Copying files...")
				err = generator.BatchCopyFile(fileArray, path.Join(fpath, generator.CurrentConfig.ImageSrcDirectory))
			}(&waitgroup)
		}

		if options.thumb == true {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				fmt.Println("Generating thumbnails...")
				err = generator.BatchImageConversion(fileArray, "thumb", path.Join(fpath, generator.CurrentConfig.ImageThumbDirectory), generator.ThumbScalingOptions)
			}(&waitgroup)
		}

		if options.large == true {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				fmt.Println("Generating display images...")
				err = generator.BatchImageConversion(fileArray, "large", path.Join(fpath, generator.CurrentConfig.ImageLargeDirectory), generator.LargeScalingOptions)
			}(&waitgroup)
		}

		folder.ItemsInFolder = fileArray
	} else {
		folder.FolderType = "folder"
		err = CopyWeb("folder", fpath)
		checkError(err)
	}

	waitgroup.Wait()
	if checkError(err) {
		panic(err)
		// if any errors occur, something wrong happened inbetween all of the batch operations
		// which is really, REALLY bad, considering how it's a bulk copy and conversion at the same time
		// therefore, we need to immediately panic before continuing onwards
	}

	err = folder.WriteFolderInfo(path.Join(fpath, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	return nil
}
