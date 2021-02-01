package install

func Install() error {
	// load or infer config
	// run install cmds (hardcode for now)
	// LOCKFILES are the thing to be using - any changes to it means we need changes
	// save identifiers to compare
	// also needs to know dirs to ignore

	// if only tracking lockfiles, do you want an "install" command?
	// pretty much need it, but maybe the thing you can't do then is "track" changes of
	// non-lockfile stuff (which also probably means we don't have good tooling anyway to know much about how to install in the first place - requirements.txt how do you set up .venv)
	return nil
}
