package config

type varOptions struct {
	VarsFile     string
	Verbose      bool
	ServicePacks servicePacks `yaml:"ServicePacks"`
}

type servicePacks struct {
	Kubernetes kubernetes `yaml:"Kubernetes"`
	AKS        aks        `yaml:"AKS"`
}

type kubernetes struct {
	KeepPods                 string `yaml:"KeepPods"` // TODO: Change type to bool, this would allow us to remove logic from kubernetes.GetKeepPodsFromConfig()
	KubeConfigPath           string `yaml:"KubeConfig"`
	KubeContext              string `yaml:"KubeContext"`
	AuthorisedContainerImage string `yaml:"AuthorisedContainerImage"`
	ProbeNamespace           string `yaml:"ProbeNamespace"`
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
