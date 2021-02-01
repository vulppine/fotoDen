package generator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Structs for fotoDen JSON files
//
// Structs for usage with fotoDen JSON files.

// Folder represents a folderInfo.json file used by fotoDen.
// It has all the needed values for fotoDen.js to operate correctly.
// fotoDen/generator does not provide functions to manage this - only to read and create these.
type Folder struct {
	FolderName          string // The name of the folder.
	FolderDesc          string // The description of a folder.
	FolderAuthor        string // The author of a folder.
	FolderShortName     string // The shortname of the folder (can be taken from the filesystem folder name)
	FolderType          string // The type of folder (currently supports only album or folder)
	FolderThumbnail     bool   // If a thumbnail exists or not. This is dictated by the generation of thumb.jpg.
	ItemAmount          int
	SubfolderShortNames []string // Any folders that are within the folder (updated whenever the generator is called in the folder)
	IsStatic            bool     // If the folder was generated statically, or has information inserted dynamically.
}

// Items represents an itemsInfo.json file used by fotoDen.
// It is used mainly in album-type folders, and contains a bool indicating whether
// metadata is being used, and a string array (potentially large) of file names.
type Items struct {
	Metadata bool // Dictates whether or not each image has its own ImageMeta object.
	// If this is false, then no metadata will be read.
	ItemsInFolder []string // All the items in a folder, by name, in an array.
}

// GeneratorConfig represents the configuration for fotoDen's generator, and where
// images will go, as well as what sizes will be generated.
// ImageSizes is a map with string keys containing ImageScale structs, which dictate
// how images will be resized.
type GeneratorConfig struct {
	ImageRootDirectory string // where all images are stored (default: img)
	ImageSrcDirectory  string // where all source images are stored (default: ImageRootDirectory/src)
	ImageMetaDirectory string // where all meta files per image are stored (default: ImageRootDirectory/meta)
	ImageSizes         map[string]ImageScale
	WebSourceLocation  string // where all html/css/js files are stored for fotoDen's functionality
	WebBaseURL         string // what the base URL is (aka, fotoDen's location)
}

// some defaults in case we never have a fotoDen config file opened

var userConfigDir, _ = os.UserConfigDir()

// FotoDenConfigDir is where the configuration files are stored.
var FotoDenConfigDir = path.Join(userConfigDir, "fotoDen")

// DefaultConfig contains a template for fotoDen to use.
// TODO: Move this to some kind of GeneratorConfig generator.
var DefaultConfig GeneratorConfig = GeneratorConfig{
	ImageRootDirectory: "img",
	ImageMetaDirectory: "meta",
	ImageSizes: map[string]ImageScale{
		"small": ImageScale{ScalePercent: 0.25},
		"medium": ImageScale{ScalePercent: 0.5},
		"large": ImageScale{ScalePercent: 0.75},
	},
	ImageSrcDirectory: "src",
	WebSourceLocation: path.Join(FotoDenConfigDir, "web"), // remember when $HOME webpage folders were a thing?
	WebBaseURL:        "",                                 // this should be set during configuration generation
}

// CurrentConfig represents the current generator config, and can be used as reference
// for any package that calls fotoDen/generator.
var CurrentConfig GeneratorConfig

// WorkingDirectory is the current working directory that fotoDen was started in.
var WorkingDirectory, _ = os.Getwd()

// Verbose is used toggle verbose statements - when toggled, prints what the generator is doing to console.
var Verbose bool // if this is set, everything important is printed

func verbose(print string) {
	if Verbose {
		log.Println(print)
	}
}

// WriteJSON writes a struct as a JSON file to a specified pathname.
// Takes a filepath, a "mode", and an interface containing something that translates to valid JSON according to encoding/json.
// Mode toggles between non-indented JSON, and indented JSON.
// Returns an error if any occur.
func WriteJSON(filePath string, mode string, iface interface{}) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	var toWrite []byte

	switch mode {
	case "single":
		toWrite, err = json.Marshal(iface)
		if err != nil {
			return err
		}
	case "multi":
		toWrite, err = json.MarshalIndent(iface, "", "\t")
		if err != nil {
			return err
		}
	}

	_, err = file.Write(toWrite)
	if err != nil {
		return err
	}

	return nil
}

// ReadJSON reads a JSON file from a pathname, and puts it into the specified interface.
// Returns an error if any occur.
func ReadJSON(filePath string, iface interface{}) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, iface)
	if err != nil {
		return err
	}

	return nil
}

// OpenfotoDenConfig sets the current fotoDen generator configuration to this
func OpenfotoDenConfig(configLocation string) error {
	err := ReadJSON(configLocation, &CurrentConfig)
	if err != nil {
		return err
	}

	return nil
}

// WritefotoDenConfig attempts to write CurrentConfig to a new file at configLocation.
func WritefotoDenConfig(config GeneratorConfig, configLocation string) error {
	err := WriteJSON(configLocation, "multi", config)
	if err != nil {
		return err
	}

	return nil
}

// RemoveItemFromStringArray removes an item from a string array at O(n) speed.
func RemoveItemFromStringArray(array []string, item string) []string {
	verbose("Attempting to remove " + item + " from an array.")
	newArray := make([]string, 0)

	for i := 0; i < len(array); i++ {
		if array[i] != item {
			newArray = append(newArray, array[i])
		}
	}

	return newArray
}
