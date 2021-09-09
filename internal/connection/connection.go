package connection

import (
	"context"
	"log"

	"github.com/probr/probr-pack-aks/internal/config"
	azureutil "github.com/probr/probr-sdk/providers/azure"
	"github.com/probr/probr-sdk/providers/azure/aks"
	azureconnection "github.com/probr/probr-sdk/providers/azure/connection"
	"github.com/probr/probr-sdk/providers/kubernetes/connection"
)

var Kubernetes *connection.Conn
var Azure *azureconnection.AzureConnection
var AKS *aks.AKS
var errStrings []string

func Init() {
	subID, tenantID, clientID, secret, good := validateConfig()
	if !good {
		log.Printf("[ERROR] Missing required config vars")
		return // Do not attempt connection without required values
	}
	Kubernetes = connection.NewConnection(
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath,
		config.Vars.ServicePacks.Kubernetes.KubeContext,
		config.Vars.ServicePacks.Kubernetes.ProbeNamespace,
	)
	AKS = aks.NewAKS(Kubernetes)
	Azure = azureconnection.NewAzureConnection(
		context.Background(),
		subID,
		tenantID,
		clientID,
		secret,
	)
	for i, conn := range []interface{}{Kubernetes, AKS, Azure} {
		if conn == nil {
			log.Printf("[ERROR] Failed initialization (%v)", i)
		}
	}
}

func validateConfig() (subID, tenantID, clientID, secret string, good bool) {
	subID, subErr := azureutil.SubscriptionID()
	tenantID, tenErr := azureutil.TenantID()
	clientID, cliErr := azureutil.ClientID()
	secret, secErr := azureutil.ClientSecret()

	for _, e := range []error{subErr, tenErr, cliErr, secErr} {
		if e != nil {
			errStrings = append(errStrings, e.Error())
		}
	}
	if errStrings != nil {
		return
	}
	good = true
	return
}

// Error reformats a collection of any errors encountered during connection
func Errors() []string {
	return errStrings
}
