package generator

import (
	"net/url"
	"os"
	"path"
	"text/template"
)

// WebConfig
//
// The structure of the JSON config file that fotoDen uses.
type WebConfig struct {
	WebsiteTitle     string
	WorkingDirectory string
	PhotoURLBase     string
	ImageRootDir     string
	ThumbnailFrom    string
	DisplayImageFrom string
	ImageSizes		 []WebImageSize
}

// WebImageSize
//
// A structure for image size types that fotoDen will call on.
type WebImageSize struct {
	SizeName string			// the semantic name of the size
	Directory string		// the directory the size is stored in, relative to ImageRootDir
	LocalBool bool          // whether to download it remotely or locally
}

// WebVars
//
// These will dictate where fotoDen gets its JavaScript and CSS files per page.
type WebVars struct {
	BaseURL     string
	JSLocation  string
	CSSLocation string
}

// GenerateWebConfig
//
// Creates a new WebConfig object, and returns a WebConfig object with a populated ImageSizes
// based on the current ScalingOptions map.
func GenerateWebConfig(source string) *WebConfig {

	webconfig := new(WebConfig)
	webconfig.PhotoURLBase = source

	for k, _ := range CurrentConfig.ImageSizes {
		webconfig.ImageSizes = append(
			webconfig.ImageSizes,
			WebImageSize{
				SizeName: k,
				Directory: k,
				LocalBool: true,
			},
		)
	}

	return webconfig
}


func (config *WebConfig) ReadWebConfig(filepath string) error {
	err := ReadJSON(filepath, config)
	if err != nil {
		return err
	}

	return nil
}

func (config *WebConfig) WriteWebConfig(filepath string) error {
	err := WriteJSON(filepath, "multi", config)
	if err != nil {
		return err
	}

	return nil
}

// GenerateWebRoot
//
// Generates the root of a fotoDen website in filepath.
// Creates the folders for JS and CSS placement.
//
// It is up to the fotoDen tool to copy over the relevant files,
// and folder configuration.
func GenerateWebRoot(filepath string) error {
	err := os.Mkdir(filepath, 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(filepath, "js"), 0755)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(filepath, "css"), 0755)
	if err != nil {
		return err
	}

	return nil
}

// MakeAlbumDirectoryStructure
//
// Makes a fotoDen-suitable album structure in the given rootDirectory (string).
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

	for k, _ := range CurrentConfig.ImageSizes {
		os.Mkdir(path.Join(CurrentConfig.ImageRootDirectory, k), 0777)
	}

	return nil
}

// NewWebVars
//
// Creates a WebVars object. Takes a single URL string, and outputs
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

	jsurl.Path = path.Join(jsurl.Path, "js", "fotoDen.js")
	webvars.JSLocation = jsurl.String()

	cssurl.Path = path.Join(cssurl.Path, "css", "theme.css")
	webvars.CSSLocation = cssurl.String()

	return webvars, nil
}


// ConfigureWebFile
//
// Configures the web variables in a template by putting it through Go's template system.
// Outputs to a destination location.
// Can be used for any fotoDen-compatible web file.
//
// This should only be done once, ideally,
// and copied over to a configuration directory
// for fotoDen to use (in essence, CurrentConfig.WebSourceDirectory)
func ConfigureWebFile(source string, dest string, config *WebVars) error {
	webpage, err := template.ParseFiles(source)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755) // just rewrite the entire thing
	if err != nil {
		return err
	}
	defer file.Close()

	err = webpage.Execute(file, config)
	if err != nil {
		panic(err)
	} // means something wrong happened during file write

	return nil
}
