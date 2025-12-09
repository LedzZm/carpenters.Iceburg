package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/goccy/go-yaml"
)

type Package struct {
	Dependencies []string `json:"Depends"`
	Description  string   `json:"Description"`
}

type AURResponse struct {
	ResultCount int       `json:"resultcount"`
	Packages    []Package `json:"results"`
	Type        string
	apiVersion  string `json:"version"`
}

type Config struct {
	Type string
}

// TODO: Error handling in all _ .
// TODO: URL Constants.
func main() {

	// TODO: Discover path
	// TODO: If Discovery is not too easy, Add parametricaly or find another way?
	configFilePath := "./packages.yml"
	configData, _ := os.ReadFile(configFilePath)
	// TODO: DEFER CLOSE see todowner
	// Unmarshal YAML into a map where each key (a, b, etc.) points to a Config
	var configMap map[string]Config
	yaml.Unmarshal(configData, &configMap)

	// Print the unmarshalled data
	for key, cfg := range configMap {
		fmt.Printf("Key: %s\n", key)
		fmt.Printf("Config: %+v\n", cfg)
	}

	os.Exit(1)
	packageName := "discord-canary"
	response, _ := http.Get("https://aur.archlinux.org/rpc/v5/info/" + packageName)

	fmt.Println("https://aur.archlinux.org/rpc/v5/info/" + packageName)

	defer response.Body.Close()
	bodyBytes, _ := io.ReadAll(response.Body)

	aurResponse := AURResponse{}
	json.Unmarshal(bodyBytes, &aurResponse)

	dependencies := aurResponse.Packages[0].Dependencies

	// TODO: Sort by dependencies.
	// TODO: Categorize by dependencies and do concurrently (?).
	for _, element := range dependencies {
		fmt.Println(element)
	}

	fmt.Println()

	// fmt.Println(string(bodyBytes))
}
