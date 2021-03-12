package commands

import (
	"fmt"
	"github.com/evilsocket/islazy/str"
	"os/exec"
	"strings"
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

	// determine current tag
	curTag := ""
	cmd := exec.Command(git, "describe", "--tags", "--abbrev=0")
	if output, err := cmd.CombinedOutput(); err == nil {
		curTag = str.Trim(string(output))
		msg("git", "current tag is %s\n", curTag)
	} else {
		msg("git", "no current tag detected\n")
	}

	// get commits from the current tag (if present) to HEAD
	tagRange := "HEAD"
	if curTag != "" {
		tagRange = fmt.Sprintf("%s..HEAD", curTag)
	}

	var fixes []string
	var features []string
	var misc []string

	cmd = exec.Command(git, "log", tagRange, "--oneline")
	if output, err := cmd.CombinedOutput(); err == nil {
		lines := strings.Split(str.Trim(string(output)), "\n")
		for _, line := range lines {
			lwrLine := strings.ToLower(line)
			if strings.Contains(lwrLine, "merge pull request") {
				continue
			} else if strings.Contains(lwrLine, "fix") {
				fixes = append(fixes, line)
			} else if strings.Contains(lwrLine, "new") || strings.Contains(line, "add") {
				features = append(features, line)
			} else {
				misc = append(misc, line)
			}
		}
	} else {
		return err
	}

	msg("git", "release changelog:\n\n")

	fmt.Print("Changelog\n===\n\n")

	if len(features) > 0 {
		fmt.Printf("**New Features**\n\n")
		for _, commit := range features {
			fmt.Printf("* %s\n", commit)
		}
		fmt.Println()
	}

	if len(fixes) > 0 {
		fmt.Printf("**Fixes**\n\n")
		for _, commit := range fixes {
			fmt.Printf("* %s\n", commit)
		}
		fmt.Println()
	}

	if len(misc) > 0 {
		fmt.Printf("**Misc**\n\n")
		for _, commit := range misc {
			fmt.Printf("* %s\n", commit)
		}
		fmt.Println()
	}

	// add and push version file in case it was changed
	txt := fmt.Sprintf("releasing v%s", version)

	msg("git", "%s ...\n", txt)

	if err = do(env.Dry, git, "add", versionFile); err != nil {
		return err
	} else if err = do(env.Dry, git, "commit", "-m", txt); err != nil {
		// if version file didn't change, this will exit with code 1 ... just ignore
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
