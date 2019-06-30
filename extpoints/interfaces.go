package extpoints

import (
	"os"
)

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Path      string     `json:"Path"`
	Type      string     `json:"Type"`
	Mode      string     `json:"Mode"`
	Size      int64      `json:"Size"`
	Modtime   string     `json:"Modtime"`
	Sum       string     `json:"Sum"`
	MediaType string     `json:MediaType`
	Content   string     `json:"Content"`
	ImageInfo ImageInfo  `json:"ImageInfo"`
	Children  []TreeFile `json:"Children"`
}

type ImageInfo struct {
	Width       int    `json:"W"`
	Height      int    `json:"H"`
	PixelInfo   string `json:"P"`
	BlockWidth  int    `json:"BW"`
	BlockHeight int    `json:"BH"`
}

// MediaFakerType is the abstraction of all the structs that fake file types
type MediaFakerType interface {
	Fake(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error
	FakeTreeFile(source TreeFile, destinationPath string) error
	GetExtensions() []string
}
