Feature: Cluster networking best practice

  Background:
    Given an Azure Kubernetes cluster we can read the configuration of

    Scenario: Private Cluster is enabled
      Then Private Cluster is enabled

    Scenario: outbound network routing is user controlled
      Then outbound network routing is user controlled

    Scenario: Kubernetes network policy is enabled
      Then Kubernetes network policy is enabled

    Scenario: Nodes do not have Public IPs
      Then Kubernetes node hosts do not have public IPs
