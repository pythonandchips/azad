package runner

import (
	"fmt"

	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/steps"
)

func handleIncludes(includeStep steps.IncludesStep, store *contextStore) error {
	rolesVal, err := store.evalVariable(includeStep.Roles, allowList)
	if err != nil {
		return err
	}
	for _, roleName := range rolesVal.AsValueSlice() {
		logger.Info("Starting Role: %s", roleName.AsString())
		err := applyRole(roleName.AsString(), store)
		if err != nil {
			return err
		}
		logger.Info("Finished Role: %s", roleName.AsString())
	}
	return nil
}

func applyRole(roleName string, store *contextStore) error {
	roleContainer, err := roleList.FindByName(roleName)
	if err != nil {
		return fmt.Errorf("%s in %s", err, store.contextPath)
	}
	user := store.user()
	if roleContainer.User != nil {
		userValue, err := store.evalVariable(roleContainer.User, allowString)
		if err != nil {
			return err
		}
		user = userValue.AsString()
	}
	childStore := store.childStore(user, roleContainer.File)
	err = runForContext(roleContainer.Steps, &childStore)
	if err != nil {
		return err
	}
	return nil
}
