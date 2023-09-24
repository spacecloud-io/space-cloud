package pkg

import (
	"log"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	clientutils "github.com/spacecloud-io/space-cloud/utils/client"
)

func newCommandGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				log.Fatal("Invalid argument: more than 1 or no resources specified")
			}

			httpClient := &http.Client{}
			creds, err := clientutils.GetCredentials()
			if err != nil {
				log.Fatal("Failed to get SpaceCloud credentials: ", err)
			}

			// Login to SpaceCloud
			if _, err := clientutils.Login(httpClient, creds); err != nil {
				log.Fatal("Failed to authenticate with SpaceCloud: ", err)
			}

			cfg := ReadConfig()
			resourceName := args[0]

			// Get all registered sources' GVR
			sourcesGVR, err := clientutils.ListAllSources(httpClient, creds.BaseUrl)
			if err != nil {
				log.Fatal("Failed to list all registered sources: ", err)
			}

			var data [][]string
			if resourceName == "all" {
				for _, gvr := range sourcesGVR {
					unstrList, err := clientutils.GetResources(httpClient, gvr, creds.BaseUrl, cfg.Name)
					if err != nil {
						log.Fatal("Failed to get resources: ", err)
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
					unstrList, err := clientutils.GetResources(httpClient, gvr, creds.BaseUrl, cfg.Name)
					if err != nil {
						log.Fatal("Failed to get resources: ", err)
					}

					for _, item := range unstrList.Items {
						data = append(data, []string{item.GetName(), "Active"})
					}

					renderTable([]string{"Name", "Status"}, data)
					return nil
				}
			}

			log.Fatal("Invalid resource name specified")
			return nil
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
