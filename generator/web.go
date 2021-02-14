package generator

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

// WebConfig is the structure of the JSON config file that fotoDen uses.
type WebConfig struct {
	WebsiteTitle     string         `json:websiteTitle`
	PhotoURLBase     string         `json:storageURL`
	ImageRootDir     string         `json:imageRoot`
	ThumbnailFrom    string         `json:thumbnailSize`
	DisplayImageFrom string         `json:displayImageSize`
	Theme            bool           `json:theme`
	DownloadSizes    []string       `json:downloadableSizes`
	ImageSizes       []WebImageSize `json:imageSizes`
}

// WebImageSize is a structure for image size types that fotoDen will call on.
type WebImageSize struct {
	SizeName  string `json:sizeName` // the semantic name of the size
	Directory string `json:dir`      // the directory the size is stored in, relative to ImageRootDir
	LocalBool bool   `json:local`    // whether to download it remotely or locally
}

// GenerateWebConfig creates a new WebConfig object, and returns a WebConfig object with a populated ImageSizes
// based on the current ScalingOptions map.
func GenerateWebConfig(source string) *WebConfig {

	webconfig := new(WebConfig)
	webconfig.PhotoURLBase = source

	for k := range CurrentConfig.ImageSizes {
		webconfig.ImageSizes = append(
			webconfig.ImageSizes,
			WebImageSize{
				SizeName:  k,
				Directory: k,
				LocalBool: true,
			},
		)
	}

	return webconfig
}

// ReadWebConfig reads a JSON file containing WebConfig fields into a WebConfig struct.
func (config *WebConfig) ReadWebConfig(fpath string) error {
	err := ReadJSON(fpath, config)
	if err != nil {
		return err
	}

	return nil
}

// WriteWebConfig writes a WebConfig struct into the specified path.
func (config *WebConfig) WriteWebConfig(fpath string) error {
	err := WriteJSON(fpath, "multi", config)
	if err != nil {
		return err
	}

	return nil
}

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
	os.Mkdir(CurrentConfig.ImageRootDirectory, 0777)
	os.Mkdir(path.Join(CurrentConfig.ImageRootDirectory, CurrentConfig.ImageSrcDirectory), 0777)
	os.Mkdir(path.Join(CurrentConfig.ImageRootDirectory, CurrentConfig.ImageMetaDirectory), 0777)

	for k := range CurrentConfig.ImageSizes {
		os.Mkdir(filepath.Join(CurrentConfig.ImageRootDirectory, k), 0777)
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
	f := new(Folder)
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
