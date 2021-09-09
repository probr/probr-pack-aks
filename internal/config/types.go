package config

import kc "github.com/probr/probr-sdk/providers/kubernetes/config"

type varOptions struct {
	VarsFile     string
	Verbose      bool
	ServicePacks servicePacks `yaml:"ServicePacks"`
}

type servicePacks struct {
	AKS        aks `yaml:"AKS"`
	Kubernetes kc.Kubernetes
}

type aks struct {
	ClusterName       string    `yaml:"ClusterName"`
	ResourceGroupName string    `yaml:"ResourceGroupName"`
	ManagedID         managedID `yaml:"ManagedPodIdentity"`
	TagInclusions     []string  `yaml:"TagInclusions"`
	TagExclusions     []string  `yaml:"TagExclusions"`
}

type managedID struct {
	DefaultNamespaceAIB string `yaml:"DefaultAzureIdentityNamespace"`
	IdentityNamespace   string `yaml:"IdentityPodNamespace"`
}
