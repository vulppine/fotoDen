package tool

import (
	"sort"
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"path"
	"sync"
	"os"
	"io/ioutil"
	"path/filepath"
)

func GenerateItems(fpath string, options *GeneratorOptions) (int, error) {
	verbose("Generating item information to " + fpath)
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

	items, err := generator.GenerateItemInfo(options.source)
	fileArray := items.ItemsInFolder

	if len(fileArray) > 0 {
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
				err = generator.BatchImageConversion(fileArray, "thumb", path.Join(fpath, generator.CurrentConfig.ImageThumbDirectory), generator.ScalingOptions["thumb"])
			}(&waitgroup)
		}

		if options.large == true {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				fmt.Println("Generating display images...")
				err = generator.BatchImageConversion(fileArray, "large", path.Join(fpath, generator.CurrentConfig.ImageLargeDirectory), generator.ScalingOptions["large"])
			}(&waitgroup)
		}

		err = items.WriteItemsInfo(path.Join(fpath, "itemsInfo.json"))
		if checkError(err) {
			return 0, err
		}
	}

	waitgroup.Wait()
	if checkError(err) {
		panic(err)
		// if any errors occur, something wrong happened inbetween all of the batch operations
		// which is really, REALLY bad, considering how it's a bulk copy and conversion at the same time
		// therefore, we need to immediately panic before continuing onwards
	}



	return len(fileArray), nil
}

// UpdateFolderImages
//
// Takes a folder path.
//
// The folder will be updated with all the images in the current folder.
//
// This should, in fact, make a difference between the arrays and copy over any new files...

func UpdateFolderImages(folder string, mode string) error {
	items := new(generator.Items)

	err := items.ReadItemsInfo(path.Join(folder, "items.json"))
	if checkError(err) {
		return err
	}

	dir, err := ioutil.ReadDir(folder)
	if checkError(err) {
		return err
	}

	items.ItemsInFolder = generator.GetArrayOfFiles(dir)
	if mode == "sort" {
		sort.Strings(items.ItemsInFolder)
	}

	return nil
}

// DeleteImage
//
// Deletes an n amount of images from the folder.
//
// DeleteImage goes through the ItemsInFolder array of folderInfo.json,
// and deletes the name of the image from the array,
// and then updates it accordingly.
//
// If the items in the folder are sorted, it uses sort.SearchStrings to find it in O(log n) time.
// Otherwise, it will go through it in O(n) time.

func DeleteImage(folder string, files []string) error {
	items := new(generator.Items)

	err := items.ReadItemsInfo(path.Join(folder, "items.json"))
	if checkError(err) {
		return err
	}

	if sort.StringsAreSorted(items.ItemsInFolder) {
		for i, _ := range files {
			result := sort.SearchStrings(files, files[i])
			if result != len(files) {
				items.ItemsInFolder = append(items.ItemsInFolder[0:i-1], items.ItemsInFolder[i+1:len(items.ItemsInFolder)]...)
			} else if files[result] == files[i] {
				copy(items.ItemsInFolder[0:len(files)-1], items.ItemsInFolder)
			} else {
				fmt.Printf("File %s not found in items, skipping.", files[i])
			}
		}
	} else {
		for i, _ := range files {
			items.ItemsInFolder = generator.RemoveItemFromStringArray(items.ItemsInFolder, files[i])
		}
	}

	err = items.WriteItemsInfo(path.Join(folder, "items.json"))
	if checkError(err) {
		return err
	}

	return nil
}

// InsertImage
//
// Inserts an n amount of images into the folder, at the very end.
//
// Accepts a n amount of file names.
// Generates thumbnails for the given filenames, and copies them over,
// and updates items.json accordingly.

func InsertImage(folder string, files []string, mode string) error {
	items := new(generator.Items)

	err := items.ReadItemsInfo(path.Join(folder, "items.json"))
	if checkError(err) {
		return err
	}

	switch (mode) {
	case "append":
		items.ItemsInFolder = append(items.ItemsInFolder, files...)
	case "sort":
		items.ItemsInFolder = append(items.ItemsInFolder, files...)

		sort.Strings(items.ItemsInFolder)
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
		err = generator.BatchImageConversion(files, "thumb", path.Join(folder, generator.CurrentConfig.ImageThumbDirectory), generator.ScalingOptions["thumb"])
	}(&waitgroup)

	waitgroup.Wait()
	if checkError(err) {
		panic(err)
		// if any errors occur, something wrong happened inbetween all of the batch operations
		// which is really, REALLY bad, considering how it's a bulk copy and conversion at the same time
		// therefore, we need to immediately panic before continuing onwards
	}

	err = items.WriteItemsInfo(path.Join(folder, "folderInfo.json"))
	checkError(err)

	return nil
}
