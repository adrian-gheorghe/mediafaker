package fakers

import (
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

	blockSize, ok := faker.Sizes[string(width)+"x"+string(height)]
	if !ok {
		blockSize = BlockSize{imageInfo.BlockWidth, imageInfo.BlockHeight}
		faker.Sizes[string(width)+"x"+string(height)] = blockSize
	}

	for i := 0; i < len(imageInfo.PixelInfo); i++ {
		draw.Draw(outImg, imageInfo.PixelInfo[i].Rectangle, &image.Uniform{imageInfo.PixelInfo[i].Color}, image.ZP, draw.Src)
	}

	destinationImage, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	err = jpeg.Encode(destinationImage, outImg, &jpeg.Options{Quality: 90})
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
