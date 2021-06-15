module github.com/probr/probr-pack-aks

go 1.14

require (
	github.com/cucumber/godog v0.11.0
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/markbates/pkger v0.17.1
	github.com/probr/probr-sdk v0.0.0-00010101000000-000000000000
	// github.com/probr/probr-sdk v0.1.5
	k8s.io/api v0.19.6
)

// For Development Only
replace github.com/probr/probr-sdk => ../probr-sdk
