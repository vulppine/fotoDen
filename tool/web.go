package tool

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/vulppine/fotoDen/generator"
)

/*
// Theme is a representation of the JSON file included with every theme,
// aka 'theme.json'.
type Theme struct {
	ThemeName   string
	Stylesheets []string
	Scripts     []string
	Other       []string
}

// ReadThemeConfig reads a theme.json file, and returns a Theme struct.
func ReadThemeConfig(fpath string) (*Theme, error) {
	t := new(Theme)

	err := generator.ReadJSON(fpath, t)
	if checkError(err) {
		return nil, fmt.Errorf("valid theme could not be read")
	}

	return t, nil
}
*/

// GenerateWeb hooks into the fotoDen/generator package, and generates web pages for fotoDen folders/albums.
// If a folder was marked as static, it will create a StaticWebVars object to put into the page, and generate
// a more static page - otherwise, it will leave those fields blank, and leave it to the fotoDen front end
// to generate the rest of the page.
func GenerateWeb(m string, dest string, f *generator.Folder, opt GeneratorOptions) error {
	verbose("Generating web pages...")
	var err error
	var pageOptions *generator.StaticWebVars

	if opt.Static || f.Static {
		verbose("Folder/album is static, generating static web vars...")
		pageOptions, err = generator.NewStaticWebVars(dest)
		if checkError(err) {
			return err
		}
	} else {
		verbose("Folder/album is dynamic.")
		pageOptions = new(generator.StaticWebVars)
		pageOptions.IsStatic = false
	}

	switch m {
	case "album":
		verbose("Album mode selected, generating album.")
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "album-template.html"), path.Join(dest, "index.html"), pageOptions)
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "photo-template.html"), path.Join(dest, "photo.html"), pageOptions)
	case "folder":
		verbose("Folder mode selected, generating folder.")
		err = generator.ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "folder-template.html"), path.Join(dest, "index.html"), pageOptions)
	default:
		return fmt.Errorf("mode was not passed to GenerateWeb")
	}

	if checkError(err) {
		return err
	}

	return nil
}

// UpdateWeb takes a folder, and updates the webpages inside of that folder.
func UpdateWeb(folder string) error {
	verbose("Updating web pages...")
	f := new(generator.Folder)

	err := f.ReadFolderInfo(path.Join(folder, "folderInfo.json"))
	if checkError(err) {
		return err
	}

	err = GenerateWeb(f.Type, folder, f, Genoptions)
	if checkError(err) {
		return err
	}

	return nil
}

// MARKDOWN SUPPORT //

// GeneratePage generates a page using a markdown document as a source.
// It will use the 'page' HTML template in the theme in order to generate
// a web page. Takes a file destination and a file location.
//
// The page template must have {{.PageContent}} in the location of where
// you want the parsed document to go.
func GeneratePage(dest string, src string, title string) error {
	f, err := os.Open(src)
	if checkError(err) {
		return err
	}

	r, err := io.ReadAll(f)
	if checkError(err) {
		return err
	}

	f.Close()

	mdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	var b bytes.Buffer
	err = mdown.Convert(r, &b)
	if checkError(err) {
		return err
	}

	r, err = io.ReadAll(&b)
	if checkError(err) {
		return err
	}

	// maybe it isn't a good idea to nest this in an anonymous struct???
	// then again, the generator package is supposed to be strictly
	// for fotoDen use, and having static web pages is completely
	// and utterly optional
	err = generator.ConfigureWebFile(
		path.Join(generator.CurrentConfig.WebSourceLocation, "html", "page-template.html"),
		dest,
		struct { PageContent string; Title string }{ string(r), title },
	)
	if checkError(err) {
		return err
	}

	return nil
}
