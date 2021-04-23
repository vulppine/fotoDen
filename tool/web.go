package tool

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/vulppine/fotoDen/generator"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// GenerateWebRoot generates the root of a fotoDen website in fpath.
// Creates the folders for JS and CSS placement.
//
// It is up to the fotoDen tool to copy over the relevant files,
// and folder configuration.
func GenerateWebRoot(fpath string) error {
	err := os.Mkdir(fpath, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(fpath, "js"), 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(fpath, "theme"), 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(fpath, "theme", "js"), 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(fpath, "theme", "css"), 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(filepath.Join(fpath, "theme", "etc"), 0755)
	if err != nil {
		return err
	}

	return nil
}

// MakeAlbumDirectoryStructure makes a fotoDen-suitable album structure in the given rootDirectory (string).
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
	os.Mkdir(generator.CurrentConfig.ImageRootDirectory, 0777)
	os.Mkdir(path.Join(generator.CurrentConfig.ImageRootDirectory, generator.CurrentConfig.ImageSrcDirectory), 0777)
	os.Mkdir(path.Join(generator.CurrentConfig.ImageRootDirectory, generator.CurrentConfig.ImageMetaDirectory), 0777)

	for k := range generator.CurrentConfig.ImageSizes {
		os.Mkdir(filepath.Join(generator.CurrentConfig.ImageRootDirectory, k), 0777)
	}

	return nil
}

// WebVars dictate where fotoDen gets its JavaScript and CSS files per page.
type WebVars struct {
	BaseURL       string
	JSLocation    string
	CSSLocation   string
	StaticWebVars map[string]string
}

// NewWebVars creates a WebVars object. Takes a single URL string, and outputs
// a set of fotoDen compatible URLs.
//
// Generates:
// BaseURL: BaseURL
// JSLocation: BaseURL/js/fotoDen.js -- This is the only expected JavaScript file.
// CSSLocation: BaseURL/css/ -- This is expected to be set according to a theme name's CSS
// 								and needs to be processed during configuration.
// 								This is really meant for themes that autoconfigure themselves,
// 								and is mostly an optional path.
//
// Also includes a string map that contains all the static vars that a page can have,
// for when a page is generated to have static parts as well.
//
// If an error occurs, returns an empty WebVars and the error, otherwise returns a filled WebVars.
func NewWebVars(u string) (*WebVars, error) {

	webvars := new(WebVars)
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

	jsurl.Path = path.Join(jsurl.Path, "js", "fotoDen.js")
	webvars.JSLocation = jsurl.String()

	cssurl.Path = path.Join(cssurl.Path, "css", "theme.css")
	webvars.CSSLocation = cssurl.String()

	webvars.StaticWebVars = map[string]string{
		"isStatic": "{{.IsStatic}}",
		"name":     "{{.PageName}}",
		"desc":     "{{.PageDesc}}",
		"auth":     "{{.PageAuthor}}",
		"sfol":     "{{.PageFolder}}",
	}

	return webvars, nil
}

// StaticWebVars are fields that a page can take in order to allow for static page generation.
// If a folder is marked for dynamic generation, these will all automatically be blank.
// Otherwise, these will have the relevant information inside. This only applies to folders.
type StaticWebVars struct {
	IsStatic   bool
	PageName   string // the current name of the page, e.g. 'My album', or 'Photo name'
	PageDesc   string // the current description of the page
	PageFolder string // the folder this is contained in
	PageAuthor string // the author of the page, i.e. the photographer
}

// NewStaticWebVars creates a new static web var set based on the folder given.
// Returns a filled webvar set - save for superFolder, which only occurs
// if the folder above is a fotoDen folder or not.
//
// If an error occurs, it returns a potentially incomplete StaticWebVars with an error.
func NewStaticWebVars(folder string) (*StaticWebVars, error) {
	swebvars := new(StaticWebVars)
	f := new(generator.Folder)
	fpath, _ := filepath.Abs(folder)

	err := f.ReadFolderInfo(filepath.Join(fpath, "folderInfo.json"))
	if err != nil {
		return swebvars, err
	}

	swebvars.IsStatic = true
	swebvars.PageName = f.Name
	swebvars.PageDesc = f.Desc

	superFolder := func() bool {
		_, err := os.Stat(filepath.Join(filepath.Dir(fpath), "folderInfo.json"))
		return os.IsNotExist(err)
	}()

	if !superFolder {
		verbose("Folder above is a fotoDen folder, using that...")
		err = f.ReadFolderInfo(filepath.Join(filepath.Dir(fpath), "folderInfo.json"))
		if err != nil {
			return swebvars, err
		}

		swebvars.PageFolder = f.Name
	}

	return swebvars, nil
}

// ConfigureWebFile configures the web variables in a template by putting it through Go's template system.
// Outputs to a destination location.
// Can be used for any fotoDen-compatible web file.
//
// This should only be done once, ideally,
// and copied over to a configuration directory
// for fotoDen to use (in essence, CurrentConfig.WebSourceDirectory)
func ConfigureWebFile(source string, dest string, vars interface{}) error {
	webpage, err := template.ParseFiles(source)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755) // just rewrite the entire thing
	if err != nil {
		return err
	}
	defer file.Close()

	err = webpage.Execute(file, vars)
	if err != nil {
		return err
	}

	return nil
}

// GenerateWeb hooks into the fotoDen/generator package, and generates web pages for fotoDen folders/albums.
// If a folder was marked as static, it will create a StaticWebVars object to put into the page, and generate
// a more static page - otherwise, it will leave those fields blank, and leave it to the fotoDen front end
// to generate the rest of the page.
func GenerateWeb(m string, dest string, f *generator.Folder, opt GeneratorOptions) error {
	verbose("Generating web pages...")
	var err error
	var pageOptions *StaticWebVars

	if opt.Static || f.Static {
		verbose("Folder/album is static, generating static web vars...")
		pageOptions, err = NewStaticWebVars(dest)
		if checkError(err) {
			return err
		}
	} else {
		verbose("Folder/album is dynamic.")
		pageOptions = new(StaticWebVars)
		pageOptions.IsStatic = false
	}

	switch m {
	case "album":
		verbose("Album mode selected, generating album.")
		err = ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "album-template.html"), path.Join(dest, "index.html"), pageOptions)
		err = ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "photo-template.html"), path.Join(dest, "photo.html"), pageOptions)
	case "folder":
		verbose("Folder mode selected, generating folder.")
		err = ConfigureWebFile(path.Join(generator.CurrentConfig.WebSourceLocation, "html", "folder-template.html"), path.Join(dest, "index.html"), pageOptions)
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
	err = ConfigureWebFile(
		path.Join(generator.CurrentConfig.WebSourceLocation, "html", "page-template.html"),
		dest,
		struct {
			PageContent string
			Title       string
		}{string(r), title},
	)
	if checkError(err) {
		return err
	}

	return nil
}
