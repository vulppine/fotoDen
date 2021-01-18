package generator

import (
	"encoding/json"
	"fmt"
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
	PhotoExtension   string
	DisplayImageFrom string
	DownloadImageFrom string
	ImageThumbDir    string
	ImageLargeDir    string
	ImageSrcDir      string
}

// WebVars
//
// These will dictate where fotoDen gets its JavaScript and CSS files per page.

type WebVars struct {
	BaseURL     string
	JSLocation  string
	CSSLocation string
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

// NewWebConfig
//
// Creates a new WebConfig object.
// Takes the a string containing the root URL of where photos are stored.
// Returns only a WebConfig object.

func NewWebConfig(rooturl string) *WebConfig {

	webconfig := new(WebConfig)
	webconfig.DisplayImageFrom = "large" // default, for saving bandwidth
	webconfig.DownloadImageFrom = "src"
	webconfig.PhotoURLBase = rooturl

	return webconfig
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

// GenerateWebConfig
//
// Generates a new web configuration file in the given path.
// This is required for fotoDen operation.

func (config *WebConfig) GenerateWebConfig(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileinfo, _ := file.Stat()
	if fileinfo.Size() > 0 {
		return fmt.Errorf("GenerateWebConfig: file already exists")
	}

	configjson, _ := json.Marshal(config)
	_, err = file.Write(configjson)
	if err != nil {
		return err
	}

	return nil
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
