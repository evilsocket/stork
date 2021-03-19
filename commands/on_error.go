package commands


func init() {
	Available["on_error:abort"] = &Command{
		Identifier: "on_error:abort",
		Argc:       0,
		Logic:      setOnErrorAbort,
	}

	Available["on_error:continue"] = &Command{
		Identifier: "on_error:continue",
		Argc:       0,
		Logic:      setOnErrorContinue,
	}
}

func setOnErrorAbort(env *Environment, args ...string) error {
	env.OnError = AbortOnError
	return nil
}

func setOnErrorContinue(env *Environment, args ...string) error {
	env.OnError = ContinueOnError
	return nil
}
