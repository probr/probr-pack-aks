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

	regoFilePath := probeengine.GetFilePath("internal", "azure", "aks", "aks.rego")
	var r bool

	r, err = opa.Eval(regoFilePath, opaPackageName, functionName, &aksJson)

	if err != nil {
		return
	}

	if r == false {
		err = fmt.Errorf("Rego function %s returned result of 'non-compliant'", functionName)
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

func opaProbe(opaFuncName string, scenario *scenarioState) (err error) {
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

	stepTrace.WriteString("Use OPA to evaluate this cluster; ")

	err = eval(opaFuncName)
	return
}

func (scenario *scenarioState) azurePolicyIsEnabled() error {
	return opaProbe("azure_policy", scenario)
}

func (scenario *scenarioState) azureADIntegrationIsEnabled() error {
	return opaProbe("enable_rbac", scenario)
}

func (scenario *scenarioState) theKubernetesWebUIIsDisabled() error {
	return opaProbe("kube_dashboard", scenario)
}

func (scenario *scenarioState) privateClusterIsEnabled() error {
	return opaProbe("private_cluster", scenario)
}

func (scenario *scenarioState) networkOutboundType() error {
	return opaProbe("network_outbound_type", scenario)
}

func (scenario *scenarioState) diskEncryption() error {
	return opaProbe("disk_encryption", scenario)
}

func (scenario *scenarioState) cniNetworkingIsEnabled() error {
	return opaProbe("network_policy", scenario)
}
