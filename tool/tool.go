package tool

import (
	"fmt"
	"os"
	"io/ioutil"

	"github.com/vulppine/fotoDen/generator"
)

// fotoDen tool:
//
// Compared to the generator, this contains all the functions
// that combine all file generation-related functions
// into something more callable,
// while equally separating more precise filesystem/folder functions
// away from the actual generation API.
//
// i.e., rather than calling all the functions to update a folder's information,
// you can just call tool.UpdateFolder(pathname).
//
// for another example, this allows for the insertion/deletion of images

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
	source   string
	copy     bool
	gensizes bool
	imagegen bool
	sort     bool
	meta     bool
	static   bool
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println("An error occured during operation: ", err)
		return true
	}

	return false
}

var Verbose bool

func verbose(print string) {
	if Verbose {
		fmt.Println(print)
	}
}

func fileCheck(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type fvisitFunction func(string) error

// RecursiveVisit
//
// Recursively visits folders, and performs a function
// inside of them. To ensure safety, this only works
// with fotoDen folders.
//
// It detects if a folder is a fotoDen folder in a lazy way,
// by seeing if a folder contains a folderInfo.json.
// If it does not, it terminates
func RecursiveVisit(folder string, fn fvisitFunction) error {
	wd, err := os.Getwd()
	if checkError(err) {
		return err
	}

	defer os.Chdir(wd)

	err = os.Chdir(folder)
	if checkError(err) {
		return err
	}

	if !fileCheck("folderInfo.json") {
		return nil
	}

	err = fn(folder)
	if checkError(err) {
		return err
	}

	folders, err := ioutil.ReadDir(folder)
	if checkError(err) {
		return err
	}

	for _, folder := range generator.GetArrayOfFolders(folders) {
		err = RecursiveVisit(folder, fn)
		if checkError(err) {
			return err
		}
	}

	return nil
}
