package extpoints

import "os"

// MediaFakerType is the abstraction of all the structs that fake file types
type MediaFakerType interface {
	Fake(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error
	GetExtensions() []string
}
