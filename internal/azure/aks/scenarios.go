package azureaks

import (
	"fmt"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/citihub/probr-sdk/probeengine/opa"
	"github.com/citihub/probr-sdk/utils"
	"log"

	"strings"
)

const opaPackageName = "aks"

func eval(functionName string) (err error) {
	// true = return nil
	// false = return new err
	// err = return err

	regoFilePath := probeengine.GetFilePath("internal", "azure", "aks", "general.rego")
	var r bool

	r, err = opa.Eval(regoFilePath, opaPackageName, functionName, &aksJson)

	if err != nil {
		return
	}

	if r == false {
		err = fmt.Errorf("Rego function %s returned an error", functionName)
	} else {
		err = nil
	}
	return
}

func loadJSON() (err error) {
	// TODO: make this configurable
	aksJson, err = azConnection.GetManagedClusterJSON("probr-demo-rg", "probr-demo-cluster")
	s := string(aksJson)
	log.Printf("[DEBUG] AKS JSON: %v", s)
	return
}

func (scenario *scenarioState) anAzureKubernetesClusterWeCanReadTheConfigurationOf() (err error) {
	var stepTrace strings.Builder

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	payload = struct {
		Placeholder string
	}{
		"placeholder",
	}

	stepTrace.WriteString("Get the configuration of the AKS cluster; ")

	err = loadJSON()
	if err != nil {
		log.Printf("Error loading JSON: %v", err)
		return
	}

	if len(aksJson) == 0 {
		err = fmt.Errorf("aksJson empty")
	}
	return
}

func (scenario *scenarioState) azureADIntegrationIsEnabled() (err error) {
	var stepTrace strings.Builder

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	payload = struct {
		Placeholder string
	}{
		"placeholder",
	}

	stepTrace.WriteString("Use OPA to evaluate whether RBAC is enabled on this cluster; ")

	err = eval("enable_rbac")
	return
}

func (scenario *scenarioState) azurePolicyIsEnabled() (err error) {
	var stepTrace strings.Builder

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	payload = struct {
		Placeholder string
	}{
		"placeholder",
	}

	stepTrace.WriteString("Use OPA to evaluate whether Azure Policy is enabled on this cluster; ")

	err = eval("azure_policy")
	return
}

func (scenario *scenarioState) theKubernetesWebUIIsDisabled() (err error) {
	var stepTrace strings.Builder

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	payload = struct {
		Placeholder string
	}{
		"placeholder",
	}

	stepTrace.WriteString("Use OPA to evaluate whether Kube Dashboard is enabled on this cluster; ")

	err = eval("kube_dashboard")
	return
}
