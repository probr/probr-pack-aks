// Package iam provides the implementation required to execute the BDD tests described in iam.feature file
package iam

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"

	"github.com/probr/probr-pack-aks/internal/common"
	"github.com/probr/probr-pack-aks/internal/config"
	"github.com/probr/probr-pack-aks/internal/connection"
	"github.com/probr/probr-pack-aks/internal/summary"

	"github.com/probr/probr-sdk/audit"
	"github.com/probr/probr-sdk/probeengine"

	k8sConnection "github.com/probr/probr-sdk/providers/kubernetes/connection"
	"github.com/probr/probr-sdk/providers/kubernetes/constructors"
	"github.com/probr/probr-sdk/providers/kubernetes/errors"
	"github.com/probr/probr-sdk/utils"
)

type probeStruct struct{}

// scenarioState holds the steps and state for any scenario in this probe
type scenarioState struct {
	common.ScenarioState
	namespace             string
	probeAudit            *audit.Probe
	pods                  []string //All pods created within the test. Should tear down at the end.
	micPodName            string
	azureIdentityBindings []string //Identity Bindings created within the test. Should tear down at the end.
}

// Probe meets the service pack interface for adding the logic from this file
var Probe probeStruct
var scenario scenarioState

func (scenario *scenarioState) aKubernetesClusterIsDeployed() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()
	stepTrace.WriteString(fmt.Sprintf("Validate that a cluster can be reached using the specified kube config and context; "))

	payload = struct {
		KubeConfigPath string
		KubeContext    string
	}{
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath,
		config.Vars.ServicePacks.Kubernetes.KubeContext,
	}

	err = connection.Kubernetes.ClusterIsDeployed() // Must be assigned to 'err' be audited
	return err
}

func (scenario *scenarioState) aResourceTypeXCalledYExistsInNamespaceCalledZ(resourceType string, resourceName string, namespace string) error {
	// Supported values for resourceType:
	//  'AzureIdentity'
	//  'AzureIdentityBinding'
	//
	// Supported values for resourceName:
	//  A string representing either an existing Azure Identity or Azure Identity Binding in K8s cluster
	//
	// Supported values for namespace:
	//	A string representing an existing namespace in K8s cluster

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	// TODO: This implementation is coupled to Azure. How should we deal with this when segregating service pack?

	var foundInNamespace bool
	var resource k8sConnection.APIResource
	var findErr error
	// Validate input
	switch resourceType {
	case "AzureIdentity":
		stepTrace.WriteString(fmt.Sprintf(
			"Retrieve Azure Identities from cluster; "))
		foundInNamespace, resource, findErr = azureIdentityExistsInNamespace(resourceName, namespace)
	case "AzureIdentityBinding":
		stepTrace.WriteString(fmt.Sprintf(
			"Retrieve Azure Identity Bindings from cluster; "))
		foundInNamespace, resource, findErr = azureIdentityBindingExistsInNamespace(resourceName, namespace)
	default:
		err = utils.ReformatError("Unexpected value provided for resourceType: %s", resourceType)
		return err
	}

	if findErr != nil {
		err = findErr
		return err
	}

	stepTrace.WriteString(fmt.Sprintf(
		"Check that %s '%s' exists in namespace '%s'; ", resourceType, resourceName, namespace))
	if !foundInNamespace {
		err = utils.ReformatError("%s '%s' was not found in namespace '%s'; ", resourceType, resourceName, namespace)
	}

	payload = struct {
		CustomResourceType string
		CustomResourceName string
		Resource           k8sConnection.APIResource
	}{
		CustomResourceType: resourceType,
		CustomResourceName: resourceName,
		Resource:           resource,
	}

	return err
}

func (scenario *scenarioState) iSucceedToCreateASimplePodInNamespaceAssignedWithThatAzureIdentityBinding(namespace, aibName string) error {
	// Supported values for namespace:
	//	'the probr'
	//	'the default'
	//
	// Supported values for aibName:
	//	'probr-aib'

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	// Validate input
	switch aibName {
	case "probr-aib":
	default:
		err = utils.ReformatError("Unexpected value provided for aibName: %s", aibName)
		return err
	}

	var aadPodIDBinding string

	// Validate input
	switch namespace {
	case "the probr":
		scenario.namespace = config.Vars.ServicePacks.Kubernetes.ProbeNamespace
		aadPodIDBinding = aibName // TODO: This value is the same in both config and feature file
	case "the default":
		scenario.namespace = "default"
		aadPodIDBinding = config.Vars.ServicePacks.AKS.ManagedID.DefaultNamespaceAIB // TODO: This value is the same in both config and feature file
	default:
		err = utils.ReformatError("Unexpected value provided for namespace: %s", namespace)
		return err
	}

	// TODO: This implementation is coupled with specific external cluster configuration, such as creation of specific namespace and aadidentitybinding
	// This is prone to error if not configured correctly.
	// Should revisit how to handle this.

	stepTrace.WriteString(fmt.Sprintf("Build a pod spec with default values; "))
	podObject := constructors.PodSpec(Probe.Name(), config.Vars.ServicePacks.Kubernetes.ProbeNamespace, config.Vars.ServicePacks.Kubernetes.AuthorisedContainerImage)
	// TODO: Delete iam-azi-test-aib-curl.yaml file from 'assets' folder

	stepTrace.WriteString(fmt.Sprintf("Add '%s' namespace to pod spec; ", scenario.namespace))
	podObject.Namespace = scenario.namespace

	stepTrace.WriteString(fmt.Sprintf("Add 'aadpodidbinding':'%s' label to pod spec; ", aadPodIDBinding))
	// For a pod to use AAD pod-managed identity, the pod needs an aadpodidbinding label with a value that matches a selector from a AzureIdentityBinding.
	// Ref: https://docs.microsoft.com/en-us/azure/aks/use-azure-ad-pod-identity
	podObject.Labels["aadpodidbinding"] = aadPodIDBinding

	stepTrace.WriteString(fmt.Sprintf("Create pod from spec; "))
	createdPodObject, creationErr := scenario.createPodfromObject(podObject)

	stepTrace.WriteString("Validate pod creation succeeds; ")
	if creationErr != nil {
		err = utils.ReformatError("Pod creation did not succeed: %v", creationErr)
	}

	payload = struct {
		Namespace      string
		AADPodIdentity string
		RequestedPod   *apiv1.Pod
		CreatedPod     *apiv1.Pod
		CreationError  error
	}{
		Namespace:      scenario.namespace,
		AADPodIdentity: aadPodIDBinding,
		RequestedPod:   podObject,
		CreatedPod:     createdPodObject,
		CreationError:  creationErr,
	}

	return err
}

func (scenario *scenarioState) anAttemptToObtainAnAccessTokenFromThatPodShouldX(expectedResult string) error {
	// Supported values for expectedResult:
	//	'Fail'
	//	'Succeed'

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	// Validate input
	var shouldReturnToken bool
	switch expectedResult {
	case "Fail":
		shouldReturnToken = false
	case "Succeed":
		shouldReturnToken = true
	default:
		err = utils.ReformatError("Unexpected value provided for expectedResult: %s", expectedResult)
		return err
	}

	// Guard clause: Ensure pod was created in previous step
	if len(scenario.pods) == 0 {
		err = utils.ReformatError("Pod failed to create in the previous step")
		return err
	}

	podName := scenario.pods[0]

	// Mechanism to get access token is executing a curl command on the pod
	// TODO: Clarify this and remove hardcoded IP
	// This is taking a long time, and failing in most cases.
	cmd := "curl http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https%3A%2F%2Fmanagement.azure.com%2F -H Metadata:true -s"

	stepTrace.WriteString(fmt.Sprintf("Attempt to run command in the pod: '%s'; ", cmd))
	_, stdOut, _, cmdErr := connection.Kubernetes.ExecCommand(cmd, scenario.namespace, podName)

	// Validate that no internal error occurred during execution of curl command
	if cmdErr != nil {
		err = utils.ReformatError("Error raised when attempting to execute curl command inside container: %v", cmdErr)
		return err
	}

	stepTrace.WriteString("Attempt to extract access token from command output; ")
	var accessToken struct {
		AccessToken string `json:"access_token,omitempty"`
	}
	jsonConvertErr := json.Unmarshal([]byte(stdOut), &accessToken)

	switch shouldReturnToken {
	case true:
		stepTrace.WriteString("Validate token was found; ")
		if jsonConvertErr != nil {
			err = utils.ReformatError("Failed to acquire token on pod %v Error: %v StdOut: %s", podName, jsonConvertErr, stdOut) //TODO: Error is being raised (see audit log)
		}
	case false:
		stepTrace.WriteString("Validate no token was found; ") //TODO: This is a potential false positve, since an error is raised by curl command (see audit log)
		if jsonConvertErr != nil && &accessToken.AccessToken != nil && len(accessToken.AccessToken) > 0 {
			err = utils.ReformatError("Token was successfully acquired on pod %v (result: %v)", podName, accessToken.AccessToken) //TODO: Adding access token to audit log until it can be tested. Remove afterwards for security reasons.
		}
	}

	return err
}

func (scenario *scenarioState) iCreateAnAzureIdentityBindingCalledInANondefaultNamespace(aibName, aiName string) error {
	// Supported values for aibName:
	//	A string representing an Azure Identity Binding to be created in K8s cluster
	//
	// Supported values for aibName:
	//	A string representing an Azure Identity to be created in K8s cluster

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	probrNameSpace := scenario.namespace

	aibName = aibName + "-test-test-test"
	stepTrace.WriteString(fmt.Sprintf(
		"Attempt to create '%s' binding in '%s' namespace bound to '%s' identity; ", aibName, probrNameSpace, aiName))
	createdAIB, err := azureCreateAIB(probrNameSpace, aibName, aiName) // create an AIB in a non-default NS if it doesn't already exist
	if err != nil {
		err = utils.ReformatError("An error occurred while creating '%s' binding: %v", aibName, err)
		log.Print(err)
	}
	// TODO:
	//	- Delete AIB at the end of test scenario
	scenario.azureIdentityBindings = append(scenario.azureIdentityBindings, aibName)

	payload = struct {
		Namespace                   string
		AzureIdentityBindingName    string
		AzureIdentityName           string
		CreatedAzureIdentityBinding k8sConnection.APIResource
	}{
		Namespace:                   probrNameSpace,
		AzureIdentityBindingName:    aibName,
		AzureIdentityName:           aiName,
		CreatedAzureIdentityBinding: createdAIB,
	}

	return err
}

func (scenario *scenarioState) theClusterHasManagedIdentityComponentsDeployed() error {

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	identityPodsNamespace := config.Vars.ServicePacks.AKS.ManagedID.IdentityNamespace
	stepTrace.WriteString(fmt.Sprintf(
		"Get pods from '%s' namespace; ", identityPodsNamespace))
	// look for the mic pods
	podList, getErr := connection.Kubernetes.GetPodsByNamespace(identityPodsNamespace)

	if getErr != nil {
		err = utils.ReformatError("An error occurred when trying to retrieve pods %v", err)
		return err
	}

	var micPodName string

	stepTrace.WriteString("Validate that at least one pod contains Label:'app.kubernetes.io/component=mic'; ")
	for _, pod := range podList.Items {
		if pod.Labels["app.kubernetes.io/component"] == "mic" {
			micPodName = pod.Name
			break
		}
	}

	if micPodName == "" {
		err = utils.ReformatError("No MIC pod found")
		return err
	}
	scenario.micPodName = micPodName

	payload = struct {
		IdentityPodsNamespace string
		MicPod                string
	}{
		IdentityPodsNamespace: identityPodsNamespace,
		MicPod:                micPodName,
	}

	return err
}

func (scenario *scenarioState) theExecutionOfAXCommandInsideTheMICPodIsY(commandType, result string) error {

	// Supported values for commandType:
	//	'get-azure-credentials'
	//
	// Supported values for result:
	//	'not allowed'

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	var cmd string
	// Validate input
	switch commandType {
	case "get-azure-credentials":
		cmd = "cat /etc/kubernetes/azure.json"
	default:
		err = utils.ReformatError("Unexpected value provided for commandType: %s", commandType)
		return err
	}

	var expectedExitCode int
	// Validate input
	switch result {
	case "not allowed":
		expectedExitCode = 126
	default:
		err = utils.ReformatError("Unexpected value provided for result: %s", result)
		return err
	}

	// Guard clause: Ensure the mic pod was found and stored in previous step
	if scenario.micPodName == "" {
		err = utils.ReformatError("MIC pod was not found in the previous step")
		return err
	}

	identityPodsNamespace := config.Vars.ServicePacks.AKS.ManagedID.IdentityNamespace
	stepTrace.WriteString(fmt.Sprintf(
		"Attempt to execute command '%s' in MIC pod '%s'; ", cmd, scenario.micPodName))
	exitCode, stdOut, _, cmdErr := connection.Kubernetes.ExecCommand(cmd, identityPodsNamespace, scenario.micPodName)

	// Validate that no internal error occurred during execution of curl command
	if cmdErr != nil && exitCode == -1 {
		err = utils.ReformatError("Error raised when attempting to execute command inside container: %v", cmdErr)
		return err
	}

	payload = struct {
		MICPodName       string
		Namespace        string
		Command          string
		ExpectedExitCode int
		ExitCode         int
		StdOut           string
	}{
		MICPodName:       scenario.micPodName,
		Namespace:        identityPodsNamespace,
		Command:          cmd,
		ExpectedExitCode: expectedExitCode,
		ExitCode:         exitCode,
		StdOut:           stdOut,
	}

	// TODO: Review this
	// I think ANY command executed against MIC pod will return same 126 exit code.
	// Potential cause is that 'mic' container image is based on 'distroless', which is
	// a minimalistic OS version containing only application and runtime dependencies.
	// It doesn't even have a shell.
	// If this assumption is correct, then this step and associated scenario should be
	// adjusted, since execution will always fail regardless of command, therefore causing
	// false positives.
	// Is there any other way to get access to Volume: /etc/kubernetes/azure.json ?
	// Ref:
	//  https://hub.docker.com/_/microsoft-k8s-aad-pod-identity-mic?tab=description
	//  https://github.com/GoogleContainerTools/distroless

	stepTrace.WriteString("Check expected exit code from command execution; ")
	if exitCode != expectedExitCode {
		err = utils.ReformatError("Unexpected exit code: %d Error: %v", exitCode, cmdErr)
		return err
	}

	return err
}

func (scenario *scenarioState) checkClusterRBACForAdminRole() error {
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	//this is the role definition name for rolename "Azure Kubernetes Service Cluster Admin Role"
	roleDefName := "0ab0b1a8-8aac-4efd-b8c2-3ee1fb270be8"

	caRoleAssigned, err := connection.Azure.ClusterHasRoleAssignment(config.Vars.ServicePacks.AKS.ResourceGroupName, config.Vars.ServicePacks.AKS.ClusterName, roleDefName)

	if err != nil {
		return err
	}

	if caRoleAssigned == true {
		return utils.ReformatError("Azure Kubernetes Service Cluster Admin Role is assigned to cluster")
	}

	payload = struct {
		Placeholder string
	}{
		Placeholder: "placeholder",
	}

	return nil
}

func (scenario *scenarioState) checkCannotObtainClusterAdminCredentials() error {
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	_, credsErr := connection.Azure.GetManagedClusterAdminCredentials(config.Vars.ServicePacks.AKS.ResourceGroupName, config.Vars.ServicePacks.AKS.ClusterName)

	if credsErr != nil {
		log.Printf("[DEBUG] Error trying to get cluster admin credentials: %v", err)
		//pass test if I got an error trying to obtain cluster admin credentials.
		// TODO: check the actual error message - risk of false positives if error was for another reason.
		err = nil
	}
	err = utils.ReformatError("I was able to get cluster admin credentials")

	payload = struct {
		Placeholder string
	}{
		Placeholder: "placeholder",
	}

	return err

}

// Name presents the name of this probe for external reference
func (probe probeStruct) Name() string {
	return "iam"
}

// Path presents the path of these feature files for external reference
func (probe probeStruct) Path() string {
	return probeengine.GetFeaturePath("internal", "azure", "kubernetes", probe.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (probe probeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		//setup AzureIdentity stuff ..??  Or should this be a pre-test setup
	})

	ctx.AfterSuite(func() {
		//tear down AzureIdentity stuff?
	})
}

// ScenarioInitialize initialises the specific test steps
func (probe probeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&scenario, probe.Name(), s)
	})

	// Background
	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, scenario.aKubernetesClusterIsDeployed)

	// Steps
	ctx.Step(`^an "([^"]*)" called "([^"]*)" exists in the namespace called "([^"]*)"$`, scenario.aResourceTypeXCalledYExistsInNamespaceCalledZ)
	ctx.Step(`^I succeed to create a simple pod in "([^"]*)" namespace assigned with the "([^"]*)" AzureIdentityBinding$`, scenario.iSucceedToCreateASimplePodInNamespaceAssignedWithThatAzureIdentityBinding)
	ctx.Step(`^an attempt to obtain an access token from that pod should "([^"]*)"$`, scenario.anAttemptToObtainAnAccessTokenFromThatPodShouldX)
	ctx.Step(`^I create an AzureIdentityBinding called "([^"]*)" in the Probr namespace bound to the "([^"]*)" AzureIdentity$`, scenario.iCreateAnAzureIdentityBindingCalledInANondefaultNamespace)
	ctx.Step(`^the cluster has managed identity components deployed$`, scenario.theClusterHasManagedIdentityComponentsDeployed)
	ctx.Step(`^the execution of a "([^"]*)" command inside the MIC pod is "([^"]*)"$`, scenario.theExecutionOfAXCommandInsideTheMICPodIsY)

	ctx.Step(`^no AAD user should have the Azure Kubernetes Service Cluster Admin Role role assigned to them for this cluster$`, scenario.checkClusterRBACForAdminRole)
	ctx.Step(`^I should not be able to obtain the cluster admin kubeconfig$`, scenario.checkCannotObtainClusterAdminCredentials)

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

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.Name = gs.Name
	s.probeAudit = summary.State.GetProbeLog(probeName)
	s.Audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.pods = make([]string, 0)
	s.namespace = config.Vars.ServicePacks.Kubernetes.ProbeNamespace
	s.azureIdentityBindings = make([]string, 0)
	probeengine.LogScenarioStart(gs)
}

func afterScenario(scenario scenarioState, probe probeStruct, gs *godog.Scenario, err error) {
	if config.Vars.ServicePacks.Kubernetes.KeepPods == "false" {
		for _, podName := range scenario.pods {
			err = connection.Kubernetes.DeletePodIfExists(podName, scenario.namespace, probe.Name())
			if err != nil {
				log.Printf(fmt.Sprintf("[ERROR] Could not retrieve pod from namespace '%s' for deletion: %s", scenario.namespace, err))
			}
		}
	}
	probeengine.LogScenarioEnd(gs)
}

func (scenario *scenarioState) createPodfromObject(podObject *apiv1.Pod) (createdPodObject *apiv1.Pod, err error) {
	createdPodObject, err = connection.Kubernetes.CreatePodFromObject(podObject, Probe.Name())
	if err == nil {
		scenario.pods = append(scenario.pods, createdPodObject.ObjectMeta.Name)
	}
	return
}

func azureIdentityExistsInNamespace(azureIdentityName, namespace string) (exists bool, resource k8sConnection.APIResource, err error) {

	resource, getError := connection.AKS.GetIdentityByNameAndNamespace(azureIdentityName, namespace)
	if getError != nil {
		if errors.IsStatusCode(404, getError) {
			exists = false
			return
		}
		err = utils.ReformatError("An error occured while retrieving Azure Identities from K8s cluster: %v", getError)
		return
	}

	exists = true
	return
}

func azureIdentityBindingExistsInNamespace(azureIdentityBindingName, namespace string) (exists bool, resource k8sConnection.APIResource, err error) {

	resource, getError := connection.AKS.GetIdentityBindingByNameAndNamespace(azureIdentityBindingName, namespace)
	if getError != nil {
		if errors.IsStatusCode(404, getError) {
			exists = false
			return
		}
		err = utils.ReformatError("An error occured while retrieving Azure Identity Bindings from K8s cluster: %v", getError)
		return
	}

	exists = true
	return
}

// azureCreateAIB creates an AzureIdentityBinding in the cluster
func azureCreateAIB(namespace, aibName, aiName string) (aibResource k8sConnection.APIResource, err error) {

	resource, createErr := connection.AKS.CreateAIB(namespace, aibName, aiName)
	if errors.IsStatusCode(409, createErr) { // Already Exists
		// TODO: Delete and recreate ?
		createErr = nil
	}

	err = createErr
	aibResource = resource

	return
}
