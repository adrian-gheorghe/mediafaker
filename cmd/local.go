package cmd

import (
	"github.com/adrian-gheorghe/mediafaker/fakers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SourcePath represents the local directory path that should be faked
var SourcePath string

func init() {
	RootCmd.AddCommand(Local)
	Local.Flags().StringVarP(&SourcePath, "source", "s", "", "Local Source Directory path which should be faked")
}

// Local represents the tree tree command
var Local = &cobra.Command{
	Use:   "local",
	Short: "runs mediafaker on a local directory",
	Long:  `This subcommand runs mediafaker on a local source directory and stores the files in a destination directory of your choice.`,
	Run: func(cmd *cobra.Command, args []string) {
		if SourcePath == "" {
			log.Fatal("Source directory has not been provided")
			return
		}

		if DestinationPath == "" {
			log.Fatal("Destination directory has not been provided")
			return
		}

		mediaFake := fakers.MediaFake{
			SourcePath:                    SourcePath,
			MoniTreePath:                  "",
			DestinationPath:               DestinationPath,
			ExtensionsToCopyAutomatically: ExtensionsToCopyAutomatically,
			MaximumSizeForCopy:            MaximumSizeForCopy,
			TotalFaked:                    0,
			TotalMissed:                   0,
			TotalCopied:                   0,
			TotalSourceSize:               0,
			TotalDestinationSize:          0,
		}
		error := mediaFake.Fake()
		if error != nil {
			log.Error(error)
		}

		log.Info("Media directory has been faked:", SourcePath)
		mediaFake.CalculateTotalDestinationSize()
		mediaFake.PrintInfo()
	},
}
