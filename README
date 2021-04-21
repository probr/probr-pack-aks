v0.1 uses OPA to assert configuration compliance against a json representation of AKS.

To use
1. Get the JSON representation of an AKS cluster: `az aks show -n ${cluster-name} -g ${resource-group} > aks.json`
1. Copy aks.json into internal/azure/aks
1. Run `make binary`
1. Run using probr-core

Todo:
1. Automate getting the Json representation via the Azure Rest API
1. Move the Azure identity tests from the Kubernetes pack into this service pack
