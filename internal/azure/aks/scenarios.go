package azureaks

import (
	"github.com/citihub/probr-pack-aks/internal/common"
)

func (scenario *scenarioState) anAzureKubernetesClusterWeCanReadTheConfigurationOf() (err error) {

	baseScenario := scenario.GetScenarioState()
	aksJson, err = common.AnAzureKubernetesClusterWeCanReadTheConfigurationOf(&baseScenario)

	return
}

func (scenario *scenarioState) azurePolicyIsEnabled() error {
	baseState := scenario.GetScenarioState()
	return common.OPAProbe("azure_policy", aksJson, &baseState)
}

func (scenario *scenarioState) azureADIntegrationIsEnabled() error {
	baseState := scenario.GetScenarioState()
	return common.OPAProbe("enable_rbac", aksJson, &baseState)
}

func (scenario *scenarioState) theKubernetesWebUIIsDisabled() error {
	baseState := scenario.GetScenarioState()
	return common.OPAProbe("kube_dashboard", aksJson, &baseState)
}
