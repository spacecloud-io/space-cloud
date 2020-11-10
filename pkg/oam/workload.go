package oam

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"github.com/oam-dev/kubevela/api/types"
	"github.com/oam-dev/kubevela/pkg/application"
	"github.com/oam-dev/kubevela/pkg/commands/util"
	cmdutil "github.com/oam-dev/kubevela/pkg/commands/util"
	"github.com/oam-dev/kubevela/pkg/plugins"
	"github.com/spf13/pflag"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RunOptions struct {
	Env          *types.EnvMeta
	WorkloadName string
	KubeClient   client.Client
	App          *application.Application
	AppName      string
	Staging      bool
	util.IOStreams
}

func LoadIfExist(envName string, workloadName string, appGroup string) (*application.Application, error) {
	var appName string
	if appGroup != "" {
		appName = appGroup
	} else {
		appName = workloadName
	}
	app, err := application.Load(envName, appName)
	if err != nil {
		return nil, err
	}
	app.Name = appName

	return app, nil
}

func BaseComplete(envName string, workloadName string, appName string, flagSet *pflag.FlagSet, workloadType string) (*application.Application, error) {
	app, err := LoadIfExist(envName, workloadName, appName)
	if err != nil {
		return nil, err
	}
	tp, workloadData := app.GetWorkload(workloadName)
	if tp == "" {
		if workloadType == "" {
			return nil, fmt.Errorf("must specify workload type for application %s", workloadName)
		}
		// Not exist
		tp = workloadType
	}
	template, err := plugins.LoadCapabilityByName(tp)
	if err != nil {
		return nil, err
	}

	for _, v := range template.Parameters {
		name := v.Name
		if v.Alias != "" {
			name = v.Alias
		}
		// Cli can check required flag before make a request to backend, but API itself could not, so validate flags here
		flag := flagSet.Lookup(name)
		if name == "name" {
			continue
		}
		if flag == nil || flag.Value.String() == "" {
			if v.Required {
				return nil, fmt.Errorf("required flag(s) \"%s\" not set", name)
			}
			continue
		}

		switch v.Type {
		case cue.IntKind:
			workloadData[v.Name], err = flagSet.GetInt64(name)
		case cue.StringKind:
			workloadData[v.Name], err = flagSet.GetString(name)
		case cue.BoolKind:
			workloadData[v.Name], err = flagSet.GetBool(name)
		case cue.NumberKind, cue.FloatKind:
			workloadData[v.Name], err = flagSet.GetFloat64(name)
		}
		if err != nil {
			if strings.Contains(err.Error(), "of flag of type string") {
				data, _ := flagSet.GetString(name)
				switch v.Type {
				case cue.IntKind:
					workloadData[v.Name], err = strconv.ParseInt(data, 10, 64)
				case cue.BoolKind:
					workloadData[v.Name], err = strconv.ParseBool(data)
				case cue.NumberKind, cue.FloatKind:
					workloadData[v.Name], err = strconv.ParseFloat(data, 64)
				}
				if err != nil {
					return nil, fmt.Errorf("get flag(s) \"%s\" err %v", v.Name, err)
				}
				continue
			}
			return nil, fmt.Errorf("get flag(s) \"%s\" err %v", v.Name, err)
		}
	}
	if err = app.SetWorkload(workloadName, tp, workloadData); err != nil {
		return app, err
	}
	return app, app.Save(envName)
}

func BaseRun(staging bool, app *application.Application, kubeClient client.Client, Env *types.EnvMeta, io cmdutil.IOStreams) (string, error) {
	if staging {
		return "Staging saved", nil
	}
	if err := app.BuildRun(context.Background(), kubeClient, Env, io); err != nil {
		err = fmt.Errorf("create app err: %s", err)
		return "", err
	}
	return fmt.Sprintf("App %s deployed", app.Name), nil
}