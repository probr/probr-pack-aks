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
