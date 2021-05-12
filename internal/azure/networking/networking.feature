@aks-net
@probes/aks/networking
Feature: Cluster networking best practice

  Background:
    Given an Azure Kubernetes cluster we can read the configuration of

    @aks-net-001
    Scenario: Private Cluster is enabled
      Then Private Cluster is enabled

    @aks-net-002
    Scenario: outbound network routing is user controlled
      Then outbound network routing is user controlled

    @aks-net-003
    Scenario: Kubernetes network policy is enabled
      Then Kubernetes network policy is enabled

    @aks-net-004
    Scenario: Nodes do not have Public IPs
      Then Kubernetes node hosts do not have public IPs
