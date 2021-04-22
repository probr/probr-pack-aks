package aksencryptionatrest

import (
	"fmt"
	//	"log"

	"github.com/cucumber/godog"
	//	apiv1 "k8s.io/api/core/v1"

	//	"github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/config"
	//	"github.com/citihub/probr-sdk/probeengine"
	//	"github.com/citihub/probr-sdk/providers/kubernetes/connection"
	"github.com/citihub/probr-sdk/providers/kubernetes/constructors"
	//	"github.com/citihub/probr-sdk/providers/kubernetes/errors"
	"github.com/citihub/probr-sdk/utils"
)

func (scenario *scenarioState) aKubernetesClusterIsDeployed() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()
	stepTrace.WriteString(fmt.Sprintf("Validate that a cluster can be reached using the specified kube config and context; "))

	payload = struct {
		KubeConfigPath string
		KubeContext    string
	}{
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath,
		config.Vars.ServicePacks.Kubernetes.KubeContext,
	}

	err = kConnection.ClusterIsDeployed() // Must be assigned to 'err' be audited
	return err
}

func (scenario *scenarioState) iCreateAPodWhichDynamicallyCreatesAnAzureDisk() error {

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Build a pod spec with default values; ")
	podObject := constructors.PodSpec(Probe.Name(), scenario.namespace)

	//TODO - make storage class configurable
	pvcObject := constructors.DynamicPersistentVolumeClaim(Probe.Name(), scenario.namespace, "default")

	/*
		 spec.containers.volumeMounts
		  [] mountPath: "/blah"
			   name: volumename
		 spec.volumes
		 	[] name: volumename
			   persistentVolumeClaim:
				   claimName: <pvc name>
	*/
	constructors.AddPVCToPod(podObject, pvcObject)

	stepTrace.WriteString("Create pod from spec; ")
	createdPVCObject, pvcCreationErr := kConnection.CreatePVCFromObject(pvcObject, Probe.Name())
	createdPodObject, podCreationErr := kConnection.CreatePodFromObject(podObject, Probe.Name()) // Pod name is saved to scenario state if successful
	if podCreationErr != nil {
		return podCreationErr
	}

	if pvcCreationErr != nil {
		return pvcCreationErr
	}

	scenario.pods = append(scenario.pods, createdPodObject.ObjectMeta.Name)
	scenario.pvcs = append(scenario.pvcs, createdPVCObject.ObjectMeta.Name)
	return nil
}

func (scenario *scenarioState) theDiskIsEncryptedUsingCustomerManagedKeys() error {
	return godog.ErrPending
}
