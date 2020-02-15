package main

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/cmd/objects"
)

func actionGetGlobalConfig(c *cli.Context) error {
	return objects.GetGlobalConfig(c)
}

func actionGetRemoteServices(c *cli.Context) error {
	return objects.GetRemoteServices(c)
}

func actionGetAuthProviders(c *cli.Context) error {
	return objects.GetAuthProviders(c)
}

func actionGetEventingTrigger(c *cli.Context) error {
	return objects.GetEventingTrigger(c)
}

func actionGetEventingConfig(c *cli.Context) error {
	return objects.GetEventingConfig(c)
}

func actionGetEventingSchema(c *cli.Context) error {
	return objects.GetEventingSchema(c)
}

func actionGetEventingSecurityRule(c *cli.Context) error {
	return objects.GetEventingSecurityRule(c)
}

func actionGetFileStoreConfig(c *cli.Context) error {
	return objects.GetFileStoreConfig(c)
}

func actionGetFileStoreRule(c *cli.Context) error {
	return objects.GetFileStoreRule(c)
}

func actionGetDbRule(c *cli.Context) error {
	return objects.GetDbRule(c)
}

func actionGetDbConfig(c *cli.Context) error {
	return objects.GetDbConfig(c)
}

func actionGetDbSchema(c *cli.Context) error {
	return objects.GetDbSchema(c)
}

func actionGetLetsEncryptDomain(c *cli.Context) error {
	return objects.GetLetsEncryptDomain(c)
}

func actionGetRoutes(c *cli.Context) error {
	return objects.GetRoutes(c)
}
