package generator

// WebConfig is the structure of the JSON config file that fotoDen uses.
type WebConfig struct {
	WebsiteTitle     string         `json:"websiteTitle"`
	PhotoURLBase     string         `json:"storageURL"`
	ImageRootDir     string         `json:"imageRoot"`
	ThumbnailFrom    string         `json:"thumbnailSize"`
	DisplayImageFrom string         `json:"displayImageSize"`
	Theme            bool           `json:"theme"`
	DownloadSizes    []string       `json:"downloadableSizes"`
	ImageSizes       []WebImageSize `json:"imageSizes"`
}

// WebImageSize is a structure for image size types that fotoDen will call on.
type WebImageSize struct {
	SizeName  string `json:"sizeName"` // the semantic name of the size
	Directory string `json:"dir"`      // the directory the size is stored in, relative to ImageRootDir
	LocalBool bool   `json:"local"`    // whether to download it remotely or locally
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
