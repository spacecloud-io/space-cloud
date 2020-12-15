package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	"github.com/spaceuptech/space-cloud/space-cli/cmd/model"
)

// HelmInstall install helm chart
func HelmInstall(chartReleaseName, chartLocation, downloadURL, namespace string, valuesFileObj map[string]interface{}) (*chart.Chart, error) {
	actionConfig, helmChart, err := CreateChart(chartLocation, downloadURL)
	if err != nil {
		return nil, err
	}

	iCli := action.NewInstall(actionConfig)
	iCli.ReleaseName = chartReleaseName
	iCli.CreateNamespace = namespace != model.HelmSpaceCloudNamespace // don't create namespace while setting up space cloud cluster
	iCli.Namespace = namespace
	iCli.Devel = true
	rel, err := iCli.Run(helmChart, valuesFileObj)
	if err != nil {
		return nil, err
	}

	LogInfo(fmt.Sprintf("Successfully installed chart: (%s)", rel.Name))
	return rel.Chart, nil
}

// HelmUninstall uninstall helm chart
func HelmUninstall(releaseName string) error {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return err
	}
	uCli := action.NewUninstall(actionConfig)
	rel, err := uCli.Run(releaseName)
	if err != nil {
		return err
	}

	LogInfo(fmt.Sprintf("Successfully removed chart: (%s)", rel.Release.Name))
	return nil
}

// HelmShow executes the helm show command
func HelmShow(chartLocation, downloadURL, releaseArg string) (string, error) {
	_, _, err := CreateChart(chartLocation, downloadURL)
	if err != nil {
		return "", err
	}

	sCli := action.NewShow(action.ShowValues)
	info, err := sCli.Run("")
	if err != nil {
		return "", err
	}
	return info, nil
}

// HelmGet gets chart info
func HelmGet(releaseName string) (*release.Release, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	iCli := action.NewGet(actionConfig)
	info, err := iCli.Run(releaseName)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// HelmList uninstall helm chart
func HelmList(filterRegex string) ([]*release.Release, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}
	lCli := action.NewList(actionConfig)
	list, err := lCli.Run()
	if err != nil {
		return nil, err
	}

	filteredList := make([]*release.Release, 0)
	for _, item := range list {
		if strings.Contains(item.Namespace, filterRegex) {
			filteredList = append(filteredList, item)
		}
	}

	return filteredList, nil
}

// HelmUpgrade upgrade space cloud chart
func HelmUpgrade(releaseName, chartLocation, downloadURL, namespace string, valuesFileObj map[string]interface{}) (*chart.Chart, error) {
	actionConfig, helmChart, err := CreateChart(chartLocation, downloadURL)
	if err != nil {
		return nil, err
	}

	uCli := action.NewUpgrade(actionConfig)
	uCli.Namespace = namespace
	rel, err := uCli.Run(releaseName, helmChart, valuesFileObj)
	if err != nil {
		return nil, err
	}

	LogInfo(fmt.Sprintf("Successfully upgraded chart: (%s)", rel.Name))
	return rel.Chart, nil
}

// CreateChart returns chart object, which describe the provide chart
func CreateChart(chartLocation, downloadURL string) (*action.Configuration, *chart.Chart, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, nil, err
	}

	var helmChart *chart.Chart
	var err error
	if chartLocation == "" {
		res, err := http.Get(downloadURL)
		if err != nil {
			return nil, nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("received invalid status code (%s)", res.Status)
		}
		helmChart, err = loader.LoadArchive(res.Body)
		if err != nil {
			return nil, nil, err
		}
	} else {
		helmChart, err = loader.Load(chartLocation)
		if err != nil {
			return nil, nil, err
		}
	}
	return actionConfig, helmChart, err
}

// ExtractValuesObj extract chart values from yaml file & cli flags
func ExtractValuesObj(setValuesFlag, valuesYamlFile string) (map[string]interface{}, error) {
	valuesFileObj := map[string]interface{}{}
	if valuesYamlFile != "" {
		var bodyInBytes []byte
		var err error
		if strings.HasPrefix(valuesYamlFile, "http") {
			// download file from the internet
			resp, err := http.Get(valuesYamlFile)
			if err != nil {
				return nil, err
			}
			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("")
			}
			bodyInBytes, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			// read locally available file
			bodyInBytes, err = ioutil.ReadFile(valuesYamlFile)
			if err != nil {
				return nil, err
			}
		}

		if err := yaml.Unmarshal(bodyInBytes, &valuesFileObj); err != nil {
			return nil, err
		}
	}

	setValuesObj := map[string]interface{}{}
	if setValuesFlag != "" {
		arr := strings.Split(setValuesFlag, ",")
		for _, element := range arr {
			tempArr := strings.Split(element, "=")
			if len(tempArr) != 2 {
				return nil, fmt.Errorf("invalid value (%s) provided for flag --set, it should be in format foo1=bar1,foo2=bar2", tempArr)
			}
			setValuesObj[tempArr[0]] = tempArr[1]
		}
	}

	// override values of yaml file
	for key, value := range setValuesObj {
		valuesFileObj[key] = value
	}

	return valuesFileObj, nil
}
