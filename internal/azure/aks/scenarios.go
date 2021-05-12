package azureaks

import (
	"github.com/citihub/probr-pack-aks/internal/common"
)

func (scenario *scenarioState) anAzureKubernetesClusterWeCanReadTheConfigurationOf() (err error) {

	aksJson, err = common.AnAzureKubernetesClusterWeCanReadTheConfigurationOf(scenario.GetScenarioState())

	return
}

func (scenario *scenarioState) azurePolicyIsEnabled() error {
	return common.OPAProbe("azure_policy", aksJson, scenario.GetScenarioState())
}

func (scenario *scenarioState) azureADIntegrationIsEnabled() error {
	return common.OPAProbe("enable_rbac", aksJson, scenario.GetScenarioState())
}

func (scenario *scenarioState) theKubernetesWebUIIsDisabled() error {
	return common.OPAProbe("kube_dashboard", aksJson, scenario.GetScenarioState())
}
