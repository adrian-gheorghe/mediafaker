package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(treeFromSSH)
}

// treeFromSSH represents the tree tree_ssh command
var treeFromSSH = &cobra.Command{
	Use:   "tree_ssh",
	Short: "runs mediafaker on a remote tree file accessible via ssh",
	Long:  `This subcommand connects to a remote host via ssh, where it downloads the moni utility. Using moni it generates a json dump of the tree and file structure in that directory. After downloading the file locally it removes both the file and moni from the ssh host and runs mediafaker on the tree json file.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("tree_ssh called")
	},
}
