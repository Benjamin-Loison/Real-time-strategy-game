package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"flag"
)
type Configuration struct {
	hostname string `json:"hostname"`
	port     int `json:"port"`
}
var client_chan chan string

func loadConfig(file_name string) Configuration {
	// Read the main configuration file
	file, err := ioutil.ReadFile(file_name)
	if err != nil {
		logging("Configuration", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}

	// Convert the json content into a Configuration structure
	var configuration Configuration
	var conf map[string]interface{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		logging("Configuration", fmt.Sprintf("Cannot parse the config file: %v", err))
		os.Exit(1)
	} else {
		configuration = Configuration {
			hostname: conf["hostname"].(string),
			port: int(conf["port"].(float64))}
	}

	// Parse the command line and overwrite the configuration if needed
	override_addr_parsed := flag.String("n", configuration.hostname, "Hostname of the server (ip address or name)")
	override_port_parsed := flag.Int("p", configuration.port, "Port of the remote server")
	flag.Parse()
	// Replace the configuration variables
	configuration.hostname = *override_addr_parsed
	configuration.port = *override_port_parsed

	logging("Configuration",
		fmt.Sprintf("Hostname: %s, port: %d", configuration.hostname, configuration.port))

	return configuration
}
