package tool

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/vulppine/fotoDen/generator"
)

// TestGenerateFolder
//
// Incidentally, this covers the GenerateItems() function too,
// so there's no need to test that one here.

func TestGenerateFolder(t *testing.T) {
	generator.CurrentConfig = generator.DefaultConfig

	dir := t.TempDir()
	genopts := GeneratorOptions{
		ImageGen: false,
	}

	err := GenerateFolder("without_images", path.Join(dir, "without_images"), genopts)
	if err != nil {
		t.Errorf("Error - GenerateFolder (no images): " + fmt.Sprint(err))
	}

	f, _ := ioutil.ReadDir(path.Join(dir, "without_images"))
	t.Log(generator.GetArrayOfFilesAndFolders(f))
	j, _ := ioutil.ReadFile(path.Join(dir, "without_images", "folderInfo.json"))
	t.Log(string(j))

	genopts = GeneratorOptions{
		Source:   "../test_images",
		Copy:     true,
		Gensizes: true,
		ImageGen: true,
		Sort:     true,
	}

	err = GenerateFolder("with_images", path.Join(dir, "with_images"), genopts)
	if err != nil {
		t.Errorf("Error - GenerateFolder (with images): " + fmt.Sprint(err))
	}

	f, _ = ioutil.ReadDir(path.Join(dir, "with_images"))
	t.Log(generator.GetArrayOfFilesAndFolders(f))
	f, _ = ioutil.ReadDir(path.Join(dir, "with_images", "img"))
	t.Log(generator.GetArrayOfFilesAndFolders(f))
	f, _ = ioutil.ReadDir(path.Join(dir, "with_images", "img", "thumb"))
	t.Log(generator.GetArrayOfFilesAndFolders(f))
	f, _ = ioutil.ReadDir(path.Join(dir, "with_images", "img", "large"))
	t.Log(generator.GetArrayOfFilesAndFolders(f))
	f, _ = ioutil.ReadDir(path.Join(dir, "with_images", "img", "src"))
	t.Log(generator.GetArrayOfFilesAndFolders(f))
	j, _ = ioutil.ReadFile(path.Join(dir, "with_images", "folderInfo.json"))
	t.Log(string(j))
	j, _ = ioutil.ReadFile(path.Join(dir, "with_images", "itemsInfo.json"))
	t.Log(string(j))
}

// holy shit, CRUD??? who could've ever guessed

func TestImageCRUD(t *testing.T) {
	generator.CurrentConfig = generator.DefaultConfig

	dir := t.TempDir()

	items, err := generator.GenerateItemInfo("../test_images")
	if err != nil {
		t.Errorf("Error - generator.GenerateItemInfo" + fmt.Sprint(err))
	}

	err = items.WriteItemsInfo(path.Join(dir, "itemsInfo.json"))
	if err != nil {
		t.Errorf("Error - generator.WriteItemInfo" + fmt.Sprint(err))
	}

	func() {
		wd, _ := os.Getwd()
		defer os.Chdir(wd)
		os.Chdir(dir)
		generator.MakeAlbumDirectoryStructure(dir)
	}()

	genopts := GeneratorOptions{
		Source:   "../test_images",
		Copy:     true,
		Gensizes: true,
		Sort:     true,
	}

	err = UpdateImages(dir, genopts)
	if err != nil {
		t.Errorf("Error - UpdateImages" + fmt.Sprint(err))
	}
	t.Log(func() string {
		j, _ := ioutil.ReadFile(path.Join(dir, "itemsInfo.json"))
		return string(j)
	}())

	err = DeleteImage(dir, items.ItemsInFolder[0])
	if err != nil {
		t.Errorf("Error - DeleteImage" + fmt.Sprint(err))
	}
	t.Log(func() string {
		j, _ := ioutil.ReadFile(path.Join(dir, "itemsInfo.json"))
		return string(j)
	}())

	err = InsertImage(dir, "sort", genopts, items.ItemsInFolder[0])
	if err != nil {
		t.Errorf("Error - InsertImage" + fmt.Sprint(err))
	}
	t.Log(func() string {
		j, _ := ioutil.ReadFile(path.Join(dir, "itemsInfo.json"))
		return string(j)
	}())
}
