package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

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
	Type    string
	Package string
}

// TODO: Error handling in all _ .
// TODO: URL Constants.
func main() {

	// TODO: Discover path
	// TODO: If Discovery is not too easy, Add parametricaly or find another way?
	configFilePath := "./packages.yml"
	configData, _ := os.ReadFile(configFilePath)
	// TODO: DEFER CLOSE see todowner

	// TODO: Does this need to be a map?
	var configMap map[string]Config
	yaml.Unmarshal(configData, &configMap)

	// TODO: Handle AUR?
	// TODO: Handle custom?
	pacmanPackages := make(map[string]Config)
	var pacmanKeys []string
	// Split packages by type
	for packageKey, cfg := range configMap {
		switch cfg.Type {
		case "pacman":
			// If a package key is specified use it.
			// Otherwise default to the yaml key.
			_packageKey := packageKey
			if cfg.Package != "" {
				_packageKey = cfg.Package
			}

			// TODO: This is a temporary skip for dev speed, remove this.
			if cfg.Package == "cachyos-gaming-meta" {
				continue
			}

			pacmanKeys = append(pacmanKeys, _packageKey)
			// TODO: Is this even needed??
			pacmanPackages[packageKey] = cfg

		default:
			fmt.Printf("Unknown type for package %s: %s\n", packageKey, cfg.Type)
		}

	}

	// TODO: Stream buffer while executing.
	command := exec.Command("sudo", append([]string{"pacman", "-S", "--noconfirm"}, pacmanKeys...)...)
	out, err := command.CombinedOutput()
	fmt.Println(string(out), err)

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
