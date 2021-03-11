package commands

var Available = make(map[string]*Command)

type Variables map[string]string

type Environment struct {
	Vars Variables
	Dry  bool
}

type Command struct {
	Identifier string
	Argc       int
	Logic      func(env *Environment, args ...string) error
}
