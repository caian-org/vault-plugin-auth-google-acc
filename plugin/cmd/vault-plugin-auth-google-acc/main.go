package main

import (
	"log"
	"os"

	googleAccountAuth "github.com/caian-org/vault-plugin-auth-google-acc"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal("could not parse flags", err)
	}

	pluginOpts := &plugin.ServeOpts{
		BackendFactoryFunc: googleAccountAuth.Factory,
		TLSProviderFunc:    api.VaultPluginTLSProvider(apiClientMeta.GetTLSConfig()),
	}

	err := plugin.Serve(pluginOpts)
	if err != nil {
		log.Fatal("plugin shutting down", err)
	}
}
