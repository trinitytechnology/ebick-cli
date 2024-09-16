package app

import (
	"fmt"
	"os"

	ebrickcli "github.com/trinitytechnology/ebrick-cli"
	"github.com/trinitytechnology/ebrick-cli/internal/constants"
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

// NewApp creates a new eBrick application
func NewApp() {

	var appConfig AppConfig

	// Check .ebrick.yaml file exists
	if !utils.FileExists(constants.AppManifestFile) {
		appConfig = NewApplicationCommandPrompts()
		err := utils.WriteYamlFile(constants.AppManifestFile, appConfig)
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
	appConfig, err := utils.ReadYamlFile[AppConfig](constants.AppManifestFile)
	if err != nil {
		fmt.Println("Error reading .ebrick.yaml:", err)
		return
	}

	fmt.Println("Creating a new eBrick application with the name:", appConfig.Name)
	generator := NewAppGenerator(&appConfig)
	generator.Generate(appConfig)

	fmt.Println("Application created successfully.")

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
		Version:       ebrickcli.FrameworkVersion,
	}

	return appConfig
}

func RunApp() {
	// Run go mod tidy
	utils.ExecCommand("go", "mod", "tidy")

	// Run go mod tidy
	utils.ExecCommand("go", "run", "cmd/main.go")
}
