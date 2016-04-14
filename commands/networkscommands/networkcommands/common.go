package networkcommands

import (
	"strings"

	"github.com/rackspace/rack/internal/github.com/fatih/structs"
	osNetworks "github.com/rackspace/rack/internal/github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

func networkSingle(network *osNetworks.Network) map[string]interface{} {
	m := structs.Map(network)
	m["Up"] = m["AdminStateUp"]
	if subnets, ok := m["Subnets"].([]string); ok {
		m["Subnets"] = strings.Join(subnets, ",")
	}
	return m
}
