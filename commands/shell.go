package commands

import (
	"fmt"
	"github.com/evilsocket/islazy/tui"
	"os"
	"os/exec"
)

var sh = ""

func init() {
	sh = os.Getenv("SHELL")

	Available["shell:do"] = &Command{
		Identifier: "shell:do",
		Argc:       1,
		Logic:      shellDo,
	}
}

func shellDo(env *Environment, args ...string) error {
	if sh == "" {
		return fmt.Errorf("$SHELL not defined")
	}

	if env.Dry {
		fmt.Printf("%s %s -c %s\n", tui.Dim("<dry>"), sh, args[0])
	} else {
		msg("shell", "%s\n", args[0])
		cmd := exec.Command(sh, "-c", args[0])

		for _, vEnv := range env.Vars.AsEnv() {
			cmd.Env = append(cmd.Env, vEnv)
		}

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
