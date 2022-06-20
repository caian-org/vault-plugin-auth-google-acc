package main

import (
	"os"

	googleAccountAuth "github.com/caian-org/vault-plugin-auth-google-acc"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	logger := hclog.New(&hclog.LoggerOptions{})

	flags := apiClientMeta.FlagSet()
	if err := flags.Parse(os.Args[1:]); err != nil {
		logger.Error("could not parse flags", "error", err)
		os.Exit(1)
	}

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: googleAccountAuth.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	})

	if err != nil {
		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}