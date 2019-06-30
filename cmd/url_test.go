package cmd_test

import (
	"github.com/adrian-gheorghe/mediafaker/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlBadParameters(t *testing.T) {
	defer func() { log.StandardLogger().ExitFunc = nil }()
	var fatal bool
	log.StandardLogger().ExitFunc = func(int) { fatal = true }

	fatal = false
	cmd.RootCmd.SetArgs([]string{"url"})
	cmd.RootCmd.Execute()
	assert.Equal(t, true, fatal)

	fatal = false
	cmd.RootCmd.SetArgs([]string{"url", "--source", "--destination"})
	cmd.RootCmd.Execute()
	assert.Equal(t, true, fatal)
}
