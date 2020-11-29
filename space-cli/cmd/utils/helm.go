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
)

// HelmInstall install helm chart
func HelmInstall(chartReleaseName, chartLocation, downloadURL, namespace string, valuesFileObj map[string]interface{}) (*chart.Chart, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}

	var helmChart *chart.Chart
	var err error
	if chartLocation == "" {
		res, err := http.Get(downloadURL)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("received invalid status code (%s)", res.Status)
		}
		helmChart, err = loader.LoadArchive(res.Body)
		if err != nil {
			return nil, err
		}
	} else {
		helmChart, err = loader.Load(chartLocation)
		if err != nil {
			return nil, err
		}
	}

	iCli := action.NewInstall(actionConfig)
	iCli.ReleaseName = chartReleaseName
	iCli.CreateNamespace = namespace != ""
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
	LogInfo(fmt.Sprintf("Successfully removed space cloud: (%s)", rel.Release.Name))
	return nil
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
