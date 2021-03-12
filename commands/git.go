package commands

import (
	"fmt"
	"os/exec"
)

func init() {
	Available["git:create_tag"] = &Command{
		Identifier: "git:create_tag",
		Argc:       1,
		Logic:      gitCreateTag,
	}
}

func gitCreateTag(env *Environment, args ...string) error {
	versionFile := env.Vars["VERSION_FILE"]
	version := env.Vars["VERSION"]

	if versionFile == "" {
		return fmt.Errorf("VERSION_FILE not set")
	} else if version == "" {
		return fmt.Errorf("VERSION not set")
	}

	git, err := exec.LookPath("git")
	if err != nil {
		return err
	}

	// TODO: create changelog

	// add and push version file in case it was changed
	txt := fmt.Sprintf("releasing v%s", version)

	msg("git", "%s ...\n", txt)

	if err = do(env.Dry, git, "add", versionFile); err != nil {
		return err
	} else if err = do(env.Dry, git, "commit", "-m", txt); err != nil {
		return err
	} else if err = do(env.Dry, git, "push"); err != nil {
		return err
	}

	// create new tag and push
	tag := fmt.Sprintf("v%s", version)
	if err = do(env.Dry, git, "tag", "-a", tag, "-m", txt); err != nil {
		return err
	} else if err = do(env.Dry, git, "push", "origin", tag); err != nil {
		return err
	}

	return nil
}
