@aks-ear
@probes/aks/encrypion-at-rest
Feature: Ensure stringent authentication and authorisation
    As a Security Auditor
    I want to ensure that stringent authentication and authorisation policies are applied to my organisation's Kubernetes clusters
    So that only approved actors have the ability to perform sensitive operations in order to prevent malicious attacks on my organization

    Background:
        Given a Kubernetes cluster exists which we can deploy into

    @aks-ear-001
    Scenario Outline: Azure Managed Disks created by AKS are encrypted using customer managed keys

        When I create a Pod which dynamically creates an Azure Disk
        Then the disk is encrypted using Customer Managed Keys
