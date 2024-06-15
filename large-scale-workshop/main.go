package main

import (
	"log"
	"os"

	"github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/config"
	RegistryService "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/service"
	TestService "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/service"
	"github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/utils"
	"gopkg.in/yaml.v2"
)

func main() {
	// read configuration file from command line argument
	if len(os.Args) != 2 {
		utils.Logger.Fatal("Expecting exactly one configuration file")
		os.Exit(1)
	}
	configFile := os.Args[1]
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
		os.Exit(2)
	}

	var config config.ConfigBase
	err = yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		log.Fatalf("error unmarshaling data: %v", err)
		os.Exit(3)
	}

	switch config.Type {
	case "TestService":
		utils.Logger.Printf("Loading service type: %v\n", config.Type)
		TestService.Start(configData)
	case "RegistryService":
		utils.Logger.Printf("Loading service type: %v\n", config.Type)
		RegistryService.RunServer(configData)
	default:
		utils.Logger.Fatalf("Unknown configuration type: %v", config.Type)
		os.Exit(4)
	}
}
