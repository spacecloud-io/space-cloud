package main

import (
	"github.com/urfave/cli"

	"github.com/spaceuptech/space-cli/cmd"
)

func actionGetGlobalConfig(c *cli.Context) error {
	return cmd.GetGlobalConfig(c)
}

func actionGetRemoteServices(c *cli.Context) error {
	return cmd.GetRemoteServices(c)
}

func actionGetAuthProviders(c *cli.Context) error {
	return cmd.GetAuthProviders(c)
}

func actionGetEventingTrigger(c *cli.Context) error {
	return cmd.GetEventingTrigger(c)
}

func actionGetEventingConfig(c *cli.Context) error {
	return cmd.GetEventingConfig(c)
}

func actionGetEventingSchema(c *cli.Context) error {
	return cmd.GetEventingSchema(c)
}

func actionGetEventingSecurityRule(c *cli.Context) error {
	return cmd.GetEventingSecurityRule(c)
}

func actionGetFileStoreConfig(c *cli.Context) error {
	return cmd.GetFileStoreConfig(c)
}

func actionGetFileStoreRule(c *cli.Context) error {
	return cmd.GetFileStoreRule(c)
}

func actionGetDbRule(c *cli.Context) error {
	return cmd.GetDbRule(c)
}

func actionGetDbConfig(c *cli.Context) error {
	return cmd.GetDbConfig(c)
}

func actionGetDbSchema(c *cli.Context) error {
	return cmd.GetDbSchema(c)
}

func actionGetLetsEncryptDomain(c *cli.Context) error {
	return cmd.GetLetsEncryptDomain(c)
}

func actionGetRoutes(c *cli.Context) error {
	return cmd.GetRoutes(c)
}
