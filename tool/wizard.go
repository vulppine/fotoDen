package tool

import (
	"fmt"
	"path/filepath"

	"github.com/vulppine/fotoDen/generator"
)

func isBlank(input string) bool {
	if input == "" {
		return true
	}

	return false
}

func setupWebsite(loc string, theme string) (WebsiteConfig, *generator.WebConfig) {
	fmt.Println("Wizard: Setup fotoDen website")
	fmt.Println("Some of these are required. Do not try to skip them.")
	w := WebsiteConfig{}
	w.RootLocation = loc

	w.Name = ReadInputReq("What is the name of your website? (required)")
	w.URL = ReadInputReq("What is the URL of your website? (required)")
	w.GeneratorConfig = setupConfig()
	w.GeneratorConfig.WebBaseURL = w.URL
	w.GeneratorConfig.WebSourceLocation, _ = filepath.Abs(
		filepath.Join(
			generator.RootConfigDir,
			"sites",
			w.Name,
			"theme",
			theme,
		))

	src := ReadInput("Are you going to remotely host your images? If so, type in the URL now, otherwise leave it blank to automatically use local hosting for all images")
	s := generator.GenerateWebConfig(src)
	s.WebsiteTitle = w.Name
	s.Theme = true

	fmt.Println("Here are your current image sizes, for reference:")
	for k := range generator.CurrentConfig.ImageSizes {
		fmt.Println(k)
	}
	s.ImageRootDir = generator.CurrentConfig.ImageSrcDirectory
	s.ThumbnailFrom = ReadInputReq("What size do you want your thumbnails to be? (required)")
	s.DisplayImageFrom = ReadInputReq("What size do you want to display your images as in a fotoDen photo viewer? (required)")
	s.DownloadSizes = ReadInputAsArray("What sizes do you want easily downloadable?", ",")

	return w, s
}

func setupConfig() generator.Config {
	config := generator.Config{}

	fmt.Println("Wizard: Setup fotoDen config")
	fmt.Println("Leave blank for default!")
	config.ImageRootDirectory = ReadInput("Where do you want your images to be stored in every folder?")
	if isBlank(config.ImageRootDirectory) {
		config.ImageRootDirectory = generator.DefaultConfig.ImageRootDirectory
	}

	config.ImageSrcDirectory = ReadInput("Where are you going to store your source-quality images?")
	if isBlank(config.ImageSrcDirectory) {
		config.ImageSrcDirectory = generator.DefaultConfig.ImageSrcDirectory
	}

	config.ImageSizes = map[string]generator.ImageScale{}
	imageSizes := ReadInputAsArray("What image sizes do you want? Separate by comma, no spaces.", ",")
	if imageSizes[0] != "" {
		fmt.Printf("Image sizes detected: %v\n", imageSizes)
		fmt.Println("Leave blank to set as zero. At least one value must be filled in. Priority: MaxHeight, MaxWidth, ScalePercent")
		for _, val := range imageSizes {
			var c bool
			for c != true {
				fmt.Println("Image size " + val)
				imageSize := generator.ImageScale{}
				imageSize.MaxHeight, _ = ReadInputAsInt("Maximium height of image?")
				imageSize.MaxWidth, _ = ReadInputAsInt("Maximum width of image?")
				scalePercent, _ := ReadInputAsFloat("Image scale percent? [0 - 100%]")
				imageSize.ScalePercent = scalePercent * 0.1
				if imageSize.MaxHeight == 0 && imageSize.MaxWidth == 0 && imageSize.ScalePercent == 0 {
					fmt.Println("You must set one value to a non-zero value!")
				} else {
					c = true
					config.ImageSizes[val] = imageSize
				}
			}
		}
	} else {
		config.ImageSizes = generator.DefaultConfig.ImageSizes
		fmt.Println("Using default image sizes: ")
		for k, v := range generator.DefaultConfig.ImageSizes {
			fmt.Printf("Size name: %s, MaxHeight: %d, MaxWidth: %d, ScalePercent: %f\n", k, v.MaxHeight, v.MaxWidth, v.ScalePercent)
		}
	}

	return config
}

func setupWebConfig(source string) *generator.WebConfig {
	config := generator.GenerateWebConfig(source)

	fmt.Println("Wizard: Setup website config")
	config.WebsiteTitle = ReadInput("What is the title of your website?")
	config.PhotoURLBase = ReadInput("Are you going to be using a remote storage provider for your photos? If so, put in the URL to the folder containing your fotoDen-structured images here.")

	fmt.Println("Here are your current image sizes, for reference:")
	for k := range generator.CurrentConfig.ImageSizes {
		fmt.Println(k)
	}

	config.Theme = true // TODO: Make selectable themes
	config.ImageRootDir = generator.CurrentConfig.ImageSrcDirectory
	config.ThumbnailFrom = ReadInput("What size do you want your thumbnails to be?")
	config.DisplayImageFrom = ReadInput("What size do you want to display your images as in a fotoDen photo viewer?")
	config.DownloadSizes = ReadInputAsArray("What sizes do you want downlodable?", ",")

	return config
}

func generateFolderWizard(directory string) (*generator.Folder, error) {
	folder, err := generator.GenerateFolderInfo(directory, "")
	if checkError(err) {
		return nil, err
	}

	fmt.Println("Wizard: Generate folder")
	folder.Name = ReadInput("What is the name of this folder/album?")
	folder.Desc = ReadInput("What is the description of this folder/album?")
	folder.Static = ReadInputAsBool("Will the folder/album webpages have some dynamic elements static?", "y")
	t := ReadInputAsBool("Will the folder have a thumbnail?", "y")
	if t {
		thumb := ReadInput("Where is the thumbnail located? (direct path or relative to current working directory)")
		if thumb == "" {
			fmt.Println("No file detected - ignoring.")
		} else {
			folder.Thumbnail = true
			generator.MakeFolderThumbnail(thumb, directory)
		}
	}

	return folder, nil
}

func updateFolderWizard(fol *generator.Folder) *generator.Folder {
	fol.Name = ReadInput("What is the new name of this folder/album?")
	fol.Desc = ReadInput("What is the new description of this folder/album?")

	return fol
}
