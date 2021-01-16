package generator

import (
	"fmt"
	"github.com/h2non/bimg"
	"path"
	"strconv"
)

// IsolateImages
//
// Isolates images in an array.
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

// ResizeImage
//
// Resizes a single image.
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

type ImageScale struct {
	maxheight    int
	maxwidth     int
	scalepercent float32
}

func ResizeImage(imageName string, newName string, scale ImageScale, dest string, imageFormat bimg.ImageType) error {
	image, err := bimg.Read(imageName)
	if err != nil {
		return err
	}

	imageType := bimg.DetermineImageType(image)
	if imageType == bimg.UNKNOWN {
		return fmt.Errorf("ResizeImage: Unknown file type. Skipping. Image: " + imageName)
	}

	newImage, err := bimg.NewImage(image).Convert(imageFormat)
	if err != nil {
		return err
	}

	size, err := bimg.NewImage(image).Size()
	width := float32(size.Width)
	height := float32(size.Height)

	switch {
	case scale.maxheight != 0:
		scale := float32(scale.maxheight) / height
		width = width * scale
		height = height * scale
	case scale.maxwidth != 0:
		scale := float32(scale.maxwidth) / width
		width = width * scale
		height = height * scale
	case scale.scalepercent != 0:
		width = width * scale.scalepercent
		height = height * scale.scalepercent
	default:
		return fmt.Errorf("ResizeImage: Image scaling undefined. Aborting. scale: " + fmt.Sprint(scale))
	}

	verbose("Resizing " + imageName + " to " + strconv.Itoa(int(width)) + "," + strconv.Itoa(int(height)) + " and attempting to place it in " + path.Join(dest, newName))
	newImage, err = bimg.NewImage(newImage).Resize(int(width), int(height))
	if err != nil {
		return err
	}

	destCheck, err := bimg.Read(path.Join(dest, newName))
	if destCheck != nil {
		return err
	}

	bimg.Write(path.Join(dest, newName), newImage)
	return nil
}

// MakeThumbnail
//
// This is only here to make fotoDen's command line tool look cleaner in code, and avoid importing more than needed.

func MakeFolderThumbnail(file string, directory string) error {
	err := ResizeImage(file, "thumb.jpg", ThumbScalingOptions, directory, bimg.JPEG)
	if err != nil {
		return err
	}

	return nil
}
