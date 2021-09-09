package azureaks

import (
	"context"
	"log"

	"github.com/cucumber/godog"
	"github.com/probr/probr-pack-aks/internal/summary"

	"github.com/probr/probr-pack-aks/internal/common"
	"github.com/probr/probr-sdk/probeengine"
)

type scenarioState struct {
	common.ScenarioState
}

// ProbeStruct allows this probe to be added to the ProbeStore
type probeStruct struct {
}

// Probe ...
var Probe probeStruct
var scenario scenarioState // Local container of scenario state
var aksJSON []byte

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.Name = gs.Name
	s.Probe = summary.State.GetProbeLog(probeName)
	s.Audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.Ctx = context.Background()
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

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		afterScenario(scenario, probe, s, err)
	})

	ctx.BeforeStep(func(st *godog.Step) {
		scenario.CurrentStep = st.Text
	})

	ctx.AfterStep(func(st *godog.Step, err error) {
		scenario.CurrentStep = ""
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
