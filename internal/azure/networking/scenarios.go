package azureaks

import (
	"github.com/citihub/probr-pack-aks/internal/common"
)

func (scenario *scenarioState) anAzureKubernetesClusterWeCanReadTheConfigurationOf() (err error) {

	baseScenario := scenario.GetScenarioState()
	aksJson, err = common.AnAzureKubernetesClusterWeCanReadTheConfigurationOf(&baseScenario)

	return
}

func (scenario *scenarioState) privateClusterIsEnabled() error {
	baseState := scenario.GetScenarioState()
	return common.OPAProbe("private_cluster", aksJson, &baseState)
}

func (scenario *scenarioState) networkOutboundType() error {
	baseState := scenario.GetScenarioState()
	return common.OPAProbe("network_outbound_type", aksJson, &baseState)
}

func (scenario *scenarioState) cniNetworkingIsEnabled() error {
	baseState := scenario.GetScenarioState()
	return common.OPAProbe("network_policy", aksJson, &baseState)
}

func (scenario *scenarioState) nodesDontHavePublicIps() error {
	baseState := scenario.GetScenarioState()
	return common.OPAProbe("node_public_ip", aksJson, &baseState)
}
