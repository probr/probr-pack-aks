package aksencryptionatrest

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"log"

	"github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/probeengine"
	azureutil "github.com/citihub/probr-sdk/providers/azure"
	azureconnection "github.com/citihub/probr-sdk/providers/azure/connection"
	k8sconnection "github.com/citihub/probr-sdk/providers/kubernetes/connection"
)

type scenarioState struct {
	name        string
	currentStep string
	audit       *audit.ScenarioAudit
	probe       *audit.Probe
	ctx         context.Context
	tags        map[string]*string
	namespace   string
	pods        []string
	pvcs        []string
	disks       []string
	//additional variables to hold state goes here
}

// ProbeStruct allows this probe to be added to the ProbeStore
type probeStruct struct {
}

var Probe probeStruct
var scenario scenarioState               // Local container of scenario state
var azConnection azureconnection.Azure   // Provides functionality to interact with Azure
var kConnection k8sconnection.Connection // Provides functionality to interact with Kubernetes

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.ctx = context.Background()
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

		azConnection = azureconnection.NewAzureConnection(
			context.Background(),
			azureutil.SubscriptionID(),
			azureutil.TenantID(),
			azureutil.ClientID(),
			azureutil.ClientSecret(),
		)

		kConnection = k8sconnection.Get()

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
	//if config.Vars.ServicePacks.Kubernetes.KeepPods == "false" {
	for _, podName := range scenario.pods {
		err := kConnection.DeletePodIfExists(podName, scenario.namespace, Probe.Name())
		if err != nil {
			log.Printf(fmt.Sprintf("[ERROR] Could not retrieve pod from namespace '%s' for deletion: %s", scenario.namespace, err))
		}
	}
	//}

	for _, pvcName := range scenario.pvcs {
		err := kConnection.DeletePVCIfExists(pvcName, scenario.namespace, Probe.Name())
		if err != nil {
			log.Printf(fmt.Sprintf("[ERROR] Could not retrieve PVC from namespace '%s' for deletion: %s", scenario.namespace, err))
		}
	}

	log.Println("[DEBUG] Teardown completed")
}