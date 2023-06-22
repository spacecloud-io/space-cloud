package pkg

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spacecloud-io/space-cloud/utils"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func newCommandApply() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "apply",
		Aliases: []string{"a"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client := &http.Client{}
			creds, err := getCredentials()
			if err != nil {
				return err
			}

			// Login to SpaceCloud
			if err := login(client, creds); err != nil {
				return err
			}

			cfg := ReadConfig()
			m := make(map[schema.GroupVersionResource][]string)

			// Get all registered sources' GVR
			path := creds["url"] + "/sc/v1/sources"
			sourcesGVR, err := listAllSources(client, path)
			if err != nil {
				return err
			}

			// Get the resources present in the SpaceCloud
			for _, gvr := range sourcesGVR {
				path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/?package=%s", creds["url"], gvr.Group, gvr.Version, gvr.Resource, cfg.Name)
				resources, err := getResources(client, path)
				if err != nil {
					return err
				}
				// Cache the resource's name in a map
				m[gvr] = make([]string, 0)

				for _, obj := range resources.Items {
					resourceName := obj.GetName()
					m[gvr] = append(m[gvr], resourceName)
				}
			}

			// Read resources from the local package.
			resDir := cfg.ResourceDir
			files, err := os.ReadDir(resDir)
			if err != nil {
				return err
			}

			for _, file := range files {
				arr, err := utils.ReadSpecObjectsFromFile(filepath.Join(resDir, file.Name()))
				if err != nil {
					return err
				}

				for _, spec := range arr {
					gvr := schema.GroupVersionResource{
						Group:    spec.GroupVersionKind().Group,
						Version:  spec.GroupVersionKind().Version,
						Resource: utils.Pluralize(spec.GetKind())}
					name := spec.GetName()

					// If resource exists in SpaceCloud, remove from the cache.
					index := findElement(m[gvr], name)
					if index != -1 {
						m[gvr] = DeleteElement(m[gvr], index)
					}

					// Inject the labels into the spec
					defaultLabel := map[string]string{
						"space-cloud.io/package": cfg.Name,
					}
					spec.SetLabels(defaultLabel)

					if cfg.Labels != nil {
						spec.SetLabels(cfg.Labels)
					}

					// Perform apply operation
					path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/", creds["url"], gvr.Group, gvr.Version, gvr.Resource)
					err := applyResources(client, path, spec)
					if err != nil {
						return err
					}

				}

				// Delete the resources in SpaceCloud which are still present in cache
				for gvr, names := range m {
					for _, name := range names {
						path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/%s", creds["url"], gvr.Group, gvr.Version, gvr.Resource, name)
						err := deleteResources(client, path)
						if err != nil {
							return err
						}
					}
				}
			}
			return nil
		},
	}

	return cmd
}
