package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestJSONRW(t *testing.T) {
	dir := t.TempDir()

	type JSON struct {
		Test string
	}

	err := WriteJSON(path.Join(dir, "tmp_json.json"), "single", JSON{"test"})
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
	t.Log(string(f))

	err = OpenfotoDenConfig(path.Join(dir, "tmp_config.json"))
	if err != nil {
		t.Errorf("Error - OpenfotoDenConfig: " + fmt.Sprint(err))
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

func TestItemsInfoCRW(t *testing.T) {
	dir := t.TempDir()

	items, err := GenerateItemInfo("../test_images")
	if err != nil {
		t.Errorf("Error - GenerateItemInfo: " + fmt.Sprint(err))
	}

	err = items.WriteItemsInfo(path.Join(dir, "itemsInfo.json"))
	if err != nil {
		t.Errorf("Error - WriteItemsInfo: " + fmt.Sprint(err))
	}

	err = items.ReadItemsInfo(path.Join(dir, "itemsInfo.json"))
	if err != nil {
		t.Errorf("Error: ReadItemsInfo: " + fmt.Sprint(err))
	}
	t.Log(fmt.Sprint(items))
}

func TestBatchCopyConvert(t *testing.T) {
	dir := t.TempDir()
	dir2 := path.Join(dir, t.TempDir())

	src, err := ioutil.ReadDir("../test_images")
	if err != nil {
		t.Errorf("Error: Opening test images folder: " + fmt.Sprint(err))
	}
	srcfiles := GetArrayOfFiles(src)

	defer os.Chdir(WorkingDirectory)
	os.Chdir("../test_images")

	err = BatchCopyFile(srcfiles, dir)
	if err != nil {
		t.Errorf("Error: BatchCopyFile: " + fmt.Sprint(err))
	}

	os.Chdir(dir)

	err = BatchImageConversion(srcfiles, "test", dir2, ImageScale{ScalePercent: 0.99})
	if err != nil {
		t.Errorf("Error: BatchImageConversion: " + fmt.Sprint(err))
	}
}

func TestWebConfigCRW(t *testing.T) {
	dir := t.TempDir()

	webconfig := GenerateWebConfig("https://localhost/")

	err := webconfig.WriteWebConfig(path.Join(dir, "config.json"))
	if err != nil {
		t.Errorf("Error - WriteWebConfig: " + fmt.Sprint(err))
	}

	err = webconfig.ReadWebConfig(path.Join(dir, "config.json"))
	if err != nil {
		t.Errorf("Error - ReadWebConfig: " + fmt.Sprint(err))
	}

	t.Log(fmt.Sprint(webconfig))
}
