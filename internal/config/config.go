package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	sdkConfig "github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/config/setter"
	"github.com/citihub/probr-sdk/utils"
)

// Vars is a stateful object containing the variables required to execute this pack
var Vars varOptions

// Init will set values with the content retrieved from a filepath, env vars, or defaults
func (ctx *varOptions) Init() (err error) {
	if ctx.VarsFile != "" {
		ctx.decode()
		if err != nil {
			return
		}
	} else {
		log.Printf("[DEBUG] No vars file provided, unexpected behavior may occur")
	}
	sdkConfig.GlobalConfig.VarsFile = ctx.VarsFile
	sdkConfig.GlobalConfig.Init()
	sdkConfig.GlobalConfig.CloudProviders.Azure.SetEnvAndDefaults()

	ctx.ServicePacks.Kubernetes.setEnvAndDefaults() //TODO: !! CHANGE ME !!
	ctx.ServicePacks.AKS.setEnvAndDefaults()

	log.Printf("[DEBUG] Config initialized by %s", utils.CallerName(1))
	return
}

// decode uses an SDK helper to create a YAML file decoder,
// parse the file to an object, then extracts the values from
// ServicePacks.Kubernetes into this context
func (ctx *varOptions) decode() (err error) {
	configDecoder, file, err := sdkConfig.NewConfigDecoder(ctx.VarsFile)
	if err != nil {
		return
	}
	err = configDecoder.Decode(&ctx)
	file.Close()
	return err
}

// LogConfigState will write the config file to the write directory
func (ctx *varOptions) LogConfigState() {
	json, _ := json.MarshalIndent(ctx, "", "  ")
	log.Printf("[INFO] Config State: %s", json)
	// path := filepath.Join("config.json")
	// if ctx.WriteConfig == "true" && utils.WriteAllowed(path) {
	// 	data := []byte(json)
	// 	ioutil.WriteFile(path, data, 0644)
	// 	//log.Printf("[NOTICE] Config State written to file %s", path)
	// }
}

func (ctx *varOptions) Tags() string {
	return sdkConfig.ParseTags(ctx.ServicePacks.AKS.TagInclusions, ctx.ServicePacks.AKS.TagExclusions)
}

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func (ctx *kubernetes) setEnvAndDefaults() {
	// Notes on SetVar's values:
	// 1. Pointer to local object; will be overwritten by env or default if empty
	// 2. Name of env var to check
	// 3. Default value to set if flags, vars file, and env have not provided a value

	setter.SetVar(&ctx.KeepPods, "PROBR_KEEP_PODS", "false")
	setter.SetVar(&ctx.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	setter.SetVar(&ctx.KubeContext, "KUBE_CONTEXT", "")
	setter.SetVar(&ctx.AuthorisedContainerImage, "PROBR_AUTHORISED_IMAGE", "")
}

func (ctx *aks) setEnvAndDefaults() {
	setter.SetVar(&ctx.ClusterName, "PROBR_AKS_CLUSTER_NAME", "")
	setter.SetVar(&ctx.ResourceGroupName, "PROBR_AKS_RG_NAME", "")
	setter.SetVar(&ctx.ManagedID.DefaultNamespaceAIB, "PROBR_AKS_AIB_NS", "default")
	setter.SetVar(&ctx.ManagedID.IdentityNamespace, "PROBR_AKS_ID_POD_NS", "kube-system")
}

func getDefaultKubeConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}
