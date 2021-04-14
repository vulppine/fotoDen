package tool

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/vulppine/fotoDen/generator"
)

// Note:
// This new method of zipped themes should immediately
// replace the old method of unzipping an uninitialized
// theme into a folder. It should make things more
// easier, as rather than having to initialize a theme
// into a directory, all a user has to do in order
// to create a site is just throw a --theme flag into
// the command. Equally, a default theme should be
// embed during build in case the user doesn't have
// a custom theme.zip to use.
//
// The old method should be deprecated immediately.

// theme represents a struct containing an open
// zip archive for reading, plus a struct containing
// an internal theme.json file.
type theme struct {
	a *zip.Reader
	s struct {
		ThemeName   string
		Stylesheets []string
		Scripts     []string
		Other       []string
	}
}

func zipFileReader(z string) (io.ReaderAt, int64, error) {
	f, err := os.Open(z)
	if checkError(err) {
		return nil, 0, err
	}

	fi, err := f.Stat()
	if checkError(err) {
		return nil, 0, err
	}

	return f, fi.Size(), nil
}

func openTheme(z io.ReaderAt, zs int64) (*theme, error) {
	var err error

	t := new(theme)

	t.a, err = zip.NewReader(z, zs)
	if checkError(err) {
		return nil, err
	}

	c, err := t.a.Open("theme.json")
	if checkError(err) {
		return nil, err
	}

	cb, err := io.ReadAll(c)
	if checkError(err) {
		return nil, err
	}

	json.Unmarshal(cb, &t.s)

	verbose("theme successfully opened")
	return t, nil
}

// writeFile directly writes a file from a theme's zip.Reader
// to the given destination.
func (t *theme) writeFile(n string, d string) error {
	verbose("attempting to write a file from zip: " + n)
	f, err := t.a.Open(n)
	if checkError(err) {
		return err
	}
	fb, err := io.ReadAll(f)
	if checkError(err) {
		return err
	}

	g, err := os.Create(d)
	if checkError(err) {
		return err
	}

	_, err = g.Write(fb)
	if checkError(err) {
		return err
	}

	return nil
}

// writeDir writes the contents of a named directory in a
// theme's zip.Reader into the named directory d, creating
// a directory in the process. A list of files can be given
// to writeDir, allowing it to skip checking the entire
// directory.
//
// NOTE: I haven't implemented the directory read yet.
// This is expected to be used in conjunction with
// the theme.json setup!
func (t *theme) writeDir(n string, d string, files ...string) error {
	verbose("attempting to write a directory from zip: " + n)
	i, err := t.a.Open(n)
	if checkError(err) {
		return err
	}
	in, err := i.Stat()
	if checkError(err) {
		return err
	}

	if !in.IsDir() {
		return errors.New("not a valid directory")
	}

	if len(files) == 0 {
		return errors.New("zip file directory checking not implemented yet")
	}

	if !fileCheck(filepath.Join(d, n)) {
		err = os.Mkdir(filepath.Join(d, n), 0755)
		if checkError(err) {
			return err
		}
	}

	for _, f := range files {
		err := t.writeFile(filepath.Join(n, f), filepath.Join(d, n, f))
		if checkError(err) {
			return err
		}
	}

	return nil
}

// copyTheme copies a theme to a destination string,
// making a directory within d containing the entire
// tree of the theme as described in theme.json,
// as well as copying over theme.json for later use.
func (t *theme) copyTheme(d string) error {

	return nil
}

// initTheme initializes a theme to a site, replacing
// all relevant variables according to what fotoDen
// needs in order to generate a gallery.
// u is the URL of the site
// e is the template directory of the site
// r is the root directory of the site
func (t *theme) initTheme(u string, e string, r string) error {
	verbose("attempting to initialize a theme")
	m, err := os.MkdirTemp("", "")
	if checkError(err) {
		return err
	}
	verbose("temp dir created in: " + m)

	wvars, err := generator.NewWebVars(u)
	wvars.StaticWebVars["pageContent"] = "{{.PageContent}}" // hacky
	checkError(err)

	hf := []string{"photo-template.html", "album-template.html", "folder-template.html", "page-template.html"}

	verbose("generating temporary template setup files")
	for _, f := range hf {
		err = t.writeFile(filepath.Join("html", f), filepath.Join(m, f))
		if checkError(err) {
			return err
		}
	}

	verbose("writing HTML templates now to: " + e)
	err = os.MkdirAll(filepath.Join(e, "html"), 0755)
	if checkError(err) {
		return err
	}
	for _, f := range hf {
		err = generator.ConfigureWebFile(
			filepath.Join(m, f),
			filepath.Join(e, "html", f),
			wvars,
		)
		if checkError(err) {
			return err
		}
	}

	copyArray := func(f []string, n string) error {
		if len(f) != 0 {
			err = t.writeDir(n, r, f...)
			if checkError(err) {
				return err
			}
		}

		return nil
	}

	verbose("copying over theme folders into: " + r)
	err = copyArray(t.s.Stylesheets, "css")
	err = copyArray(t.s.Scripts, "js")
	err = copyArray(t.s.Other, "etc")

	return err
}
