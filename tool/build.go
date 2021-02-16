package tool

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// BuildFile represents a fotoDen build file.
type BuildFile struct {
	Name     string   `yaml:"name"`
	Dir      string   `yaml:"dir"`
	Desc     string   `yaml:"desc"`
	Type     string   `yaml:"type"`
	Thumb    string   `yaml:"thumb"`
	Static   bool     `yaml:"static"`
	ImageDir string   `yaml:"imageDir"`
	Images   []string `yaml:"images,flow"`
	Options  struct {
		Copy     bool `yaml:"copy"`
		Sort     bool `yaml:"sort"`
		Meta     bool `yaml:"metadata"`
		Gensizes bool `yaml:"generateSizes"`
	} `yaml:"imageOptions,flow"`
	Subfolders []*BuildFile `yaml:"subfolders,flow"`
}

// OpenYAML opens a build file into a new BuildFile object.
func (b *BuildFile) OpenBuildYAML(file string) error {
	f, err := ioutil.ReadFile(file)
	if checkError(err) {
		return err
	}

	err = yaml.Unmarshal(f, &b)
	if checkError(err) {
		return err
	}

	return nil
}

// BuildFromYAML is the entry point to the fotoDen
// website build system. It takes a YAML file,
// with the correct structure, and creates a website
// in the given folder.
func (b *BuildFile) Build(folder string) error {
	verbose(fmt.Sprint(b))
	switch b.Type {
	case "folder":
		genopts := GeneratorOptions{
			Static: b.Static,
		}
		if b.Dir == "" {
			b.Dir = b.Name
		}
		err := GenerateFolder(
			FolderMeta{
				Name: b.Name,
				Desc: b.Desc,
			},
			filepath.Join(folder, b.Dir),
			genopts,
		)
		if checkError(err) {
			return err
		}

		for _, f := range b.Subfolders {
			if f.Dir == "" {
				f.Dir = f.Name
			}
			b.Dir, _ = filepath.Abs(filepath.Join(folder, b.Dir))
			err = f.Build(b.Dir)
			if checkError(err) {
				return err
			}
		}
	case "album":
		genopts := GeneratorOptions{
			ImageGen: true,
			Copy:     b.Options.Copy,
			Sort:     b.Options.Sort,
			Meta:     b.Options.Meta,
			Static:   b.Static,
			Gensizes: b.Options.Gensizes,
		}
		if b.Dir == "" {
			b.Dir = b.Name
		}
		if b.ImageDir != "" {
			genopts.Source = b.ImageDir
			err := GenerateFolder(
				FolderMeta{
					Name: b.Name,
					Desc: b.Desc,
				},
				filepath.Join(folder, b.Dir),
				genopts,
			)
			if checkError(err) {
				return err
			}

			InsertImage(folder, "append", Genoptions, b.Images...)
		} else {
			err := GenerateFolder(
				FolderMeta{
					Name: b.Name,
					Desc: b.Desc,
				},
				filepath.Join(folder, b.Dir),
				genopts,
			)
			if checkError(err) {
				return err
			}

			InsertImage(folder, "append", Genoptions, b.Images...)
		}

		for _, f := range b.Subfolders {
			if f.Dir == "" {
				f.Dir = f.Name
			}
			b.Dir, _ = filepath.Abs(filepath.Join(folder, b.Dir))
			err := f.Build(b.Dir)
			if checkError(err) {
				return err
			}
		}
	}

	return nil
}
