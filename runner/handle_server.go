package runner

import (
	"github.com/pythonandchips/azad/steps"
)

func handleServer(serverStep steps.ServerStep, store serverStore) error {
	addressesValue, err := store.evalVariable(serverStep.Addresses, allowList)
	if err != nil {
		return err
	}
	servers := servers{}
	for _, addressValue := range addressesValue.AsValueSlice() {
		servers = append(servers, server{
			address: addressValue.AsString(),
			group:   []string{serverStep.Name},
		})
	}
	store.addServers(servers)
	return nil
}
