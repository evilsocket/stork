package main

import (
	"flag"
	"fmt"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/tui"
	"github.com/evilsocket/stork/commands"
	"github.com/evilsocket/stork/parser"
	"os"
	"path/filepath"
)

const defFileName = "release.stork"

var (
	ctx       = make(map[string]string)
	fileName  = defFileName
	workDir   = ""
	dryRun    = false
	err       = (error)(nil)
	theParser = (*parser.Parser)(nil)
	theCode   = (*parser.AST)(nil)
)

func die(m string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", tui.Red("error"))
	fmt.Printf(prefix+m, args...)
	os.Exit(1)
}

func init() {
	flag.StringVar(&fileName, "f", fileName, "input .stork file")
	flag.BoolVar(&dryRun, "dry-run", dryRun, "will print commands instead of executing them")
}

func main() {
	flag.Parse()

	if fileName, err = filepath.Abs(fileName); err != nil {
		die("%v\n", err)
	} else if fs.Exists(fileName) == false {
		die("%s does not exist\n", fileName)
	}
	oldCwd, err := os.Getwd()
	if err != nil {
		die("%v\n", err)
	}

	workDir = filepath.Dir(fileName)
	if err = os.Chdir(workDir); err != nil {
		die("%v\n", err)
	}
	defer os.Chdir(oldCwd)

	// fmt.Printf("using %s\n", fileName)

	if theParser, err = parser.New(); err != nil {
		die("%v\n", err)
	}

	if theCode, err = theParser.ParseFile(fileName); err != nil {
		die("parsing %s: %v\n", fileName, err)
	}

	// pre check that all the command are defined
	for _, step := range theCode.Steps {
		if step.Command != nil {
			if _, found := commands.Available[step.Command.Identifier]; !found {
				die("%s:%v %s undefined command\n",
					fileName,
					step.Command.Pos,
					step.Command.Identifier)
			}
		}
	}

	env := commands.Environment{
		Vars: make(commands.Variables),
		Dry:  dryRun,
	}
	for _, step := range theCode.Steps {
		if step.Set != nil {
			if value, err := step.Set.Value.Resolve(env.Vars); err != nil {
				die("%v\n", err)
			} else {
				env.Vars[step.Set.Identifier[1:]] = value
			}
		} else {
			cmd, _ := commands.Available[step.Command.Identifier]

			if argc := len(step.Command.Parameters); argc != cmd.Argc {
				die("%s:%v %s requires %d parameters, %d provided\n",
					fileName,
					step.Command.Pos,
					step.Command.Identifier,
					cmd.Argc,
					argc)
			}

			// fmt.Printf("vars: %#v\n", env.Vars)

			var argv []string
			for _, param := range step.Command.Parameters {
				if value, err := param.Resolve(env.Vars); err != nil {
					die("%v\n", err)
				} else {
					argv = append(argv, value)
				}
			}

			// fmt.Printf("argv: %#v\n", argv)

			if err = cmd.Logic(&env, argv...); err != nil {
				die("%s: %v\n", cmd.Identifier, err)
			}
		}
	}
}
