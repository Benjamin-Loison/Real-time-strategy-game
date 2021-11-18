package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"flag"
)

type Server_t struct {
	Hostname string `json:"Hostname"`
	Port     int `json:"Port"`
}
type Keys_t struct {
    Left string `json:"Left"`
    Right string `json:"Right"`
    Up string `json:"Up"`
    Down string `json:"Down"`
}

type Configuration_t struct {
    Server Server_t `json:"Server"`
    Keys Keys_t `json:Keys"`
    Pseudo string `json:Pseudo"`
}
var client_chan chan string

func loadConfig(file_name string) Configuration_t {
	// Read the main configuration file
	file, err := ioutil.ReadFile(file_name)
	if err != nil {
		logging("Configuration", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}

	// Convert the json content into a Configuration structure
	var configuration = &Configuration_t{}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		logging("Configuration", fmt.Sprintf("Cannot parse the config file: %v", err))
		os.Exit(1)
	}
	// Parse the command line and overwrite the configuration if needed
	override_addr_parsed := flag.String("n", configuration.Server.Hostname, "Hostname of the server (ip address or name)")
	override_port_parsed := flag.Int("p", configuration.Server.Port, "Port of the remote server")
	flag.Parse()
	// Replace the configuration variables
	configuration.Server.Hostname = *override_addr_parsed
	configuration.Server.Port = *override_port_parsed

	logging("Configuration",
		fmt.Sprintf("Hostname: %s, port: %d", configuration.Server.Hostname, configuration.Server.Port))

	return *configuration
}
