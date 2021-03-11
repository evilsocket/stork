package commands

import (
	"fmt"
	"github.com/evilsocket/islazy/tui"
	"os"
	"os/exec"
	"strings"
)

func do(dry bool, app string, args ...string) error {
	if dry {
		fmt.Printf("%s %s %s\n", tui.Dim("<dry>"), app, strings.Join(args, " "))
		return nil
	}

	cmd := exec.Command(app, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
