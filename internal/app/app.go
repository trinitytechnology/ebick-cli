package app

import (
	"fmt"
	"os"

	"github.com/trinitytechnology/ebrick-cli/pkg/utils"
)

type AppConfig struct {
	Name          string   `yaml:"name"`
	Package       string   `yaml:"package"`
	Modules       []string `yaml:"modules"`
	Database      bool     `yaml:"database"`
	Observability bool     `yaml:"observability"`
	Cache         bool     `yaml:"cache"`
	Messaging     bool     `yaml:"messaging"`
	Version       string   `yaml:"version"`
}

var EBrickVersion = "v0.3.2"

const appManifest = ".ebrick.yaml"

// NewApp creates a new eBrick application
func NewApp() {

	var appConfig AppConfig

	// Check .ebrick.yaml file exists
	if !utils.FileExists(appManifest) {
		appConfig = NewApplicationCommandPrompts()
		err := utils.WriteYamlFile(appManifest, appConfig)
		if err != nil {
			os.Exit(1)
		}
	} else {
		overwrite := utils.GetYesOrNoInput("Overwrite existing configuration?", true)
		if !overwrite {
			return
		}
	}

	// Read .ebrick.yaml
	appConfig, err := utils.ReadYamlFile[AppConfig](appManifest)
	if err != nil {
		fmt.Println("Error reading .ebrick.yaml:", err)
		return
	}

	fmt.Println("Creating a new eBrick application with the name:", appConfig.Name)
	GenerateApplication(appConfig)

	fmt.Println("Application created successfully.")

	// Execute post generation tasks
	PostGenerated()
}

func NewApplicationCommandPrompts() AppConfig {

	appName := utils.GetUserInput("Enter the name of the application: ", true, "Application name is required.")
	packageName := utils.GetUserInput("Enter the application package: ", true, "Package name is required.")
	modulesInput := utils.GetUserInput("Enter the application modules (comma-separated, no spaces): ", false, "")
	modules := utils.ProcessSlicesInput(modulesInput)

	database := utils.GetYesOrNoInput("Do you need a database?", true)
	cache := utils.GetYesOrNoInput("Do you need a cache?", false)
	messaging := utils.GetYesOrNoInput("Do you need messaging?", false)
	observability := utils.GetYesOrNoInput("Do you need observability?", false)

	appConfig := AppConfig{
		Name:          appName,
		Package:       packageName,
		Modules:       modules,
		Database:      database,
		Observability: observability,
		Cache:         cache,
		Messaging:     messaging,
		Version:       EBrickVersion,
	}

	return appConfig
}

func RunApp() {
	// Run go mod tidy
	utils.ExecCommand("go", "mod", "tidy")

	// Run go mod tidy
	utils.ExecCommand("go", "run", "cmd/main.go")
}

func PostGenerated() {

	fmt.Println("Running post generation tasks...")

	// Run go mod tidy
	utils.ExecCommand("go", "mod", "tidy")

}
