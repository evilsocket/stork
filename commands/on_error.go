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

	Available["on_error:suppress"] = &Command{
		Identifier: "on_error:suppress",
		Argc:       0,
		Logic:      setOnErrorSuppress,
	}

	Available["on_error:log"] = &Command{
		Identifier: "on_error:log",
		Argc:       1,
		Logic:      setOnErrorLog,
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

func setOnErrorSuppress(env *Environment, args ...string) error {
	env.OnError = SuppressErrors
	return nil
}

func setOnErrorLog(env *Environment, args ...string) error {
	env.OnError = LogErrors
	env.ErrorLog = args[0]
	return nil
}
