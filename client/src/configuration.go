package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"flag"
)

// This file loads the configuration of the application. This configuration
// contains two parts:
// ~ The client configuration, with defaults server parameters, nickname,
//   hotkeys, ...
// ~ The game configuration (menus, actions, ...)



/*
********************************************************************************
---------------------------SERVER~CONFIGURATION~BELOW---------------------------
********************************************************************************
*/

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



/*
********************************************************************************
----------------------------MENUS~CONFIGURATION~BELOW---------------------------
********************************************************************************
*/

type MenuElement_tmp_t struct {
	Name string `json:"Name"`
	Type string `json:"Type"`
	Key string `json:"Key"`
	Ref string `json:"Ref"`
}
type Menu_tmp_t struct {
	Ref string `json:"Ref"`
	Title string `json:"Title"`
	Elements []MenuElement_tmp_t `json:"Elements"`
}
type Action_tmp_t struct {
	Type string `json:"Type"`
	Ref string `json:"Ref"`
	Title string `json:"Title"`
}
type MenuConfiguration_tmp_t struct {
	Menus []Menu_tmp_t `json:"Menus"`
	Actions []Action_tmp_t `json:"Actions"`
}

type MenuElementType int32
type ActionType int32
const(
	MenuElementAction = 0
	MenuElementSubMenu = 1
	ActionBuilding = 0
	)
type MenuElement_t struct {
	Name string
	Type MenuElementType
	Key int32
	Ref string
}
type Menu_t struct {
	Ref string
	Title string
	Elements []MenuElement_t
}
type Action_t struct {
	Type ActionType
	Ref string
	Title string
}
type MenuConfiguration_t struct {
	Menus []Menu_t
	Actions []Action_t
}

// TODO: Translate the references to sub-menus and actions into their index in
// the final arrays
// TODO bis: Check that there is no conflict between hotkeys (at least within
// the menus configuration, or with the whole game interface(should be done
// later, within the main function(that loads the different configuration
// files))
func loadTextMenus(file_name string) MenuConfiguration_t {
	// Read the menus configuration file
	file, err := ioutil.ReadFile(file_name)
	if err != nil {
		logging("Menu configuration", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}

	// Convert the json content into a Configuration structure
	var configuration_tmp = &MenuConfiguration_tmp_t{}
	var configuration = &MenuConfiguration_t{}
	err = json.Unmarshal(file, &configuration_tmp)
	if err != nil {
		logging("Menu configuration", fmt.Sprintf("Cannot parse the config file: %v", err))
		os.Exit(1)
	}

	// Verbose and build the clean (well-typed) menus configuration
	configuration.Menus = nil
	configuration.Actions = nil
	for i := 0; i < len(configuration_tmp.Menus); i++ {
		// Declaration of local variables to store the current menu to add
		var active_menu Menu_t
		var inner_elements = []MenuElement_t(nil)

		active_menu.Title = configuration_tmp.Menus[i].Title
		active_menu.Ref = configuration_tmp.Menus[i].Ref

		for j := 0 ; j < len(configuration_tmp.Menus[i].Elements); j ++ {
			inner_elements = append(inner_elements,
				MenuElement_t {
					configuration_tmp.Menus[i].Elements[j].Name,
					MenuElementTypeIfString(
						configuration_tmp.Menus[i].Elements[j].Type),
					keyOfString(configuration_tmp.Menus[i].Elements[j].Key),
					configuration_tmp.Menus[i].Elements[j].Ref })
			active_menu.Elements= inner_elements
		}
		configuration.Menus = append(configuration.Menus, active_menu)
	}
	logging("Menus configuration", "Done.")

	// Verbose and build the clean (well-typed) actions configuration
	for i := 0 ; i < len(configuration_tmp.Actions); i ++ {
		configuration.Actions = append(configuration.Actions,
			Action_t {
				ActionTypeOfString(configuration_tmp.Actions[i].Type),
				configuration_tmp.Actions[i].Title,
				configuration_tmp.Actions[i].Ref })
	}
	logging("Actions configuration", "Done.")

	return *configuration
}

