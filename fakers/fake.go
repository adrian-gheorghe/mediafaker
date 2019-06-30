package fakers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	imageProcessors "github.com/adrian-gheorghe/mediafaker-processors"
	"github.com/adrian-gheorghe/mediafaker/extpoints"
	log "github.com/sirupsen/logrus"
)

const (
	tmpPermissionForDirectory = os.FileMode(0755)
)

var (
	mediaFakerTypes = extpoints.MediaFakerTypes
)

// MediaFake structure
type MediaFake struct {
	SourcePath                    string
	MoniTreePath                  string
	DestinationPath               string
	ExtensionsToCopyAutomatically []string
	MaximumSizeForCopy            int64
}

// Fake Command
func (mediaFake *MediaFake) Fake() error {
	if mediaFake.SourcePath != "" {
		sourceInfo, err := os.Lstat(mediaFake.SourcePath)
		if err != nil {
			return err
		}
		return mediaFake.FakeItem(mediaFake.SourcePath, mediaFake.DestinationPath, sourceInfo)
	} else {
		log.Panic(mediaFake)
		return nil
		tree, err := mediaFake.GetTree()
		if err != nil {
			return err
		}

		return mediaFake.FakeTree(tree)
	}
	return nil
}

// FakeItem Command
func (mediaFake *MediaFake) FakeItem(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error {
	if sourceInfo.IsDir() {
		return mediaFake.FakeDir(sourcePath, destinationPath, sourceInfo)
	} else if sourceInfo.Mode()&os.ModeSymlink != 0 {
		return mediaFake.FakeLink(sourcePath, destinationPath, sourceInfo)
	}
	return mediaFake.FakeFile(sourcePath, strings.ToLower(destinationPath), sourceInfo)
}

// FakeFile method
func (mediaFake *MediaFake) FakeFile(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error {
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(sourcePath)), ".")
	if StringInSlice(extension, mediaFake.ExtensionsToCopyAutomatically) {
		return mediaFake.Copy(sourcePath, destinationPath, sourceInfo)
	}
	fakerFound := false
	for _, faker := range mediaFakerTypes.All() {
		if StringInSlice(extension, faker.GetExtensions()) {
			fakerFound = true
			err := faker.Fake(sourcePath, destinationPath, sourceInfo)
			if err != nil {
				return err
			}
			break
		}
	}

	if !fakerFound {
		if sourceInfo.Size() < mediaFake.MaximumSizeForCopy {
			return mediaFake.Copy(sourcePath, destinationPath, sourceInfo)
		}

		log.Warn(sourcePath, "Has could not be faked")
	}

	return nil
}

// FakeDir walks through the directory and processes each file
func (mediaFake *MediaFake) FakeDir(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error {
	originalMode := sourceInfo.Mode()
	if error := os.MkdirAll(destinationPath, tmpPermissionForDirectory); error != nil {
		return error
	}
	// Recover dir mode with original one.
	defer os.Chmod(destinationPath, originalMode)

	contents, error := ioutil.ReadDir(sourcePath)
	if error != nil {
		return error
	}

	for _, content := range contents {
		cs, cd := filepath.Join(sourcePath, content.Name()), filepath.Join(destinationPath, content.Name())
		if error := mediaFake.FakeItem(cs, cd, content); error != nil {
			// If any error, exit immediately
			return error
		}
	}
	return nil
}

// FakeLink creates a symlink fake
func (mediaFake *MediaFake) FakeLink(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}

// Copy duplicates the file from source to destination without processing it
func (mediaFake *MediaFake) Copy(src, dest string, info os.FileInfo) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, input, info.Mode())
	if err != nil {
		return err
	}
	return nil
}

// GetTree is the implementation of the tree compare method
func (mediaFake *MediaFake) GetTree() (TreeFile, error) {
	content, err := ioutil.ReadFile(mediaFake.MoniTreePath)
	if err != nil {
		return TreeFile{}, err
	}
	tree := TreeFile{}
	err = json.Unmarshal(content, &tree)
	if err != nil {
		log.Println(err)
	}
	return tree, err
}

// FakeTree Command
func (mediaFake *MediaFake) FakeTree(tree TreeFile) error {
	if tree.Children != nil {
		for _, item := range tree.Children {
			if tree.Type == "directory" {
				mediaFake.FakeDirTree(item)

			} else if tree.Type == "file" {

			}
		}
	}

	return nil
}

// FakeDirTree walks through the directory and processes each file
func (mediaFake *MediaFake) FakeDirTree(item TreeFile) error {
	// originalMode := item.Mode
	// if error := os.MkdirAll(destinationPath, tmpPermissionForDirectory); error != nil {
	// 	return error
	// }
	// // Recover dir mode with original one.
	// defer os.Chmod(destinationPath, originalMode)

	// contents, error := ioutil.ReadDir(sourcePath)
	// if error != nil {
	// 	return error
	// }

	// for _, content := range contents {
	// 	cs, cd := filepath.Join(sourcePath, content.Name()), filepath.Join(destinationPath, content.Name())
	// 	if error := mediaFake.FakeItem(cs, cd, content); error != nil {
	// 		// If any error, exit immediately
	// 		return error
	// 	}
	// }
	return nil
}

// StringInSlice checks if a string is present in a slice of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// TreeFile is a representation of a file or folder in the filesystem
type TreeFile struct {
	Path      string                    `json:"Path"`
	Type      string                    `json:"Type"`
	Mode      string                    `json:"Mode"`
	Size      int64                     `json:"Size"`
	Modtime   string                    `json:"Modtime"`
	Sum       string                    `json:"Sum"`
	MediaType string                    `json:MediaType`
	Content   string                    `json:"Content"`
	ImageInfo imageProcessors.ImageInfo `json:"ImageInfo"`
	Children  []TreeFile                `json:"Children"`
}
