package deploy

import (
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cli/cmd"
	"github.com/spaceuptech/space-cli/utils"
)

func deployService(dockerFilePath, serviceFilePath string) error {
	// Check if docker file exists
	if !utils.FileExists(dockerFilePath) {
		return utils.LogError(fmt.Sprintf("Docker file (%s) not found. Try running `space-cli deploy --prepare`", dockerFilePath), "deploy", "deploy", nil)
	}

	// Check if service config file exists
	if !utils.FileExists(serviceFilePath) {
		return utils.LogError(fmt.Sprintf("Service config file (%s) not found. Try running `space-cli deploy --prepare`", dockerFilePath), "deploy", "deploy", nil)
	}

	// Get the spec object from the file
	specObj, err := utils.ReadSpecObjectsFromFile(serviceFilePath)
	if err != nil {
		return utils.LogError("Unable to read spec object from file", "deploy", "deploy", err)
	}
	if len(specObj) != 1 {
		return utils.LogError("There can only be a single object in the service config file", "deploy", "deploy", fmt.Errorf("found %d spec objects", len(specObj)))
	}

	// Select the task with the id same as service id
	var dockerImage string
	serviceID := specObj[0].Meta["id"]
	tasks := specObj[0].Spec.(map[string]interface{})["tasks"].([]interface{})
	for _, task := range tasks {
		taskObj := task.(map[string]interface{})

		if taskObj["id"] == serviceID {
			dockerImage = taskObj["docker"].(map[string]interface{})["image"].(string)
			break
		}
	}

	// Check if docker image was found
	if dockerImage == "" {
		return utils.LogError("Unable to detect project to be built. Make sure you have a task id equal to the service id", "deploy", "deploy", nil)
	}

	// Execute the docker build command
	utils.LogInfo(fmt.Sprintf("Building docker image (%s)", dockerImage), "deploy", "deploy")
	if output, err := exec.Command("docker", "build", "--file", dockerFilePath, "--no-cache", "-t", dockerImage, ".").CombinedOutput(); err != nil {
		return utils.LogError(fmt.Sprintf("Unable to build docker image (%s) - %s", dockerImage, string(output)), "deploy", "deploy", err)
	}

	// Execute the docker push command
	utils.LogInfo(fmt.Sprintf("Pushing docker image (%s)", dockerImage), "deploy", "deploy")
	if output, err := exec.Command("docker", "push", dockerImage).CombinedOutput(); err != nil {
		logrus.Errorln(string(output))
		return utils.LogError(fmt.Sprintf("Unable to push docker image (%s). Have you logged into your registry?", dockerImage), "deploy", "deploy", err)
	}

	// Time to apply the service file config
	if err := cmd.Apply(serviceFilePath); err != nil {
		return utils.LogError("Unable to apply service file config", "deploy", "deploy", err)
	}

	return nil
}
