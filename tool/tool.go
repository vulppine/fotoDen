package tool

import (
	"fmt"
	"os"
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
	return os.IsExist(err)
}
