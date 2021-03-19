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
	ctx         = make(map[string]string)
	fileName    = defFileName
	code        = ""
	workDir     = ""
	dryRun      = false
	showVersion = false
	err         = (error)(nil)
	theParser   = (*parser.Parser)(nil)
	theCode     = (*parser.AST)(nil)
)

func die(m string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", tui.Red("error"))
	fmt.Printf(prefix+m, args...)
	os.Exit(1)
}

func perror(m string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] ", tui.Red("error"))
	fmt.Printf(prefix+m, args...)
}

func onError(env *commands.Environment, format string, args ...interface{}) {
	switch env.OnError {
	case commands.AbortOnError:
		die(format, args...)
	case commands.ContinueOnError:
		perror(format, args...)
//  case commands.SuppressErrors:
//  	do nothing :)
	case commands.LogErrors:
		perror(format, args...)

		f, err := os.OpenFile(env.ErrorLog, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			perror("main", "can't open %s for writing: %v", env.ErrorLog, err)
			return
		}
		defer f.Close()

		if _, err = f.WriteString(fmt.Sprintf(format, args...) + "\n"); err != nil {
			perror("main", "can't write to %s: %v", env.ErrorLog, err)
		}
	}

}

func init() {
	flag.StringVar(&fileName, "f", fileName, "input .stork file")
	flag.StringVar(&code, "c", code, "evaluate as code")
	flag.BoolVar(&dryRun, "dry-run", dryRun, "will print commands instead of executing them")
	flag.BoolVar(&showVersion, "v", showVersion, "print version and exit")
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("stork v%s\n", Version)
		return
	}

	if theParser, err = parser.New(); err != nil {
		die("%v\n", err)
	}

	if code == "" {
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

		if theCode, err = theParser.ParseFile(fileName); err != nil {
			die("parsing %s: %v\n", fileName, err)
		}
	} else {
		if theCode, err = theParser.ParseCode(code); err != nil {
			die("parsing %s: %v\n", code, err)
		}
	}

	// pre check that all the command are defined
	for _, step := range theCode.Steps {
		if step.Command != nil {
			if _, found := commands.Available[step.Command.Identifier]; !found {
				die("%v %s undefined command\n",
					step.Command.Pos,
					step.Command.Identifier)
			}
		}
	}

	env := commands.Environment{
		Vars:    make(commands.Variables),
		Dry:     dryRun,
		OnError: commands.AbortOnError,
	}
	for _, step := range theCode.Steps {
		if step.Set != nil {
			if value, err := step.Set.Value.Resolve(env.Vars); err != nil {
				onError(&env, "%v\n", err)
			} else {
				env.Vars[step.Set.Identifier[1:]] = value
			}
		} else {
			cmd, _ := commands.Available[step.Command.Identifier]

			if argc := len(step.Command.Parameters); argc != cmd.Argc {
				onError(&env, "%v %s requires %d parameters, %d provided\n",
					step.Command.Pos,
					step.Command.Identifier,
					cmd.Argc,
					argc)
			}

			// fmt.Printf("vars: %#v\n", env.Vars)

			var argv []string
			for _, param := range step.Command.Parameters {
				if value, err := param.Resolve(env.Vars); err != nil {
					onError(&env, "%v\n", err)
				} else {
					argv = append(argv, value)
				}
			}

			// fmt.Printf("argv: %#v\n", argv)

			if err = cmd.Logic(&env, argv...); err != nil {
				onError(&env, "%s: %v\n", cmd.Identifier, err)
			}
		}
	}
}
