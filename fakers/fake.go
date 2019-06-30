package fakers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrian-gheorghe/mediafaker/extpoints"
	"github.com/jedib0t/go-pretty/table"
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
	OptimizeImageQuality          bool
	TotalFaked                    int
	TotalMissed                   int
	TotalCopied                   int
	TotalSourceSize               int64
	TotalDestinationSize          int64
	TimeStart                     time.Time
	TimeElapsed                   time.Duration
}

// Fake Command
func (mediaFake *MediaFake) Fake() error {
	mediaFake.TimeStart = time.Now()
	log.Info("Faking ", mediaFake.SourcePath, mediaFake.MoniTreePath, " ...")
	if mediaFake.SourcePath != "" {
		sourceInfo, err := os.Lstat(mediaFake.SourcePath)
		if err != nil {
			mediaFake.TimeElapsed = time.Since(mediaFake.TimeStart)
			return err
		}
		err = mediaFake.FakeItem(mediaFake.SourcePath, mediaFake.DestinationPath, sourceInfo)
		if err != nil {
			mediaFake.TimeElapsed = time.Since(mediaFake.TimeStart)
			return err
		}
	} else {
		tree, err := mediaFake.GetTree()
		if err != nil {
			mediaFake.TimeElapsed = time.Since(mediaFake.TimeStart)
			return err
		}
		err = mediaFake.FakeTree(tree)
		if err != nil {
			mediaFake.TimeElapsed = time.Since(mediaFake.TimeStart)
			return err
		}
	}
	mediaFake.TimeElapsed = time.Since(mediaFake.TimeStart)
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

/**
 * FakeFile creates the fake version of the local source file
 */
func (mediaFake *MediaFake) FakeFile(sourcePath string, destinationPath string, sourceInfo os.FileInfo) error {
	mediaFake.TotalSourceSize = mediaFake.TotalSourceSize + sourceInfo.Size()
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
				mediaFake.TotalMissed++
				log.Error("Could not fake: ", err, sourcePath)
			}
			mediaFake.TotalFaked++
			break
		}
	}

	if !fakerFound {
		if sourceInfo.Size() < mediaFake.MaximumSizeForCopy {
			return mediaFake.Copy(sourcePath, destinationPath, sourceInfo)
		}

		mediaFake.TotalMissed++
		log.Warn(sourcePath, " could not be faked")
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
		mediaFake.TotalMissed++
		return err
	}
	mediaFake.TotalCopied++
	return os.Symlink(src, dest)
}

// Copy duplicates the file from source to destination without processing it
func (mediaFake *MediaFake) Copy(src, dest string, info os.FileInfo) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		mediaFake.TotalMissed++
		log.Error("Could not copy: ", err, src)
		return err
	}

	err = ioutil.WriteFile(dest, input, info.Mode())
	if err != nil {
		mediaFake.TotalMissed++
		return err
	}
	mediaFake.TotalCopied++
	return nil
}

// GetTree is the implementation of the tree compare method
func (mediaFake *MediaFake) GetTree() (extpoints.TreeFile, error) {
	content, err := ioutil.ReadFile(mediaFake.MoniTreePath)
	if err != nil {
		return extpoints.TreeFile{}, err
	}
	tree := extpoints.TreeFile{}
	err = json.Unmarshal(content, &tree)
	if err != nil {
		log.Println(err)
	}
	return tree, err
}

// FakeTree Command
func (mediaFake *MediaFake) FakeTree(tree extpoints.TreeFile) error {
	if tree.Children != nil {
		for _, item := range tree.Children {
			if item.Type == "directory" {
				mediaFake.FakeTreeDir(item)
			} else if item.Type == "file" {
				mediaFake.FakeTreeFile(item)
			}
		}
	}
	mediaFake.TimeElapsed = time.Since(mediaFake.TimeStart)

	return nil
}

// FakeTreeDir walks through the directory and processes each file
func (mediaFake *MediaFake) FakeTreeDir(item extpoints.TreeFile) error {
	//originalMode := item.Mode
	destinationPath := path.Join(mediaFake.DestinationPath, item.Path)
	if error := os.MkdirAll(destinationPath, tmpPermissionForDirectory); error != nil {
		return error
	}
	// Recover dir mode with original one.
	//defer os.Chmod(destinationPath, os.FileMode(originalMode))
	return nil
}

// FakeTreeFile method
func (mediaFake *MediaFake) FakeTreeFile(item extpoints.TreeFile) error {
	mediaFake.TotalSourceSize = mediaFake.TotalSourceSize + item.Size
	destinationPath := path.Join(mediaFake.DestinationPath, item.Path)
	extension := strings.TrimPrefix(strings.ToLower(path.Ext(item.Path)), ".")
	fakerFound := false
	for _, faker := range mediaFakerTypes.All() {
		if StringInSlice(extension, faker.GetExtensions()) {
			fakerFound = true
			err := faker.FakeTreeFile(item, destinationPath)
			if err != nil {
				mediaFake.TotalMissed++
				log.Error(item.Path, " - could not be faked.")
				return err
			}
			mediaFake.TotalFaked++
			break
		}
	}

	if !fakerFound {
		if item.Content != "" {
			// create file from content
			data, err := base64.StdEncoding.DecodeString(item.Content)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(destinationPath, data, 0644)
			if err != nil {
				return err
			}
			mediaFake.TotalCopied++
			log.Info(item.Path, " - faked with content from json tree.")
		} else {
			err := ioutil.WriteFile(destinationPath, []byte(""), 0644)
			if err != nil {
				return err
			}
			mediaFake.TotalMissed++
			log.Warn(item.Path, " - faked with empty content.")
		}

		return nil
	}

	return nil
}

// PrintInfo Print info
func (mediaFake *MediaFake) PrintInfo() {
	fmt.Println()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Stats", "Value"})

	t.AppendRows([]table.Row{
		{"Fake Duration", math.Round(mediaFake.TimeElapsed.Seconds())},
	})
	t.AppendRows([]table.Row{
		{"Total Faked", mediaFake.TotalFaked},
	})
	t.AppendRows([]table.Row{
		{"Total Missed", mediaFake.TotalMissed},
	})
	t.AppendRows([]table.Row{
		{"Total Copied", mediaFake.TotalCopied},
	})
	t.AppendRows([]table.Row{
		{"Total Source Size", ByteCountSI(mediaFake.TotalSourceSize)},
	})
	t.AppendRows([]table.Row{
		{"Total Destination Size", ByteCountSI(mediaFake.TotalDestinationSize)},
	})
	t.Render()
}

// CalculateTotalDestinationSize Gets the total destination size
func (mediaFake *MediaFake) CalculateTotalDestinationSize() {
	log.Info("Computing destination size ...")
	size, err := TotalDirSize(mediaFake.DestinationPath)
	if err != nil {
		log.Error("Could not compute destination size: ", err)
	}
	mediaFake.TotalDestinationSize = size
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

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func TotalDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
