package azureaks

import (
	"github.com/probr/probr-pack-aks/internal/common"
)

func (scenario *scenarioState) anAzureKubernetesClusterWeCanReadTheConfigurationOf() (err error) {

	aksJSON, err = common.AnAzureKubernetesClusterWeCanReadTheConfigurationOf(scenario.GetScenarioState())

	return
}

func (scenario *scenarioState) azurePolicyIsEnabled() error {
	return common.OPAProbe("azure_policy", aksJSON, scenario.GetScenarioState())
}

func (scenario *scenarioState) azureADIntegrationIsEnabled() error {
	return common.OPAProbe("enable_rbac", aksJSON, scenario.GetScenarioState())
}

func (scenario *scenarioState) theKubernetesWebUIIsDisabled() error {
	return common.OPAProbe("kube_dashboard", aksJSON, scenario.GetScenarioState())
}
