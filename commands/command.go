package commands

import (
	"fmt"
)

var Available = make(map[string]*Command)

type Variables map[string]string

func (v Variables) AsEnv() []string {
	env := make([]string, len(v))
	for k, v := range v {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

type ErrorHandlingType int

const (
	AbortOnError ErrorHandlingType = iota
	ContinueOnError
	SuppressErrors
	LogErrors
)

type Environment struct {
	Vars     Variables
	Dry      bool
	OnError  ErrorHandlingType
	ErrorLog string
}

type Command struct {
	Identifier string
	Argc       int
	Logic      func(env *Environment, args ...string) error
}
