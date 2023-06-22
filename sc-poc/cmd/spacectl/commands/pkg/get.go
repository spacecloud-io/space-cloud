package pkg

import (
	"fmt"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newCommandGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("more than 1 or no resources specified")
			}

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
			resourceName := args[0]

			// Get all registered sources' GVR
			path := creds["url"] + "/sc/v1/sources"
			sourcesGVR, err := listAllSources(client, path)
			if err != nil {
				return err
			}

			var data [][]string
			if resourceName == "all" {
				for _, gvr := range sourcesGVR {
					path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/?package=%s", creds["url"], gvr.Group, gvr.Version, gvr.Resource, cfg.Name)
					unstrList, err := getResources(client, path)
					if err != nil {
						return err
					}

					for _, item := range unstrList.Items {
						data = append(data, []string{item.GetName(), gvr.Resource, "Active"})
					}
				}
				renderTable([]string{"Name", "Type", "Status"}, data)
				return nil
			}

			for _, gvr := range sourcesGVR {
				if resourceName == gvr.Resource {
					path := fmt.Sprintf("%s/sc/v1/config/%s/%s/%s/?package=%s", creds["url"], gvr.Group, gvr.Version, gvr.Resource, cfg.Name)
					unstrList, err := getResources(client, path)
					if err != nil {
						return err
					}

					for _, item := range unstrList.Items {
						data = append(data, []string{item.GetName(), "Active"})
					}

					renderTable([]string{"Name", "Status"}, data)
					return nil
				}
			}

			return fmt.Errorf("invalid resource name specified")
		},
	}

	return cmd
}

func renderTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
