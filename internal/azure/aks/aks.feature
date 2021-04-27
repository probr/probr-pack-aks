Feature: General cluster Security best practices

  Background:
    Given an Azure Kubernetes cluster we can read the configuration of

    Scenario: Kubernetes Web UI is disabled
       Then the Kubernetes Web UI is disabled

    Scenario: Azure Policy is enabled
       Then Azure Policy is enabled

    Scenario: Azure AD integration is enabled
       Then Azure AD integration is enabled
