package cmd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/txn2/txeh"

	"github.com/spaceuptech/space-cli/model"
)

func getSpaceCloudIpTableDirectory() string {
	return fmt.Sprintf("%s/space-ip-table", getSpaceCloudDirectory())
}

func getSpaceCloudIpTablePath() string {
	return fmt.Sprintf("%s/hosts", getSpaceCloudIpTableDirectory())
}

func getSpaceCloudArtifactDirectory() string {
	return fmt.Sprintf("%s/space-artifact", getSpaceCloudDirectory())
}

func generateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 50
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String() // E.g. "ExcbsVQs"
}

// CodeSetup initializes development environment
func CodeSetup(username, key, url string, dev bool) error {
	// todo old keys always remain in accounts.yaml file

	createDirIfNotExist(getSpaceCloudIpTableDirectory())
	createDirIfNotExist(getSpaceCliDirectory())
	createDirIfNotExist(getSpaceCloudArtifactDirectory())
	// for now artifact.yaml need to be manually placed in this folder
	// then docker container will mount it

	log.Println("username pas", username, key, url, dev)
	if username == "" {
		username = generateRandomString()
		fmt.Printf("Your New Username: %s\n", username)
	}
	if key == "" {
		key = generateRandomString()
		fmt.Printf("Your New Key: %s\n", key)
	}

	selectedAccount := model.Account{
		ID:        username,
		UserName:  username,
		Key:       key,
		ServerUrl: url,
	}

	if err := checkCred(&selectedAccount); err != nil {
		logrus.Errorf("error in setup unable to check credentials - %v", err)
		return err
	}

	devMode := "false"
	if dev {
		devMode = "true" // even the flag set true in dev of container sc didn't start in prod mode
	}

	containersToCreate := []struct {
		dnsName        string
		containerImage string
		containerName  string
		envs           []string
		mount          []mount.Mount
		exposedPorts   nat.PortSet
		portMapping    nat.PortMap
	}{
		{
			// http://store.space-cloud.svc.cluster.local
			// gate way
			containerImage: "spaceuptech/gateway",
			containerName:  "space-cloud-gateway--service--v1--task",
			dnsName:        "gateway.space-cloud.svc.cluster.local",
			envs: []string{
				"ARTIFACT_ADDR=store.space-cloud.svc.cluster.local:4122",
				"RUNNER_ADDR=runner.space-cloud.svc.cluster.local:4050",
				"ADMIN_USER=" + username,
				"ADMIN_PASS=" + key,
				"DEV=" + devMode,
				// todo set admin-secret
			},
			exposedPorts: nat.PortSet{
				"4122": struct{}{},
			},
			portMapping: nat.PortMap{
				"4122": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "4122"}},
			},
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudIpTablePath(),
					Target: "/etc/hosts",
				},
			},
		},
		{
			// runner
			containerImage: "spaceuptech/runner",
			containerName:  "space-cloud-runner--service--v1--task",
			dnsName:        "runner.space-cloud.svc.cluster.local",
			envs: []string{
				"ARTIFACT_ADDR=store.space-cloud.svc.cluster.local:4122", // TODO Change the default value in runner it starts with http
				"DRIVER=docker",
			},
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudIpTablePath(),
					Target: "/etc/hosts",
				},
			},
			// todo set admin secret in envs
		},

		{
			// artifact store
			containerImage: "spaceuptech/gateway",
			containerName:  "space-cloud-artifact--service--v1--task",
			dnsName:        "store.space-cloud.svc.cluster.local",
			envs: []string{
				"CONFIG=/home/artifact.yaml",
			},
			// todo set admin secret in envs
			mount: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: getSpaceCloudIpTablePath(),
					Target: "/etc/hosts",
				},
				{
					Type:   mount.TypeBind, // mount artifact.yaml that is config file
					Source: getSpaceCloudArtifactDirectory(),
					Target: "/home",
				},
			},
		},
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Errorf("error cli setup unable to initialize docker client got error message - %v", err)
		return err
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		logrus.Errorf("error cli setup unable to load host file with suitable default got error message - %v", err)
		return err
	}
	// change the default host file location for crud operation to our specified path
	// default value /etc/hosts
	hosts.WriteFilePath = getSpaceCloudIpTablePath()
	if err := hosts.SaveAs(getSpaceCloudIpTablePath()); err != nil {
		logrus.Errorf("error cli setup unable to save as host file to specified path (%s) got error message - %v", getSpaceCloudIpTablePath(), err)
		return err
	}

	for _, c := range containersToCreate {
		// pull image from docker hub
		// out, err := cli.ImagePull(ctx, dockerImageSpaceCloud, types.ImagePullOptions{})
		// if err != nil {
		// 	logrus.Errorf("error cli setup unable to pull image from docker hub got error message - %v", err)
		// 	return err
		// }
		// io.Copy(os.Stdout, out)

		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image:        c.containerImage,
			ExposedPorts: c.exposedPorts,
			Env:          c.envs,
		}, &container.HostConfig{
			Mounts:       c.mount,
			PortBindings: c.portMapping,
		}, nil, c.containerName)
		if err != nil {
			logrus.Errorf("error cli setup unable to create container %s got error message  - %v", c.containerName, err)
			return err
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			logrus.Errorf("error cli setup unable to start container %s got error message - %v", c.containerName, err)
			return err
		}

		data, err := cli.ContainerInspect(ctx, c.containerName)
		if err != nil {
			logrus.Errorf("error cli setup unable to inspect container %s got error message - %v", c.containerName, err)
		}
		hosts.AddHost(data.NetworkSettings.IPAddress, c.dnsName)
	}

	if err := hosts.Save(); err != nil {
		logrus.Errorf("error cli setup unable to save host file got error message - %v", err)
		return err
	}
	return nil
}
