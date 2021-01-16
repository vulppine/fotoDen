package generator

import (
	"fmt"
	"os"
	"strings"
	"path"
	"github.com/h2non/bimg"
)

// BatchOperationOnFiles
//
// Takes two arguments, an array of file names, and a function that takes a string and an int.
//
// The string will be the file name, while the int will be the index of the file in the array.
//
// The function will iterate over every file name until the end of the array is reached,
// passing the file name into the function.

func BatchOperationOnFiles(files []string, fn func(string, int) error) error {
	for i := 0; i < len(files); i++ {
		err := fn(files[i], i)
		if err != nil {
			fmt.Println("An error occurred during operation: ", err)
		}
	}

	return nil
}

// BatchCopyFileByIndex
//
// Copies a list of file string names to the current WorkingDirectory by index.
// Returns an error if one occurs, otherwise nil.
// Also preserves the current extension of the file. (This is due to a NeoCities Free restriction)

func BatchCopyFile(files []string, directory string) error {
	wd, _ := os.Getwd()
	verbose("Attempting a batch copy from " + wd + " to " + directory)
	batchCopyFile := func (file string, index int) error {
		err := CopyFile(file, file, directory)
		if err != nil {
			return err
		}

		return nil
	}

	err := BatchOperationOnFiles(files, batchCopyFile)
	if err != nil {
		return err
	}

	return nil
}

// BatchResizeImages
//
// Resizes a set of images to thumbnail size and puts them into the given directory, as according to CurrentConfig.
// Returns an error if one occurs, otherwise nil.

func BatchImageConversion(files []string, prefix string, directory string, ScalingOptions ImageScale) error {
	wd, _ := os.Getwd()
	verbose("Generating thumbnails in " + wd + " and placing them in " + directory)
	batchResizeImage := func (file string, index int) error {
		err := ResizeImage(file, prefix + "_" + strings.Split(path.Base(file), ".")[0] + ".jpg", ScalingOptions, directory, bimg.JPEG)
		if err != nil && err != fmt.Errorf("skip") {
			return err
		}

		return nil
	}

	err := BatchOperationOnFiles(files, batchResizeImage)
	if err != nil {
		return err

	}

	return nil
}
