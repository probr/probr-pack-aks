package azureaks

import (
	"github.com/probr/probr-pack-aks/internal/common"
)

func (scenario *scenarioState) anAzureKubernetesClusterWeCanReadTheConfigurationOf() (err error) {
	aksJSON, err = common.AnAzureKubernetesClusterWeCanReadTheConfigurationOf(scenario.GetScenarioState())
	return
}

func (scenario *scenarioState) privateClusterIsEnabled() error {
	return common.OPAProbe("private_cluster", aksJSON, scenario.GetScenarioState())
}

func (scenario *scenarioState) networkOutboundType() error {
	return common.OPAProbe("network_outbound_type", aksJSON, scenario.GetScenarioState())
}

func (scenario *scenarioState) cniNetworkingIsEnabled() error {
	return common.OPAProbe("network_policy", aksJSON, scenario.GetScenarioState())
}

func (scenario *scenarioState) nodesDontHavePublicIps() error {
	return common.OPAProbe("node_public_ip", aksJSON, scenario.GetScenarioState())
}
