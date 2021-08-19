# Probr AKS Service Pack
## Probes Provenance

The AKS Probr service pack has been built based on the [CIS GKE Benchmark 1.5.1](https://www.cisecurity.org/cis-benchmarks/), transposed into the context of AKS.

It compliments the Kubernetes Probr service pack, which covers additional controls not covered by the AKS-specific pack.

## Controls covered

| CIS ID | CIS Policy Statement | Probr Implementation | Suggested further improvements |
| ------ | ------               | -------------------- | ------------------- |
| 5.1.1  | Ensure that the cluster-admin role is only used where required | Check role assignments in Azure Control plane; fail if Azure Kubernetes Admin role is assigned. Try to get cluster admin credentials via Azure control plane (equivalent to `az aks get-credentials --admin`). Fail if Probr can get admin creds. |  [1] Add Check for clusterrolebindings to the cluster-admin clusterrole [2] inspect all roles assigned to the cluster role hierarchy for the cluster admin action |
| 5.3.1 | Ensure that the CNI in use supports Network Policies | OPA to check cluster configuration | - |
| 5.3.2 | Ensure that all Namespaces have Network Policies defined | To Do | To Do |
| 6.2 | Identity and Access Management (IAM) | This check is not part of the GKE benchmarks but is AKS-specific. We perform several tests to ensure Azure Pod Managed Identity is setup securely | It is a bit clunky to set up - look at making it simpler by discovering some of the config items |
| 6.5 | Nodes don't have public IP addresses | OPA check against configuration | Try to connect from Internet to the public endpoints (to see if mitigating controls have been put in place) |
| 6.6.4 | Ensure clusters are created with Private Endpoint Enabled and Public Access Disabled | OPA check against configuration | - |
| 6.6.5 | Ensure clusters are created with Private Nodes | OPA check against configuration to ensure agent pool nodes don't have public IP | - |
| 6.6.7 | Ensure Network Policy is Enabled and set as appropriate | OPA check against configuration | - |
| 6.8.1 | Ensure Basic Authentication using static passwords is Disabled | OPA check to ensure AAD RBAC integration is enabled | - |
| 6.8.2 | Ensure authentication using client certificates is Disabled | OPA check to ensure AAD RBAC integration is enabled | - |
| 6.9.1 | Enable Customer-Managed Encryption Keys for GKE Persistent Disks | Create a dynamic persistent volume claim from inside of AKS, then check the Disk to ensure it has customer managed keys | - |
| 6.10.1 | Ensure Kubernetes Web UI is Disabled | OPA check against configuration; look for kubernetes dashboard pod in kube-system namespace | - |
| 6.10.2 | Ensure that Alpha clusters are not used for production workloads | To Do | Check that preview features aren't being used in production?? |
| 6.10.3 | Ensure Pod Security Policy is Enabled and set as appropriate | OPA check that Azure Policy is enabled (plus deep PSP checks in Kubernetes Service Pack) | - |
