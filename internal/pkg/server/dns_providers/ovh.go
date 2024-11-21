package dns_providers

import (
	"log"

	"github.com/bmehdi777/moon/internal/pkg/server/config"
	"github.com/ovh/go-ovh/ovh"
)

type OVHProvider struct {
	client *ovh.Client
}

func NewClient() *OVHProvider {
	var provider OVHProvider

	client, err := ovh.NewClient(
		config.GlobalConfig.OvhProvider.Endpoint,
		config.GlobalConfig.OvhProvider.ApplicationKey,
		config.GlobalConfig.OvhProvider.ApplicationSecret,
		config.GlobalConfig.OvhProvider.ConsumerKey,
	)
	if err != nil {
		log.Fatalf("Can't initialize OVH client : %v", err)
	}

	provider = OVHProvider{
		client: client,
	}

	return &provider
}
