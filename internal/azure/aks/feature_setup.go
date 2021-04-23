package azureaks

import (
	"context"
	"log"

	"github.com/cucumber/godog"

	"github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/probeengine"
	azureutil "github.com/citihub/probr-sdk/providers/azure"
	"github.com/citihub/probr-sdk/providers/azure/connection"
)

type scenarioState struct {
	name        string
	currentStep string
	audit       *audit.ScenarioAudit
	probe       *audit.Probe
	ctx         context.Context
	tags        map[string]*string
	//additional variables to hold state goes here
}

// ProbeStruct allows this probe to be added to the ProbeStore
type probeStruct struct {
}

var Probe probeStruct
var scenario scenarioState // Local container of scenario state
var aksJson []byte
var azConnection connection.Azure // Provides functionality to interact with Azure

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.ctx = context.Background()
	probeengine.LogScenarioStart(gs)
}

// Name will return this probe's name
func (probe probeStruct) Name() string {
	return "aks"
}

// Path will return this probe's feature path
func (probe probeStruct) Path() string {
	return probeengine.GetFeaturePath("internal", "azure", probe.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (probe probeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {

	ctx.BeforeSuite(func() {

		azConnection = connection.NewAzureConnection(
			context.Background(),
			azureutil.SubscriptionID(),
			azureutil.TenantID(),
			azureutil.ClientID(),
			azureutil.ClientSecret(),
		)

	})

	ctx.AfterSuite(func() {
	})
}

// ScenarioInitialize initialises the scenario
func (probe probeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&scenario, probe.Name(), s)
	})

	// Background
	ctx.Step(`^an Azure Kubernetes cluster we can read the configuration of$`, scenario.anAzureKubernetesClusterWeCanReadTheConfigurationOf)

	// Steps
	ctx.Step(`^Azure AD integration is enabled$`, scenario.azureADIntegrationIsEnabled)
	ctx.Step(`^Azure Policy is enabled$`, scenario.azurePolicyIsEnabled)
	ctx.Step(`^the Kubernetes Web UI is disabled$`, scenario.theKubernetesWebUIIsDisabled)
	ctx.Step(`^Private Cluster is enabled$`, scenario.privateClusterIsEnabled)
	ctx.Step(`^Disk Encryption is enabled$`, scenario.diskEncryption)
	ctx.Step(`^outbound network routing is user controlled$`, scenario.networkOutboundType)
	ctx.Step(`^CNI network policy is enabled$`, scenario.cniNetworkingIsEnabled)
	ctx.Step(`^Kubernetes node hosts do not have public IPs$`, scenario.nodesDontHavePublicIps)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		afterScenario(scenario, probe, s, err)
	})

	ctx.BeforeStep(func(st *godog.Step) {
		scenario.currentStep = st.Text
	})

	ctx.AfterStep(func(st *godog.Step, err error) {
		scenario.currentStep = ""
	})
}

func afterScenario(scenario scenarioState, probe probeStruct, gs *godog.Scenario, err error) {

	teardown()
	probeengine.LogScenarioEnd(gs)
}

func teardown() {

	log.Printf("[DEBUG] Teardown - removing resources used during tests")

	//delete any resources you created here

	log.Println("[DEBUG] Teardown completed")
}
