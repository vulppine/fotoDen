package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// GetArrayOfFilesAndFolders takes an array of os.FileInfo (usually from from os.Readdir()), and returns a string array of all non-directories.
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

// GetArrayOfFiles takes an array of os.FileInfo, and runs it through GetArrayOfFilesAndFolders, returning only the array of files.
func GetArrayOfFiles(directory []os.FileInfo) []string {
	fileArray, _ := GetArrayOfFilesAndFolders(directory)
	return fileArray
}

// GetArrayOfFolders takes an array of os.FileInfo, and runs it through GetArrayOfFilesAndFolders, returning only the array of folders.
func GetArrayOfFolders(directory []os.FileInfo) []string {
	_, folderArray := GetArrayOfFilesAndFolders(directory)
	return folderArray
}

// CopyFile takes three arguments - the name of the file, the name of the new file, and the destination.
// This assumes the file we're renaming is in the current working directory, or it is reachable
// via the current working directory.
//
// Returns an error if one occurs - otherwise returns nil.
func CopyFile(file string, dest string) error {
	verbose("Copying " + file + " to " + filepath.Join(dest))
	fileReader, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	toWrite, err := os.Create(filepath.Join(dest))
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
