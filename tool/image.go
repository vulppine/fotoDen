package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"path"
	"sync"
)

// UpdateFolderImages
//
// Takes a folder path.
//
// The folder will be updated with all the images in the current folder.
//
// This should, in fact, make a difference between the arrays and copy over any new files...

// DeleteImage
//
// Deletes an n amount of images from the folder.
//
// DeleteImage goes through the ItemsInFolder array of folderInfo.json,
// and deletes the name of the image from the array,
// and then updates it accordingly.
//
// TODO: Similar to GenerateFolder, give this a mode
// to allow switching between working with files,
// and working with cloud objects

func DeleteImage(files ...string) {

}

// InsertImage
//
// Inserts an n amount of images into the folder, at the very end.
//
// Accepts a n amount of file names.
// Generates thumbnails for the given filenames, and copies them over,
// and updates folderInfo.json accordingly.

func InsertImage(folder string, files []string) error {
	folderInfo := new(generator.Folder)

	err := folderInfo.ReadFolderInfo(path.Join(folder, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	for file := range files {
		folderInfo.ItemsInFolder = append(folderInfo.ItemsInFolder, files[file])
	}

	var waitgroup sync.WaitGroup

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("Copying files...")
		err = generator.BatchCopyFile(files, path.Join(folder, generator.CurrentConfig.ImageSrcDirectory))
	}(&waitgroup)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("Generating thumbnails...")
		err = generator.BatchImageConversion(files, "thumb", path.Join(folder, generator.CurrentConfig.ImageThumbDirectory), generator.ThumbScalingOptions)
	}(&waitgroup)

	waitgroup.Wait()
	if checkError(err) {
		panic(err)
		// if any errors occur, something wrong happened inbetween all of the batch operations
		// which is really, REALLY bad, considering how it's a bulk copy and conversion at the same time
		// therefore, we need to immediately panic before continuing onwards
	}

	err = folderInfo.WriteFolderInfo(path.Join(folder, "folderInfo.json"))
	checkError(err)

	return nil
}
