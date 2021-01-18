package generator

import (
	"testing"
	"io/ioutil"
	"path"
	"fmt"
)

func TestJSONRW(t *testing.T) {
	dir := t.TempDir()

	type JSON struct {
		Test string
	}

	err := WriteJSON(path.Join(dir, "tmp_json.json"), JSON{"test"})
	if err != nil {
		t.Errorf("Error - WriteJSON: " + fmt.Sprint(err))
	}

	json := new(JSON)

	err = ReadJSON(path.Join(dir, "tmp_json.json"), json)
	if err != nil {
		t.Errorf("Error - ReadJSON: " + fmt.Sprint(err))
	}

	if json.Test != "test" {
		t.Errorf("Error - ReadJSON: test string does not match: " + json.Test)
	}

	t.Log(fmt.Sprint(json))
}

func TestFotoDenConfigRW(t *testing.T) {
	dir := t.TempDir()

	err := WritefotoDenConfig(DefaultConfig, path.Join(dir, "tmp_config.json"))
	if err != nil {
		t.Errorf("Error - WritefotoDenConfig: " + fmt.Sprint(err))
	}

	f, _ := ioutil.ReadFile(path.Join(dir, "tmp_config.json"))
	t.Log(fmt.Printf("%s", f))

	err = OpenfotoDenConfig(path.Join(dir, "tmp_config.json"))
	if err != nil {
		t.Errorf("Error - OpenfotoDenConfig: "+ fmt.Sprint(err))
	}
	t.Log(fmt.Sprint(CurrentConfig))
}

func TestFolderInfoCRW(t *testing.T) {
	dir := t.TempDir()

	folder, err := GenerateFolderInfo(dir, dir)
	if err != nil {
		t.Errorf("Error - GenerateFolderInfo: " + fmt.Sprint(err))
	}

	err = folder.WriteFolderInfo((path.Join(dir, "folderInfo.json")))
	if err != nil {
		t.Errorf("Error - WriteFolderInfo: " + fmt.Sprint(err))
	}

	err = folder.ReadFolderInfo((path.Join(dir, "folderInfo.json")))
	if err != nil {
		t.Errorf("Error - ReadFolderInfo: " + fmt.Sprint(err))
	}
	t.Log(fmt.Sprint(folder))
}
