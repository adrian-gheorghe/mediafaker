package cmd

import (
	"github.com/spf13/cobra"
)

// SourcePath represents the local directory path that should be faked
var SourcePath string

func init() {
	RootCmd.AddCommand(local)
	local.LocalFlags().StringVarP(&SourcePath, "source", "s", "", "Source Directory path which should be faked")
}

// local represents the tree tree command
var local = &cobra.Command{
	Use:   "local",
	Short: "runs mediafaker on a local directory",
	Long:  `This subcommand runs mediafaker on a local source directory and stores the files in a destination directory of your choice.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if SourcePath == "" {
		// 	log.Fatal("Source directory has not been provided")
		// 	os.Exit(1)
		// }

		// if DestinationPath == "" {
		// 	log.Fatal("Destination directory has not been provided")
		// 	os.Exit(1)
		// }

		// mediaFake := fakers.MediaFake{
		// 	SourcePath:                    SourcePath,
		// 	MoniTreePath:                  "",
		// 	DestinationPath:               DestinationPath,
		// 	ExtensionsToCopyAutomatically: ExtensionsToCopyAutomatically,
		// 	MaximumSizeForCopy:            MaximumSizeForCopy,
		// }
		// error := mediaFake.Fake()
		// if error != nil {
		// 	log.Fatal(error)
		// 	os.Exit(1)
		// }

		// log.Info("Media directory has been faked successfully:", SourcePath)
	},
}
