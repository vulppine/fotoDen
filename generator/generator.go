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
	SupFolderName       string   // The name of the folder above it (for presentation purposes).
	FolderType          string   // The type of folder (currently supports only album or folder)
	FolderThumbnail     bool     // If the folder has a thumbnail or not
	ItemsInFolder       []string // All the items in the folder, by name.
	SubfolderShortNames []string // Any folders that are within the folder (updated whenever the generator is called in the folder)
}

type GeneratorConfig struct {
	ThumbMaxHeight      int
	ThumbMaxWidth       int
	ThumbScalePercent   float32
	LargeScalePercent   float32
	ImageRootDirectory  string // where all images are stored (default: img)
	ImageThumbDirectory string // where all thumbnails are stored (default: ImageRootDirectory/thumb)
	ImageLargeDirectory string // where all 'large' display images are stored (default: ImageRootDirectory/large)
	ImageSrcDirectory   string // where all source images are stored (default: ImageRootDirectory/src)
	ImageJSONDirectory  string // where all JSON files per image are stored (default: ImageRootDirectory/json)
	WebSourceLocation   string // where all html/css/js files are stored for fotoDen's functionality
	WebBaseURL          string // what the base URL is (aka, fotoDen's location)
}

// some defaults in case we never have a fotoDen config file opened

var userConfigDir, _ = os.UserConfigDir()
var FotoDenConfigDir = path.Join(userConfigDir, "fotoDen")
var DefaultConfig GeneratorConfig = GeneratorConfig{
	ThumbMaxHeight:      800,
	LargeScalePercent:   0.5,
	ImageRootDirectory:  "img",
	ImageThumbDirectory: "img/thumb",
	ImageLargeDirectory: "img/large",
	ImageSrcDirectory:   "img/src",
	ImageJSONDirectory:  "img/json",
	WebSourceLocation:   path.Join(FotoDenConfigDir, "web"), // remember when $HOME webpage folders were a thing?
	WebBaseURL:          "",                                 // this should be set during configuration generation
}
var CurrentConfig GeneratorConfig
var ThumbScalingOptions ImageScale
var LargeScalingOptions ImageScale

var WorkingDirectory, _ = os.Getwd()
var Verbose bool // if this is set, everything important is printed

func verbose(print string) {
	if Verbose {
		fmt.Println(print)
	}
}

// sets the current fotoDen configuration to this

func OpenfotoDenConfig(configLocation string) error {
	configFile, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configFile, &CurrentConfig)
	if err != nil {
		return err
	}

	ThumbScalingOptions = ImageScale{
		maxheight:    CurrentConfig.ThumbMaxHeight,
		maxwidth:     CurrentConfig.ThumbMaxWidth,
		scalepercent: CurrentConfig.ThumbScalePercent,
	}

	LargeScalingOptions = ImageScale{
		scalepercent: CurrentConfig.LargeScalePercent,
	}

	return nil
}

// WritefotoDenConfig
//
// Attempts to write CurrentConfig to a new file at configLocation.

func WritefotoDenConfig(config GeneratorConfig, configLocation string) error {
	configFile, err := os.OpenFile(configLocation, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer configFile.Close()

	toWrite, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = configFile.Write(toWrite)
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
