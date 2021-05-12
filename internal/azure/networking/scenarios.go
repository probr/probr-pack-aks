package azureaks

import (
	"github.com/citihub/probr-pack-aks/internal/common"
)

func (scenario *scenarioState) anAzureKubernetesClusterWeCanReadTheConfigurationOf() (err error) {
	aksJson, err = common.AnAzureKubernetesClusterWeCanReadTheConfigurationOf(scenario.GetScenarioState())
	return
}

func (scenario *scenarioState) privateClusterIsEnabled() error {
	return common.OPAProbe("private_cluster", aksJson, scenario.GetScenarioState())
}

func (scenario *scenarioState) networkOutboundType() error {
	return common.OPAProbe("network_outbound_type", aksJson, scenario.GetScenarioState())
}

func (scenario *scenarioState) cniNetworkingIsEnabled() error {
	return common.OPAProbe("network_policy", aksJson, scenario.GetScenarioState())
}

func (scenario *scenarioState) nodesDontHavePublicIps() error {
	return common.OPAProbe("node_public_ip", aksJson, scenario.GetScenarioState())
}
