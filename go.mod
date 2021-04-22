module github.com/citihub/probr-pack-aks

go 1.14

require (
	github.com/Azure/aad-pod-identity v1.7.5 // indirect
	github.com/citihub/probr-sdk v0.0.18
	github.com/cucumber/godog v0.11.0
	github.com/markbates/pkger v0.17.1
	k8s.io/api v0.19.6
)

replace github.com/citihub/probr-sdk => ../probr-sdk
