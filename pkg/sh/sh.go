package sh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(verbose bool, name string, arg ...string) error {
	if verbose {
		cmdText := name + " " + strings.Join(arg, " ")
		fmt.Fprintln(os.Stderr, " + ", cmdText)
	}
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	return cmd.Run()
}
