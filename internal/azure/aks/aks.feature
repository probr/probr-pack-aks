@aks-gen
@probes/aks/general
Feature: General cluster Security best practices

  Background:
    Given an Azure Kubernetes cluster we can read the configuration of

    @aks-gen-001
    Scenario: Kubernetes Web UI is disabled
       Then the Kubernetes Web UI is disabled

    @aks-gen-002
    Scenario: Azure Policy is enabled
       Then Azure Policy is enabled

    @aks-gen-003
    Scenario: Azure AD integration is enabled
       Then Azure AD integration is enabled
