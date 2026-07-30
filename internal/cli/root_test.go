package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.ExecuteContext(context.Background())
	return buf.String(), err
}

func TestRootCommandHelp(t *testing.T) {
	root := NewRootCommand()
	output, err := executeCommand(root, "--help")
	if err != nil {
		t.Fatalf("help: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("expected help output")
	}
}

func TestStartCommand(t *testing.T) {
	cmd := newStartCommand()
	_, err := executeCommand(cmd)
	if err != nil {
		t.Fatalf("start: %v", err)
	}
}

func TestStopCommand(t *testing.T) {
	cmd := newStopCommand()
	_, err := executeCommand(cmd)
	if err != nil {
		t.Fatalf("stop: %v", err)
	}
}

func TestStatusCommand(t *testing.T) {
	cmd := newStatusCommand()
	_, err := executeCommand(cmd)
	if err != nil {
		t.Fatalf("status: %v", err)
	}
}
