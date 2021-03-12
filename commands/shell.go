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
	if  sh == "" {
		return fmt.Errorf("$SHELL not defined")
	}

	if env.Dry {
		fmt.Printf("%s %s -c %s\n", tui.Dim("<dry>"), sh, args[0])
	} else if out, err := exec.Command(sh, "-c", args[0]).Output(); err != nil {
		return err
	} else {
		fmt.Print(string(out))
	}

	return nil
}
