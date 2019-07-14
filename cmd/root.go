package cmd

import (
	"fmt"
	"os"

	figure "github.com/common-nighthawk/go-figure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var AppVersion = "0.1.2"

var RootCmd = &cobra.Command{
	Use:     "mediafaker",
	Short:   "fake your media files for local use",
	Long:    `This utility creates a fake simplified version of a directory tree of your chosing, making it easier to work locally on legacy projects that have a large media asset folder.`,
	Version: AppVersion,
}
var (
	jsonlog bool
)

// DestinationPath represents the local directory path where the new faked files should be stored
var DestinationPath string

// MaximumSizeForCopy represents the maximum size a file should have to be copied automatically
var MaximumSizeForCopy int64

// ExtensionsToCopyAutomatically filters files by extension and copies automatically instead of faking them
var ExtensionsToCopyAutomatically []string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once in the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&DestinationPath, "destination", "d", "", "Local Destination directory path where mediafaker should store the files")
	RootCmd.PersistentFlags().StringSliceVarP(&ExtensionsToCopyAutomatically, "extcopy", "e", []string{}, "List of extensions that should be copied automatically")
	RootCmd.PersistentFlags().Int64VarP(&MaximumSizeForCopy, "maxcopy", "m", 30000, "Maximum Size(in bytes) a file should have to be copied automatically if it cannot be faked")
	RootCmd.PersistentFlags().BoolVarP(&jsonlog, "jsonlog", "j", false, "Change logger format to json")
}

func initConfig() {
	myFigure := figure.NewFigure("MediaFaker", "", true)
	myFigure.Print()
	fmt.Println("")

	if jsonlog {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetOutput(os.Stdout)
	}
}
