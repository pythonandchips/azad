package parser

import "github.com/pythonandchips/azad/azad"

func unpackServer(serverDescriptions []serverDescription) ([]azad.Server, error) {
	servers := []azad.Server{}
	for _, serverDescription := range serverDescriptions {
		server := azad.Server{
			Group:     serverDescription.Group,
			Addresses: serverDescription.Addresses,
		}
		servers = append(servers, server)
	}
	return servers, nil
}
