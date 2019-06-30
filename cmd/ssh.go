package cmd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	sshClient "github.com/adrian-gheorghe/go-sshclient"
	"github.com/adrian-gheorghe/mediafaker/fakers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	spin "github.com/tj/go-spin"
)

// Source represents remote absolute path that should be faked
var Source string

// SSHHost represents the ssh port mediafaker should connect to
var SSHHost string

// SSHUser represents the ssh user mediafaker should use to connect to the remote
var SSHUser string

// SSHKey represents the ssh key mediafaker should use to connect to the remote
var SSHKey string

// SSHPort represents the ssh port mediafaker should connect to
var SSHPort string

func init() {
	RootCmd.AddCommand(Ssh)
	Ssh.Flags().StringVarP(&Source, "source", "s", "", "Remote SSH absolute path that should be faked")
	Ssh.Flags().StringVar(&SSHHost, "ssh-host", "", "SSH Host to be used for the connection")
	Ssh.Flags().StringVar(&SSHUser, "ssh-user", "", "SSH User to be used for the connection")
	Ssh.Flags().StringVar(&SSHKey, "ssh-key", "", "SSH Key to be used for the connection")
	Ssh.Flags().StringVar(&SSHPort, "ssh-port", "22", "SSH Port to be used for the connection")
}

// treeFromSSH represents the tree tree_ssh command
var Ssh = &cobra.Command{
	Use:   "ssh",
	Short: "runs mediafaker on a remote tree file accessible via ssh",
	Long:  `This subcommand connects to a remote host via ssh, where it downloads the moni utility. Using moni it generates a json dump of the tree and file structure in that directory. After downloading the file locally it removes both the file and moni from the ssh host and runs mediafaker on the tree json file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Source == "" {
			log.Fatal("Remote source directory has not been provided")
			return
		}

		if SSHHost == "" {
			log.Fatal("SSH Host has not been provided")
			return
		}

		if SSHUser == "" {
			log.Fatal("SSH User has not been provided")
			return
		}

		if SSHPort == "" {
			log.Fatal("SSH Port has not been provided")
			return
		}

		if SSHKey == "" {
			log.Fatal("SSH Key has not been provided")
			return
		}

		log.Info("Attempting connection to: ", SSHHost+":"+string(SSHPort))

		client, err := sshClient.DialWithKey(SSHHost+":"+SSHPort, SSHUser, SSHKey)
		if err != nil {
			log.Fatal("There was an error connecting to the ssh host: ", err)
			return
		}
		log.Info("Successfully connected to ssh remote host!")

		log.Info("Attempting to download moni ...")
		err = client.Cmd("sh -c \"$(curl -fsSL https://raw.githubusercontent.com/adrian-gheorghe/moni/master/download.sh)\"").Run()
		if err != nil {
			log.Fatal("There was an error downloading moni: ", err)
			return
		}
		log.Info("Successfully downloaded moni!")

		log.Info("Attempting to generate tree json ...")
		err = client.Cmd("./moni --periodic=false --show_tree_diff=false --gzip=true --tree_store=moni_output.json.gz --path=\"" + Source + "\" --algorithm_name=\"MediafakerTreeWalk\" --content_store_max_size=\"" + strconv.Itoa(int(MaximumSizeForCopy)) + "\"").Run()
		if err != nil {
			log.Fatal("There was an error generating tree using moni: ", err)
			return
		}
		log.Info("Successfully generated json tree!")

		log.Info("Attempting to download the tree file ...")
		output, err := client.Cmd("cat moni_output.json.gz").Output()
		if err != nil {
			log.Fatal("There was an error while downloading the tree output: ", err)
			return
		}

		b := bytes.NewBuffer(output)

		var r io.Reader
		r, err = gzip.NewReader(b)
		if err != nil {
			return
		}

		var resB bytes.Buffer
		_, err = resB.ReadFrom(r)
		if err != nil {
			return
		}

		resultData := resB.Bytes()

		treeFile, err := ioutil.TempFile(os.TempDir(), "output-*.json")
		if err != nil {
			log.Fatal("Temporary local file could not be created ...", err)
			return
		}

		log.Info("Attempting to write to temp file ...", treeFile.Name())
		if _, err = treeFile.Write(resultData); err != nil {
			log.Fatal("Failed to write to temporary file", err)
			return
		}

		// Close the file
		if err := treeFile.Close(); err != nil {
			log.Fatal("Temporary file could not be closed", err)
			return
		}

		log.Info("Attempting to remove the remote tree file ...")
		err = client.Cmd("rm moni_output.json.gz").Run()
		if err != nil {
			log.Error("There was an error while removing the tree output: ", err)
		}

		log.Info("Attempting to remove moni from remote ...")
		err = client.Cmd("rm moni").Run()
		if err != nil {
			log.Error("There was an error while removing moni: ", err)
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
		defer client.Close()
		log.Info("Media directory has been faked", SourcePath)
		mediaFake.CalculateTotalDestinationSize()
		mediaFake.PrintInfo()
	},
}

func showSpinner(s *spin.Spinner, name, frames string) {
	s.Set(frames)
	fmt.Printf("\n\n  %s: %s\n\n", name, frames)
	for i := 0; i < 30; i++ {
		fmt.Printf("\r  \033[36m\033[m %s ", s.Next())
		time.Sleep(100 * time.Millisecond)
	}
}
