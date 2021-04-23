package aks

default kube_dashboard = false
kube_dashboard {
   m := input.properties.addonProfiles.kubeDashboard.enabled
   m == false
}

default enable_rbac = false
enable_rbac {
  m := input.properties.enableRBAC
  m == true
}

default azure_policy = false
azure_policy {
   m := input.properties.addonProfiles.azurepolicy.enabled
   m == true
}

default network_policy = false
network_policy {
  m := input.properties.networkProfile.networkPolicy
  m == "azure"
}
network_policy {
  m := input.properties.networkProfile.networkPolicy
  m == "calico"
}

default private_cluster = false
private_cluster {
  m := input.properties.apiServerAccessProfile.enablePrivateCluster
  m == true
}

default network_outbound_type = false
network_outbound_type {
  m := input.properties.networkProfile.outboundType
  m == "userDefinedRouting"
}

default disk_encryption = false
disk_encryption {
  m := input.properties.diskEncryptionSetId
  startswith(m, "/subscriptions")
}

default node_public_ip = false
node_public_ip {
  node_public_ips := {agentPool | agentPool := input.properties.agentPoolProfiles[_]; agentPool.enableNodePublicIp == true}
  count(node_public_ips) == 0
}
