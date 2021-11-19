package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"flag"
	"strings"
)

type Server_t struct {
	Hostname string `json:"Hostname"`
	Port     int `json:"Port"`
}
type Keys_t struct {
    Left int32
    Right int32
    Up int32
    Down int32
}
// Keys are represented by strings in the configuration file. This will allows
// us to use words representing non-ascii keys.
type Keys_tmp_t struct {
    Left string `json:"Left"`
    Right string `json:"Right"`
    Up string `json:"Up"`
    Down string `json:"Down"`
}

// Represents json structure
type Configuration_tmp_t struct {
    Server Server_t `json:"Server"`
    Keys Keys_tmp_t `json:Keys"`
    Pseudo string `json:Pseudo"`
}

// The actual config file: the keys are replaced by their raylib values
type Configuration_t struct {
    Server Server_t `json:"Server"`
    Keys Keys_t `json:Keys"`
    Pseudo string `json:Pseudo"`
}

func loadConfig(file_name string) Configuration_t {
	// Read the main configuration file
	file, err := ioutil.ReadFile(file_name)
	if err != nil {
		logging("Configuration", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}

	// Convert the json content into a Configuration structure
	var configuration = &Configuration_t{}
	var configuration_tmp = &Configuration_tmp_t{}
	err = json.Unmarshal(file, &configuration_tmp)
	if err != nil {
		logging("Configuration", fmt.Sprintf("Cannot parse the config file: %v", err))
		os.Exit(1)
	}
	// Fills the Configration structure usinf the temp structure
	configuration.Server = configuration_tmp.Server
	configuration.Keys.Left = keyOfString(configuration_tmp.Keys.Left)
	configuration.Keys.Right = keyOfString(configuration_tmp.Keys.Right)
	configuration.Keys.Up = keyOfString(configuration_tmp.Keys.Up)
	configuration.Keys.Down = keyOfString(configuration_tmp.Keys.Down)
	configuration.Pseudo = configuration_tmp.Pseudo

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

//Raylib represents keys as int32 values
func keyOfString(s string)(int32) {
	switch l := len(s); l {
		case 0:
			panic("An empty string cannot represent a key in the configuration file.")
		case 1:
			return int32([]rune(strings.ToUpper(s))[0])
		default:
			panic("Not implemented: recognition of non-ascii characters and description")
	}
}

