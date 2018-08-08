package runner

import (
	"path/filepath"
	"strings"

	"github.com/pythonandchips/azad/logger"
	"github.com/pythonandchips/azad/parser"
	"github.com/pythonandchips/azad/steps"
)

var parsePlaybook = func(playbookFilePath string) (steps.PlaybookSteps, error) {
	return parser.PlaybookSteps(playbookFilePath)
}

var roleList steps.RoleContainers

// RunPlaybook run playbook with given path and config
var RunPlaybook = func(playbookFilePath string, globalConfig Config) error {
	logger.Info("Starting playbook run with %s", playbookFilePath)

	playbookSteps, err := parsePlaybook(playbookFilePath)
	if err != nil {
		logger.ErrorAndExit("Playbook has invalid syntax: %s", err)
	}
	roleList = playbookSteps.RoleList
	roleNames := []string{}
	for _, role := range roleList {
		roleNames = append(roleNames, role.Name)
	}
	logger.Debug("Available roles %s", strings.Join(roleNames, ", "))
	err = validateSteps(playbookSteps.StepList)
	if err != nil {
		logger.ErrorAndExit("Playbook is invalid: %s", err)
	}
	playbookPath, _ := filepath.Abs(filepath.Dir(playbookFilePath))
	globalStore := &store{
		config: globalConfig,
		path:   playbookPath,
	}
	for _, step := range playbookSteps.StepList {
		switch val := step.(type) {
		case steps.ServerStep:
			if err = handleServer(val, globalStore); err != nil {
				return err
			}
		case steps.InventoryStep:
			if err = handleInventory(val, globalStore); err != nil {
				return err
			}
		case steps.VariableStep:
			if err = handleVariable(val, globalStore); err != nil {
				return err
			}
		case steps.ContextContainer:
			if err = handleContainerContext(val, globalStore); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleContainerContext(contextContainer steps.ContextContainer, store *store) error {
	applyToValue, evalErr := store.evalVariable(contextContainer.ApplyTo, allowList)
	if evalErr != nil {
		return evalErr
	}
	applyTo := []string{}
	for _, val := range applyToValue.AsValueSlice() {
		applyTo = append(applyTo, val.AsString())
	}
	user := ""
	if contextContainer.User != nil {
		userValue, evalErr := store.evalVariable(contextContainer.User, allowString)
		if evalErr != nil {
			return evalErr
		}
		user = userValue.AsString()
	}
	contextStore, contextErr := store.contextStore(applyTo, user, store.path)
	defer func() {
		contextStore.closeConnections()
	}()
	if contextErr != nil {
		return contextErr
	}
	logger.Info("Starting running context %s", contextContainer.Name)
	err := runForContext(contextContainer.Steps, &contextStore)
	if err != nil {
		logger.Error("Failed to apply context %s:", contextContainer.Name)
		logger.Error("Error message: %s", err)
		return err
	}
	logger.Info("Finished running context %s", contextContainer.Name)
	return nil
}

func validateSteps(steps.StepList) error {
	return nil
}

func runForContext(steplist steps.StepList, store *contextStore) error {
	for _, step := range steplist {
		switch val := step.(type) {
		case steps.VariableStep:
			if err := handleVariable(val, store); err != nil {
				return err
			}
		case steps.TaskStep:
			if err := handleTask(val, store); err != nil {
				return err
			}
		case steps.IncludesStep:
			if err := handleIncludes(val, store); err != nil {
				return err
			}
		}
	}
	return nil
}
