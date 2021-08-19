# Probr AKS Service Pack

The [Probr](https://github.com/citihub/probr-core) AKS Service pack compliments the [Kubernetes service pack](https://github.com/citihub/probr-pack-kubernetes) with Azure Kubernetes Service (AKS)-specific compliance checks.

Click [here](./COVERAGE.md) to see the current state of the probes in this pack.

## To Build

The following will build a binary named "aks":
```
git clone https://github.com/citihub/probr-pack-aks.git
cd probr-pack-aks
make binary
```

Move the `aks` binary into your probr service pack location (default is `${HOME}/probr/binaries`)

## Pre-Requisites

You will need
1. Probr Core (https://github.com/citihub/probr-core)
1. An AKS Cluster
1. An active kubeconfig against the cluster, that can deploy into the probe namespace (see config below. Default is probr-general-test-ns)
1. A service principle that has the "Reader" Azure role on the cluster
1. For the IAM probes, you will need:
  - Managed Pod Identity to be configured in your cluster
  - An `AzureIdentity` called `probr-ai` and an `AzureIdentityBinding` called `probr-aib` to be deployed in the default namespace (or the namespace configured in the runtime config, see below).

## Configuration

### Minimum configuration

The minimum required additions to your Probr runtime configuration is as follows:

```
Run:
  - "aks"
ServicePacks:
  Kubernetes:
    AuthorisedContainerImage: "yourprivateregistry.io/citihub/probr-probe"
  AKS:
    ClusterName: "your-clustername-here"
    ResourceGroupName: "your-resource-group-name-here"      
    CloudProviders:
      Azure:
        TenantID: "UUID of your tenant"
        SubscriptionID: "UUID of your subscription"
        ClientID: "Client ID UUID of your service principle"
        ClientSecret: "Recommend leaving this blank and using envvar"
```
We recommend _not_ storing the ClientSecret in the config.yml, instead use the `PROBR_AZURE_CLIENT_SECRET` environment variable.

### Full configuration

If you don't want to use the defaults you can add the following to your Probr config.yml:

```
Run:
  - "aks"
ServicePacks:
  Kubernetes:
    KubeConfig: "location of your kubeconfig if not the default"
    KubeContext: "specific kubecontext if not the current context"
    AuthorisedContainerImage: "yourprivateregistry.io/citihub/probr-probe"
    ProbeNamespace: "namespace Probr deploys into. Defaults to 'probr-general-test-ns'"
  AKS:
    ClusterName: "your-clustername-here"
    ResourceGroupName: "your-resource-group-name-here"
    ManagedPodIdentity:
      DefaultAzureIdentityNamespace: "Namespace where the probr-ai and probr-aib live for cross-namespace identity tests. Defaults to 'default'"
      IdentityPodNamespace: "namespace where the MIC and NMI pods live. Defaults to 'kube-system'"
CloudProviders:
  Azure:
    TenantID: "UUID of your tenant"
    SubscriptionID: "UUID of your subscription"
    ClientID: "Client ID UUID of your service principle"
    ClientSecret: "Recommend leaving this blank and using envvar"
```



## Running the Service Pack

If all of the instructions above have been followed, then you should be able to run `./probr` and the service pack will run.
