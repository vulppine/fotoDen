// fotoDen tool:
//
// Reference implementation for generation of fotoDen-based websites.
package tool

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/vulppine/fotoDen/generator"
)

// WizardFlag specifies if fotoDen tool functions should have interactive input or not.
var WizardFlag bool

// NameFlag sets the name for a folder/album. If this is not set, fotoDen will automatically use the folder's name.
var NameFlag string

// Recurse toggles recursive functions on directories. This is primarily used for the update command.
var Recurse bool

// URLFlag sets the URL for functions that require a URL. This is mostly used in initialization.
var URLFlag string

// GeneratorOptions is a set of options for the generator.
//
// Includes:
// - source
// - copy
// - thumb
// - large
//
// from the flags.
type GeneratorOptions struct {
	Source   string
	Copy     bool
	Gensizes bool
	ImageGen bool
	Sort     bool
	Meta     bool
	Static   bool
}

// Genoptions is a global variable for functions that use GeneratorOptions.
var Genoptions GeneratorOptions

func checkError(err error) bool {
	if err != nil {
		log.Println("An error occured during operation: ", err)
		return true
	}

	return false
}

// Verbose toggles the verbosity of the command line tool.
var Verbose bool

func verbose(print string) {
	if Verbose {
		log.Println(print)
	}
}

func fileCheck(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type fvisitFunction func(string) error

// RecursiveVisit recursively visits folders, and performs a function
// inside of them. To ensure safety, this only works
// with fotoDen folders.
//
// It detects if a folder is a fotoDen folder in a lazy way,
// by seeing if a folder contains a folderInfo.json.
// If it does not, it terminates
//
// TODO: Replace this with the new fs library function in Go 1.16
func RecursiveVisit(folder string, fn fvisitFunction) error {
	wd, err := os.Getwd()
	if checkError(err) {
		return err
	}

	if !fileCheck(filepath.Join(folder, "folderInfo.json")) {
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

	defer os.Chdir(wd)

	err = os.Chdir(folder)
	if checkError(err) {
		return err
	}

	for _, f := range generator.GetArrayOfFolders(folders) {
		err = RecursiveVisit(f, fn)
		if checkError(err) {
			return err
		}
	}

	return nil
}
