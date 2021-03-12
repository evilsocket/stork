package commands

import (
	"fmt"
	"os/exec"
)

func init() {
	Available["docker:repository"] = &Command{
		Identifier: "docker:repository",
		Argc:       1,
		Logic:      dockerSetRepo,
	}

	Available["docker:build"] = &Command{
		Identifier: "docker:build",
		Argc:       3,
		Logic:      dockerBuild,
	}

	Available["docker:push"] = &Command{
		Identifier: "docker:push",
		Argc:       1,
		Logic:      dockerPush,
	}

	Available["docker:create_tag"] = &Command{
		Identifier: "docker:create_tag",
		Argc:       3,
		Logic:      dockerCreateTag,
	}
}

func dockerSetRepo(env *Environment, args ...string) error {
	env.Vars["DOCKER_REPOSITORY"] = args[0]
	return nil
}

func dockerBuild(env *Environment, args ...string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	image := args[0]
	dockerfile := args[1]
	tag := args[2]
	tagName := fmt.Sprintf("%s:%s", image, tag)

	msg("docker", "building %s ...\n", tagName)

	if err = do(env.Dry, docker, "build", "-t", tagName, dockerfile); err != nil {
		return err
	}

	return nil
}

func dockerPush(env *Environment, args ...string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	targetName := args[0]
	if repo := env.Vars["DOCKER_REPOSITORY"]; repo != "" {
		targetName = fmt.Sprintf("%s/%s", repo, targetName)
	}

	msg("docker", "pushing %s ...\n", targetName)
	if err = do(env.Dry, docker, "push", targetName); err != nil {
		return err
	}

	return nil
}

func dockerCreateTag(env *Environment, args ...string) error {
	docker, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	image := args[0]
	srcTag := args[1]
	version := args[2]

	sourceName := fmt.Sprintf("%s:%s", image, srcTag)
	remoteSourceName := sourceName
	targetName := fmt.Sprintf("%s:%s", image, version)

	if repo := env.Vars["DOCKER_REPOSITORY"]; repo != "" {
		targetName = fmt.Sprintf("%s/%s", repo, targetName)
		remoteSourceName = fmt.Sprintf("%s/%s", repo, sourceName)
	}

	msg("docker", "tagging %s from %s ...\n", targetName, sourceName)
	if err = do(env.Dry, docker, "tag", sourceName, targetName); err != nil {
		return err
	}

	if sourceName != remoteSourceName {
		msg("docker", "tagging %s from %s ...\n", remoteSourceName, sourceName)
		if err = do(env.Dry, docker, "tag", sourceName, remoteSourceName); err != nil {
			return err
		}
	}

	msg("docker", "pushing %s and %s ...\n", targetName, remoteSourceName)

	if err = do(env.Dry, docker, "push", targetName); err != nil {
		return err
	} else if err = do(env.Dry, docker, "push", remoteSourceName); err != nil {
		return err
	}

	return nil
}
