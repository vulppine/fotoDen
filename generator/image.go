package generator

import (
	"fmt"
	"github.com/h2non/bimg"
	"path"
	"strconv"
)

// IsolateImages isolates images in an array.
//
// Checks all image files at O(n), if a file is not an image, removes it from the current slice.
func IsolateImages(files []string) []string {
	for i := 0; i < len(files); i++ {
		image, err := bimg.Read(files[i])
		if err != nil {
			fmt.Println(err)
		} else {
			if bimg.DetermineImageType(image) == bimg.UNKNOWN {
				verbose("File " + files[i] + " is not an image. Removing.")
				files = RemoveItemFromStringArray(files, files[i]) // replace this with an append later
				i--                                                // because now everything is shifted one backwards
			}
		}
	}

	return files
}

// ImageScale represents scaling options to be used by ResizeImage.
// See ResizeImage for more information.
type ImageScale struct {
	MaxHeight    int
	MaxWidth     int
	ScalePercent float32
}

// ResizeImage resizes a single image.
//
// You'll have to pass it a ImageScale object,
// which contains values for either a scale percentage, or a max height/width.
// In order of usage:
// maxheight, maxwidth, scalepercent
//
// maxheight is first due to a restriction with CSS Flex and mixed height images,
// maxwidth is second for the same thing, but with flex set to column mode,
// scalepercent is final for when the first two don't apply.
//
// The function will output the image to the given directory, without changing the name.
// It will return an error if the filename given already exists in the destination directory.
func ResizeImage(file string, imageName string, scale ImageScale, dest string, imageFormat bimg.ImageType) error {
	image, err := bimg.Read(file)
	if err != nil {
		return err
	}

	imageType := bimg.DetermineImageType(image)
	if imageType == bimg.UNKNOWN {
		return fmt.Errorf("ResizeImage: Unknown file type. Skipping. Image: %s", imageName)
	}

	newImage, err := bimg.NewImage(image).Convert(imageFormat)
	if err != nil {
		return err
	}

	size, err := bimg.NewImage(image).Size()
	width := float32(size.Width)
	height := float32(size.Height)

	switch {
	case scale.MaxHeight != 0:
		scale := float32(scale.MaxHeight) / height
		width = width * scale
		height = height * scale
	case scale.MaxWidth != 0:
		scale := float32(scale.MaxWidth) / width
		width = width * scale
		height = height * scale
	case scale.ScalePercent != 0:
		width = width * scale.ScalePercent
		height = height * scale.ScalePercent
	default:
		return fmt.Errorf("ResizeImage: Image scaling undefined. Aborting. scale: " + fmt.Sprint(scale))
	}

	verbose("Resizing " + imageName + " to " + strconv.Itoa(int(width)) + "," + strconv.Itoa(int(height)) + " and attempting to place it in " + path.Join(dest, imageName))
	newImage, err = bimg.NewImage(newImage).Resize(int(width), int(height))
	if err != nil {
		return err
	}

	destCheck, err := bimg.Read(path.Join(dest, imageName))
	if destCheck != nil {
		return err
	}

	bimg.Write(path.Join(dest, imageName), newImage)
	return nil
}

// ImageMeta provides metadata such as name and description for an image file.
// This is used by fotoDen on the web frontend to display custom per-image information.
type ImageMeta struct {
	ImageName string // The name of an image.
	ImageDesc string // The description of an image.
}

// WriteImageMeta writes an ImageMeta struct to a file.
//
// Takes two arguments: a folder destination, and a name.
// The name is automatically combined to create a [name].json file,
// in order to ensure compatibility with fotoDen.
// Writes the json file into the given folder.
func (meta *ImageMeta) WriteImageMeta(folder string, name string) error {
	err := WriteJSON(path.Join(folder, name+".json"), "multi", meta)
	if err != nil {
		return err
	}

	return nil
}

// MakeFolderThumbnail creates a thumbnail from a file into a destination directory.
// This is only here to make fotoDen's command line tool look cleaner in code, and avoid importing more than needed.
func MakeFolderThumbnail(file string, directory string) error {
	err := ResizeImage(file, "thumb.jpg", CurrentConfig.ImageSizes["thumb"], directory, bimg.JPEG)
	if err != nil {
		return err
	}

	return nil
}
