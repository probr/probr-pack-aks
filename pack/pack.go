package pack

import (
	"github.com/markbates/pkger"
	azureaks "github.com/probr/probr-pack-aks/internal/azure/aks"
	kubeear "github.com/probr/probr-pack-aks/internal/azure/kubernetes/encryption-at-rest"
	kubeiam "github.com/probr/probr-pack-aks/internal/azure/kubernetes/iam"
	azurenw "github.com/probr/probr-pack-aks/internal/azure/networking"
	"github.com/probr/probr-sdk/probeengine"
)

// GetProbes returns a list of probe objects
func GetProbes() []probeengine.Probe {
	return []probeengine.Probe{
		azureaks.Probe,
		kubeiam.Probe,
		kubeear.Probe,
		azurenw.Probe,
	}
}

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/internal/azure/aks/aks.feature")
	pkger.Include("/internal/azure/networking/networking.feature")
	pkger.Include("/internal/common/aks.rego")
	pkger.Include("/internal/azure/kubernetes/iam/iam.feature")
	pkger.Include("/internal/azure/kubernetes/encryption-at-rest/encryption-at-rest.feature")
}
