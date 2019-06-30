package cmd_test

import (
	"io/ioutil"
	"testing"

	"github.com/adrian-gheorghe/mediafaker/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	m.Run()
}

func TestVersionWorks(t *testing.T) {
	cmd.RootCmd.SetArgs([]string{"--version"})
	assert.NotPanics(t, func() {
		assert.NoError(t, cmd.RootCmd.Execute())
	})
}
