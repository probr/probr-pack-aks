package aks

default kube_dashboard = false
kube_dashboard {
   m := input.addonProfiles.kubeDashboard.enabled
   m == false
}

default enable_rbac = false
enable_rbac {
  m := input.enableRbac
  m == true
}

default azure_policy = false
azure_policy {
   m := input.addonProfiles.azurepolicy.enabled
   m == true
}
