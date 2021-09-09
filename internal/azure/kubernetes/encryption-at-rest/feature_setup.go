package aksencryptionatrest

import (
	"context"
	"fmt"
	"log"

	"github.com/cucumber/godog"

	"github.com/probr/probr-pack-aks/internal/common"
	"github.com/probr/probr-pack-aks/internal/config"
	"github.com/probr/probr-pack-aks/internal/connection"
	"github.com/probr/probr-pack-aks/internal/summary"

	"github.com/probr/probr-sdk/probeengine"
)

type scenarioState struct {
	common.ScenarioState
	namespace string
	pods      []string
	pvcs      []string
	disks     []string
	//additional variables to hold state goes here
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
	s.namespace = config.Vars.ServicePacks.Kubernetes.ProbeNamespace
	probeengine.LogScenarioStart(gs)
}

// Name will return this probe's name
func (probe probeStruct) Name() string {
	return "encryption-at-rest"
}

// Path will return this probe's feature path
func (probe probeStruct) Path() string {
	return probeengine.GetFeaturePath("internal", "azure", "kubernetes", probe.Name())
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
	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, scenario.aKubernetesClusterIsDeployed)
	ctx.Step(`^I create a Pod which dynamically creates an Azure Disk$`, scenario.iCreateAPodWhichDynamicallyCreatesAnAzureDisk)
	ctx.Step(`^the disk is encrypted using Customer Managed Keys$`, scenario.theDiskIsEncryptedUsingCustomerManagedKeys)

	ctx.Step(`^an Azure Kubernetes cluster we can read the configuration of$`, scenario.anAzureKubernetesClusterWeCanReadTheConfigurationOf)
	ctx.Step(`^Disk Encryption is enabled$`, scenario.diskEncryption)

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
	//if config.Vars.ServicePacks.Kubernetes.KeepPods == "false" {
	for _, podName := range scenario.pods {
		err := connection.Kubernetes.DeletePodIfExists(podName, scenario.namespace, Probe.Name())
		if err != nil {
			log.Printf(fmt.Sprintf("[ERROR] Could not retrieve pod from namespace '%s' for deletion: %s", scenario.namespace, err))
		}
	}
	//}

	for _, pvcName := range scenario.pvcs {
		err := connection.Kubernetes.DeletePVCIfExists(pvcName, scenario.namespace, Probe.Name())
		if err != nil {
			log.Printf(fmt.Sprintf("[ERROR] Could not retrieve PVC from namespace '%s' for deletion: %s", scenario.namespace, err))
		}
	}

	log.Println("[DEBUG] Teardown completed")
}
