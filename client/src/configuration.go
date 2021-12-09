package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"flag"
    "rts/utils"
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
	Port	 int `json:"Port"`
}
// Keys are represented by strings in the configuration file. This will allows
// us to use words representing non-ascii keys.
type Keys_t struct {
	Left int32 `json:"Left"`
	Right int32 `json:"Right"`
	Up int32 `json:"Up"`
	Down int32 `json:"Down"`
	ZoomIn int32 `json:"ZoomIn"`
	ZoomOut int32 `json:"ZoomOut"`
	Menu int32 `json:"Menu"`
	ResetCamera int32 `json:"ResetCamera"`
	Chat int32 `json:"Chat"`
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
		utils.Logging("Configuration", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}

	// Convert the json content into a Configuration structure
	var configuration = &Configuration_t{}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		utils.Logging("Configuration", fmt.Sprintf("Cannot parse the config file: %v", err))
		os.Exit(1)
	}

	// Parse the command line and overwrite the configuration if needed
	override_addr_parsed := flag.String("n", configuration.Server.Hostname, "Hostname of the server (ip address or name)")
	override_port_parsed := flag.Int("p", configuration.Server.Port, "Port of the remote server")
	flag.Parse()
	// Replace the configuration variables
	configuration.Server.Hostname = *override_addr_parsed
	configuration.Server.Port = *override_port_parsed

	utils.Logging("Configuration",
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
	Ref int
}
type Menu_t struct {
	Ref int
	Title string
	Elements []MenuElement_t
}
type Action_t struct {
	Type ActionType
	Ref int
	Title string
}
type MenuConfiguration_t struct {
	Menus []Menu_t
	Actions []Action_t
}

func loadTextMenus(file_name string) MenuConfiguration_t {
	// Read the menus configuration file
	file, err := ioutil.ReadFile(file_name)
	if err != nil {
		utils.Logging("Menu configuration", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}

	// Convert the json content into a Configuration structure
	var configuration_tmp = &MenuConfiguration_tmp_t{}
	var configuration = &MenuConfiguration_t{}
	err = json.Unmarshal(file, &configuration_tmp)
	if err != nil {
		utils.Logging("Menu configuration", fmt.Sprintf("Cannot parse the config file: %v", err))
		os.Exit(1)
	}

	// List the different references of the menus and actions (dirty-typed)
	// in order to build the mapping "dirty_ref" ->~n_{"dirty ref"}.
	var menus_ref_mapping = make(map[string]int)
	for i:= 0 ; i < len(configuration_tmp.Menus); i ++ {
		// The index 0 is reserved for the Main menu
		if configuration_tmp.Menus[i].Ref == "Main" {
			menus_ref_mapping[configuration_tmp.Menus[i].Ref] = 0
		} else {
			menus_ref_mapping[configuration_tmp.Menus[i].Ref] = i + 1
		}
	}
	var actions_ref_mapping = make(map[string]int)
	for i:= 0 ; i < len(configuration_tmp.Actions); i ++ {
		actions_ref_mapping[configuration_tmp.Actions[i].Ref] = i + 1
	}

	// Verbose and build the clean (well-typed) menus configuration
	configuration.Menus = nil
	configuration.Actions = nil
	for i := 0; i < len(configuration_tmp.Menus); i++ {
		// Declaration of local variables to store the current menu to add
		var active_menu Menu_t
		var inner_elements = []MenuElement_t(nil)

		active_menu.Title = configuration_tmp.Menus[i].Title
		active_menu.Ref = menus_ref_mapping[configuration_tmp.Menus[i].Ref]

		for j := 0 ; j < len(configuration_tmp.Menus[i].Elements); j ++ {
			loc_type :=  MenuElementTypeIfString(configuration_tmp.Menus[i].Elements[j].Type)
			var idx int

			switch (loc_type) {
				case MenuElementAction:
				idx = actions_ref_mapping[configuration_tmp.Menus[i].Elements[j].Ref]
				case MenuElementSubMenu:
				idx = menus_ref_mapping[configuration_tmp.Menus[i].Elements[j].Ref]
			}

			inner_elements = append(inner_elements,
				MenuElement_t {
					configuration_tmp.Menus[i].Elements[j].Name,
					MenuElementTypeIfString(
						configuration_tmp.Menus[i].Elements[j].Type),
					utils.KeyOfString(configuration_tmp.Menus[i].Elements[j].Key),
					idx })
		}

		active_menu.Elements= inner_elements
		configuration.Menus = append(configuration.Menus, active_menu)
	}
	utils.Logging("Menus configuration", "Done.")

	// Verbose and build the clean (well-typed) actions configuration
	for i := 0 ; i < len(configuration_tmp.Actions); i ++ {
		configuration.Actions = append(configuration.Actions,
			Action_t {
				Type: ActionTypeOfString(configuration_tmp.Actions[i].Type),
				Title: configuration_tmp.Actions[i].Title,
				Ref: actions_ref_mapping[configuration_tmp.Actions[i].Ref] })
	}
	utils.Logging("Actions configuration", "Done.")

	return *configuration
}

