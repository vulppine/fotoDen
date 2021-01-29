package generator

import (
	"io/ioutil"
	"os"
	"path"
)

// GenerateFolderInfo
//
// Generates a Folder object that can be used for folder configuration.
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
		folder.FolderName = path.Base(currentDirectory)
	} else {
		folder.FolderName = name
	}

	folder.FolderShortName = path.Base(currentDirectory)
	folder.SubfolderShortNames = []string{}

	/* Move this to tool
	_, err := os.Stat(path.Join(path.Dir(currentDirectory), "folderInfo.json"))
	if err != nil {
		fmt.Println("Directory above either is not a fotoDen directory, or is missing folderInfo.json. Skipping. Folder: ", directory)
	} else {
		supFolder := new(Folder)
		err = supFolder.ReadFolderInfo(path.Join(path.Dir(currentDirectory), "folderInfo.json"))
		if err != nil {
			fmt.Println("An error occurred during folder generation: ", err)
		} else {
			supFolder.SubfolderShortNames = append(supFolder.SubfolderShortNames, folder.FolderShortName)
			supFolder.WriteFolderInfo(path.Join(path.Dir(currentDirectory), "folderInfo.json"))
		}
	}	*/

	/* Move this to tool
	_, err = os.Stat(path.Join(currentDirectory, "thumb"))
	if err != nil {
		fmt.Println("No thumbnail detected. (You can set this manually by placing a valid image into the folder named as 'thumb', and setting folderInfo.json's 'FolderThumbnail' to true.)")
		folder.FolderThumbnail = false
	} else {
		folder.FolderThumbnail = true
	}
	*/

	return folder, nil
}


/// ReadFolderInfo
//
// A method for reading folder info from a file.
// Returns an error if any occur.
func (folder *Folder) ReadFolderInfo(filePath string) error {
	verbose("Reading folder infomation from " + filePath)
	err := ReadJSON(filePath, folder)
	if err != nil {
		return err
	}

	return nil
}

// WriteFolderInfo
//
// A method for writing folder info to a file.
// Returns an error if any occur.
func (folder *Folder) WriteFolderInfo(filePath string) error {
	verbose("Writing folder (" + folder.FolderShortName + ") to " + filePath)
	err := WriteJSON(filePath, "multi", folder)
	if err != nil {
		return err
	}

	return nil
}

// GenerateItemInfo
//
// Generates an Items object based on the contents of the directory.
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

// ReadItemsInfo
//
// A method for reading items info from a file.
// Returns an error if any occur.
func (items *Items) ReadItemsInfo(filePath string) error {
	verbose("Reading items infomation from " + filePath)
	err := ReadJSON(filePath, items)
	if err != nil {
		return err
	}

	return nil
}

// WriteItemsInfo
//
// A method for writing items info to a file.
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

// UpdateSubdirectories
//
// Updates a Folder object's subdirectories according to the given directory.
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
	folder.SubfolderShortNames = RemoveItemFromStringArray(GetArrayOfFolders(fileArray), CurrentConfig.ImageRootDirectory)

	// really lazy, find a better way to do this
	folder.SubfolderShortNames = RemoveItemFromStringArray(folder.SubfolderShortNames, "theme")
	folder.SubfolderShortNames = RemoveItemFromStringArray(folder.SubfolderShortNames, "js")
	return len(folder.SubfolderShortNames), nil
}
