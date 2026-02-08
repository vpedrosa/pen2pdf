package cmd

import (
	"testing"
)

func TestInfoCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "info [input.pen]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("info command not registered")
	}
}

func TestInfoCommandRequiresArg(t *testing.T) {
	err := infoCmd.Args(infoCmd, []string{})
	if err == nil {
		t.Error("expected error for missing argument")
	}
}
