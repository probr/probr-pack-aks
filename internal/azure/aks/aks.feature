Feature: General cluster Security best practices

  Background:
    Given an Azure Kubernetes cluster we can read the configuration of

    Scenario:
       Then the Kubernetes Web UI is disabled

    Scenario:
       Then Azure Policy is enabled

    Scenario:
       Then Azure AD integration is enabled
