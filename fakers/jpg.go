package fakers

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	mediaFakerProcessors "github.com/adrian-gheorghe/mediafaker-processors"
	"github.com/adrian-gheorghe/mediafaker/extpoints"
	log "github.com/sirupsen/logrus"
)

func init() {
	faker := new(Jpg)
	faker.Name = "JPEG"
	faker.Extensions = []string{"jpg", "jpeg"}
	faker.Sizes = make(map[string]BlockSize)
	extpoints.RegisterExtension(faker, "jpg")

}

// Jpg structure
type Jpg struct {
	Name       string
	Extensions []string
	Sizes      map[string]BlockSize
	Quality    int
}

// Init Method
func (faker *Jpg) init() {
	log.Info("Media Faker Registered: ", faker.Name)
}

// GetExtensions Method
func (faker *Jpg) GetExtensions() []string {
	return faker.Extensions
}

// Fake Command
func (faker *Jpg) Fake(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error {
	log.Info(sourcePath)

	imageProcessor := mediaFakerProcessors.ImageProcessor{}
	imageInfo, err := imageProcessor.Inspect(sourcePath)
	if err != nil {
		return err
	}

	width := imageInfo.Width
	height := imageInfo.Height
	outImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < len(imageInfo.PixelInfo); i++ {
		pixelInfo, err := imageProcessor.ExtractRectangleInfo(imageInfo.PixelInfo[i])
		if err != nil {
			log.Println(err)
			return err
		}
		colorInfo, err := imageProcessor.ParseHexColorFast("#" + pixelInfo.Color)
		if err != nil {
			return err
		}
		draw.Draw(outImg, pixelInfo.Rectangle, &image.Uniform{colorInfo}, image.ZP, draw.Src)
	}

	destinationImage, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	err = jpeg.Encode(destinationImage, outImg, &jpeg.Options{Quality: faker.Quality})
	if err != nil {
		return err
	}
	return nil
}

// FakeTreeFile Command
func (faker *Jpg) FakeTreeFile(item extpoints.TreeFile, destinationPath string) error {
	log.Info(item.Path, " - faked with JPG")
	width := item.ImageInfo.Width
	height := item.ImageInfo.Height
	outImg := image.NewRGBA(image.Rect(0, 0, width, height))
	imageProcessor := mediaFakerProcessors.ImageProcessor{}
	rectangles, err := imageProcessor.ExtractPixelInfo(item.ImageInfo.PixelInfo)
	if err != nil {
		return errors.New("Pixel Information could not be extracted correctly")
	}
	for i := 0; i < len(rectangles); i++ {
		rectangle := rectangles[i].Rectangle
		colorInfo, err := imageProcessor.ParseHexColorFast("#" + rectangles[i].Color)
		if err != nil {
			return errors.New("Pixel color information is incorrect")
		}
		draw.Draw(outImg, rectangle, &image.Uniform{colorInfo}, image.ZP, draw.Src)
	}

	destinationImage, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	err = jpeg.Encode(destinationImage, outImg, &jpeg.Options{Quality: 65})
	if err != nil {
		return err
	}

	return nil
}

// BlockSize is the abstraction of a widthxHeight map
type BlockSize struct {
	Width  int
	Height int
}
