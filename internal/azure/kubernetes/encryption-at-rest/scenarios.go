package aksencryptionatrest

import (
	"fmt"
	"log"
	"time"

	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/providers/kubernetes/constructors"
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

	//TODO: make storage class configurable. Hardcode to 'default' at the moment
	pvcObject := constructors.DynamicPersistentVolumeClaim(Probe.Name(), scenario.namespace, "default")

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
	pvc, err := kConnection.GetPVCFromPVCName(scenario.pvcs[0], scenario.namespace)
	if err != nil {
		log.Printf("[DEBUG] Error getting PVC from PVC Name")
		return err
	}

	/*
		the volumename isn't available until the PVC changes into "Bound" status, which takes a few seconds
		so wait for the status to change
	*/
	for pvc.Status.Phase == "Pending" {
		log.Printf("[DEBUG] PVC Status.Phase: %s; Waiting...", pvc.Status.Phase)

		time.Sleep(2 * time.Second)
		pvc, _ = kConnection.GetPVCFromPVCName(scenario.pvcs[0], scenario.namespace)
	}

	log.Printf("[DEBUG] PVC name is %s. PV name is %s.", scenario.pvcs[0], pvc.Spec.VolumeName)

	pv, err := kConnection.GetPVFromPVName(pvc.Spec.VolumeName)
	if err != nil {
		log.Printf("[DEBUG] Error getting PV from PV Name")
		log.Printf("[DEBUG] PVC trace: %v", pvc)
		return err
	}

	log.Printf("[DEBUG] Disk URI is %s", pv.Spec.AzureDisk.DataDiskURI)

	rgName, diskName := azConnection.ParseDiskDetails(pv.Spec.AzureDisk.DataDiskURI)
	log.Printf("[DEBUG] Disk details are rgName: %s. diskName: %s", rgName, diskName)

	azureDisk, err := azConnection.GetDisk(rgName, diskName)
	if err != nil {
		log.Printf("Error getting disk client")
		return err
	}

	encryptionType := azureDisk.DiskProperties.Encryption.Type
	log.Printf("[DEBUG] Disk Encryption Type: %s", encryptionType)
	if encryptionType == "EncryptionTypeEncryptionAtRestWithCustomerKey" {
		return nil
	}

	return fmt.Errorf("Disk %s in resource group %s not encrypted with customer key. Encryption type was %s", diskName, rgName, encryptionType)

}
