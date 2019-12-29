package taskner

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Job struct {
	Command   Command
	StdinFile string
}

func (j *Job) Run(ctx context.Context) ([]byte, error) {
	cmd := j.Command.build(ctx)

	if j.StdinFile != "" {
		stdin, _ := cmd.StdinPipe()
		f, _ := os.Open(j.StdinFile)
		defer f.Close()
		io.Copy(stdin, f)

		stdin.Close()
	}

	return cmd.CombinedOutput()
}

type Command struct {
	command string
	args    []string
}

func getCommand(s string) Command {
	cs := strings.Split(s, " ")
	return Command{
		command: cs[0],
		args:    cs[1:],
	}
}

func (c *Command) Run() ([]byte, error) {
	return c.build(context.Background()).CombinedOutput()
}

func (c *Command) build(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, c.command, c.args...)
	return cmd
}
