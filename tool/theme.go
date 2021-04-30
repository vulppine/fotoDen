package tool

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"text/template"

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
// an internal theme.json file. It also caches some
// files for later use via the text/template system.
type theme struct {
	a *zip.Reader
	s struct {
		ThemeName   string
		Stylesheets []string
		Scripts     []string
		Other       []string
	}
	f map[string]string // file cache
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

	verbose("theme config successfully opened")

	t.f = make(map[string]string)
	f := []string{
		"album-template.html",
		"folder-template.html",
		"photo-template.html",
		"page-template.html",
	}

	for _, i := range f {
		f, err := t.a.Open(filepath.Join("html", i))
		if err != nil {
			verbose("could not open a required page template: " + i)
			return nil, err
		}

		b, err := io.ReadAll(f)
		if err != nil {
			verbose("could not open a required page template: " + i)
			return nil, err
		}

		t.f[i] = string(b)
	}

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

// webVars dictate where fotoDen gets its JavaScript and CSS files per page.
// Four variables are definite - and BaseURL is the most important one.
// PageVars indicate optional variables to be passed to the Go template engine.
type webVars struct {
	BaseURL     string
	JSLocation  string
	CSSLocation string
	IsStatic    bool
	PageVars    map[string]string
}

// newWebVars creates a WebVars object. Takes a single URL string, and outputs
// a set of fotoDen compatible URLs.
func newWebVars(u, folder string) (*webVars, error) {

	webvars := new(webVars)
	url, err := url.Parse(u)
	jsurl, err := url.Parse(u)
	cssurl, err := url.Parse(u)
	if err != nil {
		return webvars, err
	}

	webvars.BaseURL = url.String()
	if len(webvars.BaseURL) > 0 && webvars.BaseURL[len(webvars.BaseURL)-1] == '/' {
		webvars.BaseURL = webvars.BaseURL[0 : len(webvars.BaseURL)-1]
	}

	f := new(generator.Folder)
	fpath, _ := filepath.Abs(folder)

	err = f.ReadFolderInfo(filepath.Join(fpath, "folderInfo.json"))
	if err != nil {
		return nil, err
	}

	superFolder, err := func() (string, error) {
		f := new(generator.Folder)

		_, err := os.Stat(filepath.Join(filepath.Dir(fpath), "folderInfo.json"))
		if os.IsNotExist(err) {
			return "", nil
		} else if checkError(err) {
			return "", err
		}

		verbose("Folder above is a fotoDen folder, using that...")
		err = f.ReadFolderInfo(filepath.Join(filepath.Dir(fpath), "folderInfo.json"))
		if checkError(err) {
			return "", err
		}

		return f.Name, nil
	}()

	if checkError(err) {
		verbose("Could not read folder above the current one")
	}

	jsurl.Path = path.Join(jsurl.Path, "js", "fotoDen.js")
	webvars.JSLocation = jsurl.String()

	cssurl.Path = path.Join(cssurl.Path, "css", "theme.css")
	webvars.CSSLocation = cssurl.String()

	if f.Static {
		webvars.PageVars = map[string]string{
			"name": f.Name,
			"desc": f.Desc,
			"sfol": superFolder,
		}
	}

	return webvars, nil
}

type page int

const (
	photo page = iota
	album
	folder
	info
)

// configurePage configures the various Go template
// variables within a page according to a specific type.
// u is the URL of a website,
// d is the destination that the result goes into,
// t is the type of page,
// i is the set of variables to use
func (t *theme) configurePage(u, d string, y page, i *webVars) error {
	var f string
	p := template.New("result")

	r, err := os.Create(d)
	if err != nil {
		return err
	}

	switch y {
	case photo:
		f = t.f["photo-template.html"]
	case album:
		f = t.f["album-template.html"]
	case folder:
		f = t.f["folder-template.html"]
	case info:
		f = t.f["page-template.html"]
	}

	m, err := p.Parse(f)
	if err != nil {
		return err
	}

	return m.Execute(r, i)
}

// generateWeb takes a mode, a destinatination, and an optional map[string]string.
// If i is not nil, that map will be merged into the WebVars PageVars field.
func (t *theme) generateWeb(m, dest string, i map[string]string) error {
	var err error
	var v *webVars

	if m == "folder" || m == "album" {
		v, err = newWebVars(generator.CurrentConfig.WebBaseURL, dest)
		if err != nil {
			return err
		}
	} else {
		v = new(webVars)
		v.BaseURL = generator.CurrentConfig.WebBaseURL
		v.PageVars = make(map[string]string)
	}

	if i != nil {
		for k, a := range i {
			v.PageVars[k] = a
		}
	}

	switch m {
	case "folder":
		err = t.configurePage(generator.CurrentConfig.WebBaseURL, path.Join(dest, "index.html"), folder, v)
		if checkError(err) {
			return err
		}
	case "album":
		err = t.configurePage(generator.CurrentConfig.WebBaseURL, path.Join(dest, "index.html"), album, v)
		if checkError(err) {
			return err
		}

		err = t.configurePage(generator.CurrentConfig.WebBaseURL, path.Join(dest, "photo.html"), photo, v)
		if checkError(err) {
			return err
		}
	case "page":
		err = t.configurePage(generator.CurrentConfig.WebBaseURL, dest, info, v)

		if checkError(err) {
			return err
		}
	}

	return nil
}

/// GLOBAL THEME VAR ///
// this probably *could* be mitigated by
// including the current theme's name in the
// config, allowing for functions to refer to
// the current *website*'s theme

var currentTheme *theme

func setCurrentTheme(t string) error {
	c, err := os.UserConfigDir()
	if checkError(err) {
		return err
	}

	if t != "Default" {
		p := filepath.Join(c, "fotoDen", "themes", t+".zip")
		z, s, err := zipFileReader(p)
		if checkError(err) {
			return err
		}

		currentTheme, err = openTheme(z, s)
		if checkError(err) {
			return err
		}
	} else {
		if isEmbed {
			currentTheme, _ = openTheme(defaultThemeZipReader(), defaultThemeZipLen)
		} else {
			return fmt.Errorf("could not find a fotoDen theme to use")
		}
	}

	return nil
}

func openDefaultTheme() error {
	if fileCheck(path.Join(generator.RootConfigDir, "defaulttheme")) {
		f, err := os.Open(path.Join(generator.RootConfigDir, "defaulttheme"))
		d, err := ioutil.ReadAll(f)
		if checkError(err) {
			return err
		}

		err = setCurrentTheme(string(d))
		if checkError(err) {
			return err
		}
	} else {
		return fmt.Errorf("warning: could not find a default fotoDen theme to use")
	}

	return nil
}

// initTheme initializes a theme to a site, replacing
// all relevant variables according to what fotoDen
// needs in order to generate a gallery.
// u is the URL of the site
// e is the template directory of the site
// r is the root directory of the site
func (t *theme) initTheme(u string, e string, r string) error {
	var err error
	verbose("attempting to initialize a theme")
	/*m, err := os.MkdirTemp("", "")
	if checkError(err) {
		return err
	}
	verbose("temp dir created in: " + m)

	wvars, err := NewWebVars(u)
	wvars.PageVars["pageContent"] = "{{.PageContent}}" // hacky
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
		err = ConfigureWebFile(
			filepath.Join(m, f),
			filepath.Join(e, "html", f),
			wvars,
		)
		if checkError(err) {
			return err
		}
	}*/

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
