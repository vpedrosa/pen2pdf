package cmd

import (
	"testing"
)

func TestRenderCommandRegistered(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "render [input.pen]" {
			found = true
			break
		}
	}
	if !found {
		t.Error("render command not registered")
	}
}

func TestRenderCommandHasOutputFlag(t *testing.T) {
	f := renderCmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("expected --output flag")
	}
	if f.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got '%s'", f.Shorthand)
	}
}

func TestRenderCommandHasPagesFlag(t *testing.T) {
	f := renderCmd.Flags().Lookup("pages")
	if f == nil {
		t.Error("expected --pages flag")
	}
}

func TestRenderCommandRequiresArg(t *testing.T) {
	err := renderCmd.Args(renderCmd, []string{})
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestRenderCommandHasNoPromptFlag(t *testing.T) {
	f := renderCmd.Flags().Lookup("no-prompt")
	if f == nil {
		t.Error("expected --no-prompt flag")
	}
}
