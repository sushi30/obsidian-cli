package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandAliasesAreUnique(t *testing.T) {
	seen := map[string]string{}
	for _, cmd := range rootCmd.Commands() {
		for _, alias := range cmd.Aliases {
			if existing, ok := seen[alias]; ok {
				t.Errorf("alias %q is claimed by both %q and %q", alias, existing, cmd.Name())
			}
			seen[alias] = cmd.Name()
		}
	}
}

func TestRemoveAlias(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"rm"})
	assert.NoError(t, err)
	assert.Equal(t, "remove", cmd.Name())
}

func TestMoveAlias(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"mv"})
	assert.NoError(t, err)
	assert.Equal(t, "move", cmd.Name())
}
