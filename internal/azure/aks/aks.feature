Feature: General cluster Security best practices

  Background:
    Given an Azure Kubernetes cluster we can read the configuration of

    Scenario:
       Then the Kubernetes Web UI is disabled

    Scenario:
       Then Azure Policy is enabled

    Scenario:
       Then Azure AD integration is enabled

    Scenario:
      Then Private Cluster is enabled

    Scenario:
      Then Disk Encryption is enabled

    Scenario:
      Then outbound network routing is user controlled

    Scenario:
      Then CNI network policy is enabled

    Scenario:
      Then Kubernetes node hosts do not have public IPs
