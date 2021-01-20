package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Structs for fotoDen JSON files
//
// Structs for usage with fotoDen JSON files.

/* -- Completely deprecated, refer to Folder for this information
type Album struct {
	albumName string		// The name of the album.
	albumShortName string	// The album's shortname (can be taken from the folder name)
	supAlbumName string		// The name of the album above it (for presentation purposes)
	photoAmount int			// The amount of photos in the album.
}
*/

type Folder struct {
	FolderName          string   // The name of the folder.
	FolderShortName     string   // The shortname of the folder (can be taken from the filesystem folder name)
	FolderType          string   // The type of folder (currently supports only album or folder)
	FolderThumbnail     string   // The path to the thumbnail, relative to the current folder.
	SubfolderShortNames []string // Any folders that are within the folder (updated whenever the generator is called in the folder)
}

type Items struct {
	Metadata			bool	 // Dictates whether or not each image has its own ImageMeta object.
								 // If this is false, then no metadata will be read.
	ItemsInFolder       []string // All the items in a folder, by name, in an array.
}

// TODO: Implement this!

type ImageMeta struct {
	ImageName			string   // The name of an image.
	ImageDesc			string   // The description of an image.
}

type GeneratorConfig struct {
	ImageRootDirectory  string // where all images are stored (default: img)
	ImageSrcDirectory   string // where all source images are stored (default: ImageRootDirectory/src)
	ImageMetaDirectory  string // where all meta files per image are stored (default: ImageRootDirectory/meta)
	ImageSizes			map[string]ImageScale
	WebSourceLocation   string // where all html/css/js files are stored for fotoDen's functionality
	WebBaseURL          string // what the base URL is (aka, fotoDen's location)
}

// some defaults in case we never have a fotoDen config file opened

var userConfigDir, _ = os.UserConfigDir()
var FotoDenConfigDir = path.Join(userConfigDir, "fotoDen")
var DefaultConfig GeneratorConfig = GeneratorConfig{
	ImageRootDirectory:  "img",
	ImageMetaDirectory:  "meta",
	ImageSizes: map[string]ImageScale{
		"thumb": ImageScale{MaxHeight: 800},
		"large": ImageScale{ScalePercent: 0.5},
	},
	ImageSrcDirectory:   "src",
	WebSourceLocation:   path.Join(FotoDenConfigDir, "web"), // remember when $HOME webpage folders were a thing?
	WebBaseURL:          "",                                 // this should be set during configuration generation
}
var CurrentConfig GeneratorConfig

var WorkingDirectory, _ = os.Getwd()
var Verbose bool // if this is set, everything important is printed

func verbose(print string) {
	if Verbose {
		fmt.Println(print)
	}
}

func WriteJSON(filePath string, iface interface{}) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	toWrite, err := json.Marshal(iface)
	if err != nil {
		return err
	}

	_, err = file.Write(toWrite)
	if err != nil {
		return err
	}

	return nil
}

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

// sets the current fotoDen configuration to this

func OpenfotoDenConfig(configLocation string) error {
	err := ReadJSON(configLocation, &CurrentConfig)
	if err != nil {
		return err
	}

	return nil
}

// WritefotoDenConfig
//
// Attempts to write CurrentConfig to a new file at configLocation.

func WritefotoDenConfig(config GeneratorConfig, configLocation string) error {
	err := WriteJSON(configLocation, config)
	if err != nil {
		return err
	}

	return nil
}

// removeItemFromStringArray
//
// Now externally accesible, pending refactor of how images
// are stored and accessed.

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
