package commands

import (
	"fmt"
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
		fmt.Printf("%s -c %s\n", sh, args[0])
	} else if out, err := exec.Command(sh, "-c", args[0]).Output(); err != nil {
		return err
	} else {
		fmt.Print(string(out))
	}

	return nil
}
