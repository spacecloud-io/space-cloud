package deploy

import (
	"fmt"
	"os/exec"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/modules/operations"
	"github.com/spaceuptech/space-cloud/space-cli/cmd/utils"
)

func deployService(dockerFilePath, serviceFilePath string) error {
	// Check if docker file exists
	if !utils.FileExists(dockerFilePath) {
		return utils.LogError(fmt.Sprintf("Docker file (%s) not found. Try running `space-cli deploy --prepare`", dockerFilePath), nil)
	}

	// Check if service config file exists
	if !utils.FileExists(serviceFilePath) {
		return utils.LogError(fmt.Sprintf("Service config file (%s) not found. Try running `space-cli deploy --prepare`", dockerFilePath), nil)
	}

	// Get the spec object from the file
	specObj, err := utils.ReadSpecObjectsFromFile(serviceFilePath)
	if err != nil {
		return utils.LogError("Unable to read spec object from file", err)
	}
	if len(specObj) != 1 {
		return utils.LogError("There can only be a single object in the service config file", fmt.Errorf("found %d spec objects", len(specObj)))
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
		return utils.LogError("Unable to detect project to be built. Make sure you have a task id equal to the service id", nil)
	}

	// Execute the docker build command
	utils.LogInfo(fmt.Sprintf("Building docker image (%s)", dockerImage))
	if output, err := exec.Command("docker", "build", "--file", dockerFilePath, "--no-cache", "-t", dockerImage, ".").CombinedOutput(); err != nil {
		return utils.LogError(fmt.Sprintf("Unable to build docker image (%s) - %s", dockerImage, string(output)), err)
	}

	// Execute the docker push command
	utils.LogInfo(fmt.Sprintf("Pushing docker image (%s)", dockerImage))
	if output, err := exec.Command("docker", "push", dockerImage).CombinedOutput(); err != nil {
		_ = utils.LogError(string(output), nil)
		return utils.LogError(fmt.Sprintf("Unable to push docker image (%s). Have you logged into your registry?", dockerImage), err)
	}

	// Time to apply the service file config
	if err := operations.Apply(serviceFilePath, true, model.ApplyWithNoDelay, 1); err != nil {
		return utils.LogError("Unable to apply service file config", err)
	}

	return nil
}
