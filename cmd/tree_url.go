package cmd

import (
	"github.com/spf13/cobra"
)

// SourceURL represents the local directory path that should be faked
var SourceURL string

// Gzip represents if the tree file is gzipped or not
var Gzip bool

func init() {
	RootCmd.AddCommand(url)
	url.LocalFlags().StringVarP(&SourceURL, "source", "s", "", "Remote URL path where the tree json is stored")
	url.LocalFlags().BoolVarP(&Gzip, "gzip", "g", false, "Load gzipped json tree")
}

// url represents the tree from url command
var url = &cobra.Command{
	Use:   "url",
	Short: "runs mediafaker from a tree json stored remotely",
	Long:  `This subcommand runs mediafaker on a tree formatted by moni that is accessible via http`,
	Run: func(cmd *cobra.Command, args []string) {
		// if SourceURL == "" {
		// 	log.Fatal("Source tree json has not been provided")
		// 	os.Exit(1)
		// }

		// if DestinationPath == "" {
		// 	log.Fatal("Destination directory has not been provided")
		// 	os.Exit(1)
		// }

		// httpClient := http.Client{
		// 	Timeout: time.Second * 10, // Maximum of 2 secs
		// }

		// req, err := http.NewRequest(http.MethodGet, SourceURL, nil)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// res, err := httpClient.Do(req)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// treeBody, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// treeFile, err := ioutil.TempFile(os.TempDir(), "output-*.json")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// if _, err = treeFile.Write(treeBody); err != nil {
		// 	log.Fatal("Failed to write to temporary file", err)
		// }

		// // Close the file
		// if err := treeFile.Close(); err != nil {
		// 	log.Fatal(err)
		// }

		// mediaFake := fakers.MediaFake{
		// 	SourcePath:                    "",
		// 	MoniTreePath:                  path.Join(os.TempDir(), treeFile.Name()),
		// 	DestinationPath:               DestinationPath,
		// 	ExtensionsToCopyAutomatically: ExtensionsToCopyAutomatically,
		// 	MaximumSizeForCopy:            MaximumSizeForCopy,
		// }
		// error := mediaFake.Fake()
		// if error != nil {
		// 	log.Fatal(error)
		// 	os.Exit(1)
		// }
		// defer os.Remove(treeFile.Name())
		// log.Info("Media directory has been faked successfully:", SourcePath)
	},
}
