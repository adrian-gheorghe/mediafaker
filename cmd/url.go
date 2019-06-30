package cmd

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/adrian-gheorghe/mediafaker/fakers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SourceURL represents the local directory path that should be faked
var SourceURL string

func init() {
	RootCmd.AddCommand(Url)
	Url.Flags().StringVarP(&SourceURL, "source", "s", "", "Remote URL path where the tree json is stored")
}

// Url represents the tree from url command
var Url = &cobra.Command{
	Use:   "url",
	Short: "runs mediafaker from a tree json stored remotely",
	Long:  `This subcommand runs mediafaker on a tree formatted by moni that is accessible via http`,
	Run: func(cmd *cobra.Command, args []string) {
		if SourceURL == "" {
			log.Fatal("Source tree json has not been provided")
			return
		}

		if DestinationPath == "" {
			log.Fatal("Destination directory has not been provided")
			return
		}

		log.Info("Attempting to download remote json tree ...")
		httpClient := http.Client{
			Timeout: time.Second * 10, // Maximum of 10 secs
		}

		req, err := http.NewRequest(http.MethodGet, SourceURL, nil)
		if err != nil {
			log.Fatal("Remote json tree could not be loaded ...", err)
			return
		}

		res, err := httpClient.Do(req)
		if err != nil {
			log.Fatal("Remote json tree could not be loaded ...", err)
			return
		}

		treeBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		treeFile, err := ioutil.TempFile(os.TempDir(), "output-*.json")
		if err != nil {
			log.Fatal("Temporary file could not be created ...", err)
			return
		}

		if _, err = treeFile.Write(treeBody); err != nil {
			log.Fatal("Failed to write to temporary file", err)
			return
		}

		// Close the file
		if err := treeFile.Close(); err != nil {
			log.Fatal("Temporary file could not be closed", err)
			return
		}

		mediaFake := fakers.MediaFake{
			SourcePath:                    "",
			MoniTreePath:                  treeFile.Name(),
			DestinationPath:               DestinationPath,
			ExtensionsToCopyAutomatically: ExtensionsToCopyAutomatically,
			MaximumSizeForCopy:            MaximumSizeForCopy,
		}
		error := mediaFake.Fake()
		if error != nil {
			log.Error(error)
		}
		defer os.Remove(treeFile.Name())
		log.Info("Media directory has been faked", SourcePath)
		mediaFake.CalculateTotalDestinationSize()
		mediaFake.PrintInfo()
	},
}
