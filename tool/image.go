package tool

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"

	"github.com/h2non/bimg"
	"github.com/vulppine/fotoDen/generator"
)


// GenerateItems generates the album part of a fotoDen folder in fpath.
//
// It takes an options struct (which is just a set of options condensed into one struct)
// and generates an image folder based on the options given (e.g., source, etc.)
func GenerateItems(fpath string, options GeneratorOptions) (int, error) {
	verbose("GenerateItems: Current generator options: " + fmt.Sprint(options))
	verbose("Generating item information to " + fpath)
	var waitgroup sync.WaitGroup

	items, err := generator.GenerateItemInfo(options.Source)
	verbose("Current images in folder: " + fmt.Sprint(items.ItemsInFolder))

	if len(items.ItemsInFolder) > 0 {
		if options.Sort == true {
			sort.Strings(items.ItemsInFolder)
		}
		err = generator.MakeAlbumDirectoryStructure(fpath)
		if checkError(err) {
			panic(err)
		}

		if options.Source != "" {
			verbose("Changing to directory: " + options.Source)
			wd, err := os.Getwd()
			checkError(err)
			fpath, err = filepath.Abs(fpath)
			checkError(err)
			source, err := filepath.Abs(options.Source)
			checkError(err)

			defer os.Chdir(wd)
			os.Chdir(source)
			verbose("Current directory: " + func() string {
				dir, _ := os.Getwd()
				return dir
			}())
		}

		if options.Copy == true {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				fmt.Println("Copying files...")
				err = generator.BatchCopyFile(items.ItemsInFolder, path.Join(fpath, generator.CurrentConfig.ImageRootDirectory, generator.CurrentConfig.ImageSrcDirectory))
			}(&waitgroup)
		}

		if options.Gensizes == true {
			verbose("Attempting to generate from sizes: " + fmt.Sprint(generator.CurrentConfig.ImageSizes))
			for k, v := range generator.CurrentConfig.ImageSizes {
				sizeName := k
				sizeOpts := v
				waitgroup.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					fmt.Printf("Generating size %s...\n", sizeName)
					err = generator.BatchImageConversion(items.ItemsInFolder, sizeName, path.Join(fpath, generator.CurrentConfig.ImageRootDirectory, sizeName), sizeOpts)
				}(&waitgroup)
			}
		}

		if options.Meta == true {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				verbose("Generating metadata to: " + path.Join(fpath, generator.CurrentConfig.ImageRootDirectory, generator.CurrentConfig.ImageMetaDirectory))
				err = generator.BatchImageMeta(items.ItemsInFolder, path.Join(fpath, generator.CurrentConfig.ImageRootDirectory, generator.CurrentConfig.ImageMetaDirectory))
			}(&waitgroup)
			items.Metadata = true
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

	return len(items.ItemsInFolder), nil
}

// UpdateImages updates all the images in a fotoDen folder.
func UpdateImages(folder string, options GeneratorOptions) error {
	items := new(generator.Items)

	err := items.ReadItemsInfo(path.Join(folder, "itemsInfo.json"))
	if checkError(err) {
		return err
	}

	dir, err := ioutil.ReadDir(options.Source)
	if checkError(err) {
		return err
	}

	if options.Source != "" {
		verbose("Changing to directory: " + options.Source)
		wd, err := os.Getwd()
		checkError(err)
		folder, err = filepath.Abs(folder)
		checkError(err)
		source, err := filepath.Abs(options.Source)
		checkError(err)

		defer os.Chdir(wd)
		os.Chdir(source)
		verbose("Current directory: " + func() string {
			dir, _ := os.Getwd()
			return dir
		}())
	}

	items.ItemsInFolder = generator.IsolateImages(generator.GetArrayOfFiles(dir))
	if options.Sort {
		sort.Strings(items.ItemsInFolder)
	}

	err = items.WriteItemsInfo(path.Join(folder, "itemsInfo.json"))
	if checkError(err) {
		return err
	}

	return nil
}

// DeleteImage deletes an image from the folder.
//
// DeleteImage goes through the ItemsInFolder array of folderInfo.json,
// and deletes the name of the image from the array,
// and then updates it accordingly.
//
// If the items in the folder are sorted, it uses sort.SearchStrings to find it in O(log n) time.
// Otherwise, it will go through it in O(n) time.
func DeleteImage(folder string, files ...string) error {
	items := new(generator.Items)

	err := items.ReadItemsInfo(path.Join(folder, "itemsInfo.json"))
	if checkError(err) {
		return err
	}

	for _, file := range files {
		if sort.StringsAreSorted(items.ItemsInFolder) {
			result := sort.SearchStrings(items.ItemsInFolder, file)
			if result == len(items.ItemsInFolder) {
				copy(items.ItemsInFolder[0:result-1], items.ItemsInFolder)
			} else if result == 0 {
				copy(items.ItemsInFolder[1:len(items.ItemsInFolder)], items.ItemsInFolder)
			} else if result != len(items.ItemsInFolder) && items.ItemsInFolder[result] == file {
				items.ItemsInFolder = append(items.ItemsInFolder[0:result], items.ItemsInFolder[result+1:len(items.ItemsInFolder)]...)
			} else {
				log.Printf("File %s not found in items.", file)
			}
		} else {
			items.ItemsInFolder = generator.RemoveItemFromStringArray(items.ItemsInFolder, file)
		}
	}

	err = items.WriteItemsInfo(path.Join(folder, "itemsInfo.json"))
	if checkError(err) {
		return err
	}

	return nil
}

// InsertImage inserts an image into a fotoDen folder.
func InsertImage(folder string, mode string, options GeneratorOptions, files ...string) error {
	items := new(generator.Items)

	err := items.ReadItemsInfo(path.Join(folder, "itemsInfo.json"))
	if checkError(err) {
		return err
	}

	var waitgroup sync.WaitGroup

	folder, err = filepath.Abs(folder)
	checkError(err)
	wd, err := os.Getwd()
	checkError(err)

	for _, fv := range files {
		fi, err := filepath.Abs(fv)
		if checkError(err) {
			return err
		}

		if filepath.Dir(fi) != wd {
			verbose("Changing to directory: " + filepath.Dir(fi))
			err := os.Chdir(filepath.Dir(fi))
			if checkError(err) {
				return err
			}
		}

		f := filepath.Base(fv)
		verbose("Current file: " + f)
		switch mode {
		case "append":
			items.ItemsInFolder = append(items.ItemsInFolder, f)
		case "sort":
			items.ItemsInFolder = append(items.ItemsInFolder, f)

			sort.Strings(items.ItemsInFolder)
		}

		verbose("Current directory: " + func() string {
			dir, _ := os.Getwd()
			return dir
		}())

		if options.Copy {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				fmt.Println("Copying file...")
				err = generator.CopyFile(
					f, f,
					path.Join(
						folder,
						generator.CurrentConfig.ImageRootDirectory,
						generator.CurrentConfig.ImageSrcDirectory,
					),
				)
			}(&waitgroup)
		}

		if options.Gensizes {
			for k, v := range generator.CurrentConfig.ImageSizes {
				sizeName := k
				sizeOpts := v
				waitgroup.Add(1)
				go func(wg *sync.WaitGroup) {
					defer wg.Done()
					fmt.Printf("Generating size %s...\n", sizeName)
					err = generator.ResizeImage(
						f, sizeName + "_" + f, sizeOpts,
						path.Join(
							folder,
							generator.CurrentConfig.ImageRootDirectory,
							sizeName,
						),
						bimg.JPEG,
					)
				}(&waitgroup)
			}
		}

		if options.Meta == true {
			waitgroup.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				verbose("Generating metadata to: " + generator.CurrentConfig.ImageMetaDirectory)
				m := new(generator.ImageMeta)
				err = m.WriteImageMeta(
					filepath.Join(
						folder,
						generator.CurrentConfig.ImageRootDirectory,
						generator.CurrentConfig.ImageMetaDirectory,
					),
					f,
				)
			}(&waitgroup)
		}

		os.Chdir(wd)
		if checkError(err) {
			return err
		}
	}

	waitgroup.Wait()
	if checkError(err) {
		panic(err)
		// if any errors occur, something wrong happened inbetween all of the batch operations
		// which is really, REALLY bad, considering how it's a bulk copy and conversion at the same time
		// therefore, we need to immediately panic before continuing onwards
	}

	err = items.WriteItemsInfo(path.Join(folder, "itemsInfo.json"))
	checkError(err)

	return nil
}
