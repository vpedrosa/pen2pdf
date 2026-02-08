package cmd

import (
	"testing"
)

func TestValidateCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "validate [input.pen]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("validate command not registered")
	}
}

func TestValidateCommandRequiresArg(t *testing.T) {
	err := validateCmd.Args(validateCmd, []string{})
	if err == nil {
		t.Error("expected error for missing argument")
	}
}
