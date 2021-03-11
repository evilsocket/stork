package commands

import (
	"fmt"
	"os/exec"
)

func init() {
	Available["docker:create_tag"] = &Command{
		Identifier: "docker:create_tag",
		Argc:       2,
		Logic:      dockerCreateTag,
	}
}

func dockerCreateTag(env *Environment, args ...string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	image := args[0]
	version := args[1]
	tagName := fmt.Sprintf("%s:v%s", image, version)
	latest := fmt.Sprintf("%s:latest", image)

	fmt.Printf("[docker] building %s ...\n", tagName)

	if err = do(env.Dry, docker, "build", "-t", tagName, "."); err != nil {
		return err
	} else if err = do(env.Dry, docker, "tag", tagName, latest); err != nil {
		return err
	}

	fmt.Printf("[doker] pushing %s ...\n", tagName)

	if err = do(env.Dry, docker, "push", tagName); err != nil {
		return err
	} else if err = do(env.Dry, docker, "push", latest); err != nil {
		return err
	}

	return nil
}
