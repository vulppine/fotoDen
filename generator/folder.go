package generator

import (
	"io/ioutil"
	"os"
	"path"
)

// Folder represents a folderInfo.json file used by fotoDen.
// It has all the needed values for fotoDen.js to operate correctly.
// fotoDen/generator does not provide functions to manage this - only to read and create these.
type Folder struct {
	Name       string   `json:"name"`      // The name of the folder.
	Desc       string   `json:"desc"`      // The description of a folder.
	ShortName  string   `json:"shortName"` // The shortname of the folder (can be taken from the filesystem folder name)
	Type       string   `json:"type"`      // The type of folder (currently supports only album or folder)
	Thumbnail  bool     `json:"thumbnail"` // If a thumbnail exists or not. This is dictated by the generation of thumb.jpg.
	ItemAmount int      `json:"itemAmount"`
	Subfolders []string `json:"subfolders"` // Any folders that are within the folder (updated whenever the generator is called in the folder)
	Static     bool     `json:"static"`     // If the folder was generated statically, or has information inserted dynamically.
}

// GenerateFolderInfo generates a Folder object that can be used for folder configuration.
// If directory is an empty string, it does it in the current directory.
// Otherwise, it attempts to reach the directory from the current working directory.
//
// If name is an empty string, it uses the target directory's name.
//
// Does not check if folderType is valid.
//
// Returns a Folder object if successful, and a nil error,
// otherwise returns a potentially incomplete object, and an error.
func GenerateFolderInfo(directory string, name string) (*Folder, error) {
	folder := new(Folder)

	if directory != "" {
		wd, _ := os.Getwd()
		defer os.Chdir(wd)
		err := os.Chdir(directory)
		if err != nil {
			return folder, err
		}
	}

	currentDirectory, _ := os.Getwd()
	verbose("Generating folder info from: " + currentDirectory)

	if name == "" {
		folder.Name = path.Base(currentDirectory)
	} else {
		folder.Name = name
	}

	folder.ShortName = path.Base(currentDirectory)
	folder.Subfolders = []string{}

	return folder, nil
}

// ReadFolderInfo is a method for reading folder info from a file.
// Returns an error if any occur.
func (folder *Folder) ReadFolderInfo(filePath string) error {
	verbose("Reading folder infomation from " + filePath)
	err := ReadJSON(filePath, folder)
	if err != nil {
		return err
	}

	return nil
}

// WriteFolderInfo is a method for writing fotoDen folder info to a file.
// Returns an error if any occur.
func (folder *Folder) WriteFolderInfo(filePath string) error {
	verbose("Writing folder (" + folder.ShortName + ") to " + filePath)
	err := WriteJSON(filePath, "multi", folder)
	if err != nil {
		return err
	}

	return nil
}

// Items represents an itemsInfo.json file used by fotoDen.
// It is used mainly in album-type folders, and contains a bool indicating whether
// metadata is being used, and a string array (potentially large) of file names.
type Items struct {
	Metadata bool `json:"metadata"` // Dictates whether or not each image has its own ImageMeta object.
	// If this is false, then no metadata will be read.
	ItemsInFolder []string `json:"items"` // All the items in a folder, by name, in an array.
}

// GenerateItemInfo generates an Items object based on the contents of the directory.
// This automatically strips non-images.
func GenerateItemInfo(directory string) (*Items, error) {
	items := new(Items)

	verbose("Reading items in folder: " + directory)

	dir, err := os.Open(directory)
	defer dir.Close()
	if err != nil {
		return items, err
	}

	dirContents, err := dir.Readdir(0)
	if err != nil {
		return items, err
	}
	defer os.Chdir(WorkingDirectory)
	os.Chdir(directory)
	items.ItemsInFolder = IsolateImages(GetArrayOfFiles(dirContents))

	return items, nil
}

// ReadItemsInfo is a method for reading items info from a file.
// Returns an error if any occur.
func (items *Items) ReadItemsInfo(filePath string) error {
	verbose("Reading items infomation from " + filePath)
	err := ReadJSON(filePath, items)
	if err != nil {
		return err
	}

	return nil
}

// WriteItemsInfo is a method for writing items info to a file.
// Returns an error if any occur.
func (items *Items) WriteItemsInfo(filePath string) error {
	verbose("Writing items to " + filePath)
	err := WriteJSON(filePath, "single", items)
	if err != nil {
		return err
	}

	return nil
}

// perhaps this should be moved to tool?
// no - it's useful here

// UpdateSubdirectories updates a Folder object's subdirectories according to the given directory.
// If directory is an empty string, will attempt to update the Folder from the current working directory.
// Returns an error, if any occurs, otherwise the number of directories and a nil error.
func (folder *Folder) UpdateSubdirectories(directory string) (int, error) {
	var currentDirectory string
	switch directory {
	case "":
		currentDirectory, _ = os.Getwd()
	default:
		currentDirectory = directory
	}

	verbose("Updating subdirectories in " + directory)

	fileArray, err := ioutil.ReadDir(currentDirectory)
	if err != nil {
		return 0, err
	}

	// special cases for the root css/js directories, and the album image location
	folder.Subfolders = RemoveItemFromStringArray(GetArrayOfFolders(fileArray), CurrentConfig.ImageRootDirectory)

	// really lazy, find a better way to do this
	folder.Subfolders = RemoveItemFromStringArray(folder.Subfolders, "theme")
	folder.Subfolders = RemoveItemFromStringArray(folder.Subfolders, "js")
	return len(folder.Subfolders), nil
}
