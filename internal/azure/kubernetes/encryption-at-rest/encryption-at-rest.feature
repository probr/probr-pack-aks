@aks-ear
@probes/aks/encrypion-at-rest
Feature: Ensure data is encrypted using customer managed keys
    As a Security Auditor
    I want to ensure that stringent authentication and authorisation policies are applied to my organisation's Kubernetes clusters
    So that only approved actors have the ability to perform sensitive operations in order to prevent malicious attacks on my organization

    @aks-ear-001
    Scenario: Azure Managed Disks created by AKS are encrypted using customer managed keys
        Given a Kubernetes cluster exists which we can deploy into
        When I create a Pod which dynamically creates an Azure Disk
        Then the disk is encrypted using Customer Managed Keys

    @aks-ear-002
    Scenario: Disk Encryption is enabled
        Given an Azure Kubernetes cluster we can read the configuration of
        Then Disk Encryption is enabled
