package tool

import (
	"fmt"

	"github.com/vulppine/fotoDen/generator"
)

func isBlank(input string) bool {
	if input == "" {
		return true
	}

	return false
}

func SetupConfig() generator.GeneratorConfig {
	config := generator.GeneratorConfig{}

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
			fmt.Println("Image size " + val)
			imageSize := generator.ImageScale{}
			imageSize.MaxHeight, _ = ReadInputAsInt("Maximium height of image?")
			imageSize.MaxWidth, _ = ReadInputAsInt("Maximum width of image?")
			scalePercent, _ := ReadInputAsFloat("Image scale percent? [0 - 100%]")
			imageSize.ScalePercent = scalePercent * 0.1
			config.ImageSizes[val] = imageSize
		}
	} else {
		config.ImageSizes = generator.DefaultConfig.ImageSizes
		fmt.Println("Using default image sizes: ")
		for k, v := range generator.DefaultConfig.ImageSizes {
			fmt.Printf("Size name: %s, MaxHeight: %d, MaxWidth: %d, ScalePercent: %f\n", k, v.MaxHeight, v.MaxWidth, v.ScalePercent)
		}
	}

	config.WebBaseURL = ReadInput("What is the URL of your website?")

	return config
}

func SetupWebConfig(source string) *generator.WebConfig {
	config := generator.GenerateWebConfig(source)

	fmt.Println("Wizard: Setup website config")
	config.WebsiteTitle = ReadInput("What is the title of your website?")
	config.PhotoURLBase = ReadInput("Are you going to be using a remote storage provider for your photos? If so, put in the URL to the folder containing your fotoDen-structured images here.")

	fmt.Println("Here are your current image sizes, for reference:")
	for k := range generator.CurrentConfig.ImageSizes {
		fmt.Println(k)
	}

	config.ThumbnailFrom = ReadInput("What size do you want your thumbnails to be?")
	config.DisplayImageFrom = ReadInput("What size do you want to display your images as in a fotoDen photo viewer?")

	return config
}
