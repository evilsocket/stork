package commands

import (
	"fmt"
	"github.com/evilsocket/islazy/tui"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	// TODO: make this more versatile
	versionFileParser = regexp.MustCompile(`[Vv]ersion\s*=\s*['"]([\d\.ab]+)["']`)
	versionParser = regexp.MustCompile(`\d+\.\d+\.\d+[ab]?`)
)

func init() {
	Available["version:file"] = &Command{
		Identifier: "version:file",
		Argc:       1,
		Logic:      versionFile,
	}

	Available["version:from_user"] = &Command{
		Identifier: "version:from_user",
		Logic:      versionFromUser,
	}
}

func versionFile(env *Environment, args ...string) error {
	fileName := args[0]

	// fmt.Printf("reading version from %s\n", fileName)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if matches := versionFileParser.FindStringSubmatch(string(data)); matches == nil {
		return fmt.Errorf("can't parse version from %s", fileName)
	} else {
		env.Vars["VERSION_FILE"] = fileName
		env.Vars["VERSION"] = matches[1]
	}

	return nil
}

func versionFromUser(env *Environment, args ...string) error {
	versionFile := env.Vars["VERSION_FILE"]
	version := env.Vars["VERSION"]

	if versionFile == "" {
		return fmt.Errorf("VERSION_FILE not set")
	} else if version == "" {
		return fmt.Errorf("VERSION not set")
	}

	fmt.Printf("[v%s] enter new version (major.minor.patch): ", version)
	var newVersion string
	fmt.Scanln(&newVersion)
	if versionParser.MatchString(newVersion) == false {
		return fmt.Errorf("'%s' is not a valid version, use the major.minor.patch format", newVersion)
	}

	if !env.Dry {
		data, err := ioutil.ReadFile(versionFile)
		if err != nil {
			return err
		}

		newData := strings.ReplaceAll(string(data), version, newVersion)
		// FIXME: save the original permissions somewhere and restore them here
		if err = ioutil.WriteFile(versionFile, []byte(newData), os.ModePerm); err != nil {
			return err
		}
	} else {
		fmt.Printf("%s update %s: %s -> %s\n", tui.Dim("<dry>"), versionFile, version, newVersion)
	}

	env.Vars["VERSION"] = newVersion

	return nil
}
