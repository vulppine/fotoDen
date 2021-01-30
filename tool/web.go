package tool

import (
	"fmt"
	"github.com/vulppine/fotoDen/generator"
	"path"
)

type Theme struct {
	ThemeName   string
	Stylesheets []string
	Scripts     []string
}

func ReadThemeConfig(fpath string) (*Theme, error) {
	t := new(Theme)

	err := generator.ReadJSON(fpath, t)
	if checkError(err) {
		return nil, fmt.Errorf("valid theme could not be read")
	}

	return t, nil
}

func GenerateWeb(m string, dest string, f *generator.Folder, opt GeneratorOptions) error {
	var err error
	var pageOptions *generator.StaticWebVars

	if opt.static || f.IsStatic {
		pageOptions, err = generator.NewStaticWebVars(dest)
		if checkError(err) {
			return err
		}
	} else {
		pageOptions = new(generator.StaticWebVars)
		pageOptions.IsStatic = false
	}

	switch m {
	case "album":
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "album-template.html"), path.Join(dest, "index.html"), pageOptions)
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "photo-template.html"), path.Join(dest, "photo.html"), pageOptions)
	case "folder":
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "folder-template.html"), path.Join(dest, "index.html"), pageOptions)
	default:
		return fmt.Errorf("mode was not passed to GenerateWeb")
	}

	if checkError(err) {
		return err
	}

	return nil
}

func UpdateWeb(folder string) error {
	f := new(generator.Folder)

	err := f.ReadFolderInfo(path.Join(folder, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	err = GenerateWeb(f.FolderType, folder, f, genoptions)
	if checkError(err) {
		return err
	}

	return nil
}
