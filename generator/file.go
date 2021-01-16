package generator

import (
	"path"
	"os"
	"io/ioutil"
)

// GetArrayOfFilesAndFolders
//
// Takes an array of os.FileInfo (from os.Readdir()), and returns a string array of all non-directories.
// Also returns a string array of directories, so we don't have to copy and paste this function.

func GetArrayOfFilesAndFolders(directory []os.FileInfo) ([]string, []string) {
	fileArray := make([]string, 0)
	folderArray := make([]string, 0)

	for i := 0; i < len(directory); i++ {
		if directory[i].Mode().IsDir() == true {
			folderArray = append(folderArray, directory[i].Name())
		} else {
			fileArray = append(fileArray, directory[i].Name())
		}
	}

	return fileArray, folderArray
}

// Wrappers for both, just in case only one of the two is needed

func GetArrayOfFiles(directory []os.FileInfo) []string {
	fileArray, _ := GetArrayOfFilesAndFolders(directory)
	return fileArray
}

func GetArrayOfFolders(directory []os.FileInfo) []string {
	_, folderArray := GetArrayOfFilesAndFolders(directory)
	return folderArray
}


// CopyFile
//
// Takes three arguments - the name of the file, the name of the new file, and the destination.
// This assumes the file we're renaming is in the current working directory, or it is reachable
// via the current working directory.
//
// Returns an error if one occurs - otherwise returns nil.

func CopyFile(file string, name string, dest string) error {
	verbose("Copying " + file + " to " + path.Join(dest, name))
	fileReader, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	toWrite, err := os.Create(path.Join(dest, name))
	defer toWrite.Close()
	if err != nil {
		return err
	}

	_, err = toWrite.Write(fileReader)
	if err != nil {
		return err
	}

	return nil
}

// MakeAlbumDirectoryStructure
//
// Makes a fotoDen-suitable album structure in the given rootDirectory (string).
// The directory must exist beforehand.

func MakeAlbumDirectoryStructure(rootDirectory string) error {

	currentDirectory, _ := os.Getwd()

	defer func() {
		verbose("Changing back to " + currentDirectory)
		os.Chdir(currentDirectory)
	}()

	verbose("Attempting to change to " + rootDirectory)
	err := os.Chdir(rootDirectory)
	if err != nil {
		return err
	}

	verbose("Creating directories in " + rootDirectory)
	os.Mkdir(CurrentConfig.ImageRootDirectory, 0777)
	os.Mkdir(CurrentConfig.ImageThumbDirectory, 0777)
	os.Mkdir(CurrentConfig.ImageSrcDirectory, 0777)
	os.Mkdir(CurrentConfig.ImageLargeDirectory, 0777)
	os.Mkdir(CurrentConfig.ImageJSONDirectory, 0777)

	return nil
}
