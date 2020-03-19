package deploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/spaceuptech/space-cli/modules/services"
	"github.com/spaceuptech/space-cli/utils"
)

func prepareService(projectID, dockerFilePath, serviceFilePath string) error {

	// Check if a docker file exists
	if !utils.FileExists(dockerFilePath) {
		utils.LogInfo(fmt.Sprintf("Could not find docker file (%s)", dockerFilePath), "deploy", "prepare")
		utils.LogInfo(fmt.Sprintf("Creating docker file (%s)", dockerFilePath), "deploy", "prepare")
		// We need to create a docker file at the path provided. First lets try and detect the programming language
		lang, err := getLanguage()
		if err != nil {
			return utils.LogError("Could not detect programing language. Only python, js and golang are currently supported. For other languages, make sure you have a Dockerfile prepared", "deploy", "prepare", err)
		}

		utils.LogInfo(fmt.Sprintf("Language detected (%s)", lang), "deploy", "prepare")

		var dockerFileContents string
		switch lang {
		case "golang":
			dockerFileContents = utils.DockerfileGolang
		case "javascript":
			dockerFileContents = utils.DockerfileNodejs
		case "python":
			dockerFileContents = utils.DockerfilePython
		default:
			return utils.LogError(fmt.Sprintf("Lnaguage (%s) not supported. Consider making a Dockerfile yourself.", lang), "deploy", "prepare", nil)
		}

		// Create the docker file
		utils.LogInfo("Creating docker file with following contents:", "deploy", "prepare")
		fmt.Println()
		fmt.Println(dockerFileContents)
		fmt.Println()
		if err := utils.CreateFileIfNotExist(dockerFilePath, dockerFileContents); err != nil {
			return utils.LogError(fmt.Sprintf("Could not create docker file (%s)", dockerFilePath), "deploy", "prepare", err)
		}
	}

	// Check if a services file exist
	if !utils.FileExists(serviceFilePath) {
		utils.LogInfo(fmt.Sprintf("Could not find service file (%s)", serviceFilePath), "deploy", "prepare")

		svc, err := services.GenerateService(projectID, "auto")
		if err != nil {
			return utils.LogError("Could not generate service config", "deploy", "prepare", err)
		}

		data, _ := yaml.Marshal(svc)
		utils.LogInfo("Creating service config file with following contents:", "deploy", "prepare")
		fmt.Println()
		fmt.Println(string(data))
		fmt.Println()
		if err := utils.AppendConfigToDisk(svc, serviceFilePath); err != nil {
			return utils.LogError(fmt.Sprintf("Could not create service config file (%s)", dockerFilePath), "deploy", "prepare", err)
		}
	}

	utils.LogInfo("All configuration has been saved successfully. Run `space-cli deploy` to deploy your service!", "deploy", "prepare")
	return nil
}

func getLanguage() (string, error) {
	// Iterate over all files in the working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(workingDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		// Skip all directories
		if file.IsDir() {
			continue
		}

		// Check if folder has a requirements.txt file (python)
		if file.Name() == "requirements.txt" {
			return "python", nil
		}

		// Check if folder has a package.json file (js)
		if file.Name() == "package.json" {
			return "javascript", nil
		}

		// Check if folder has a .go file (golang)
		if strings.HasSuffix(file.Name(), ".go") {
			return "golang", nil
		}
	}

	return "", fmt.Errorf("could not detect programing language")
}
