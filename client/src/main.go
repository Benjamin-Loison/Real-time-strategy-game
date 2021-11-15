package main

import (
	"io/ioutil"
	"fmt"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
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
		logging("Conf loader", fmt.Sprintf("Cannot open the config file: %v", err))
		os.Exit(-1)
	}
	//defer file.Close()

	var configuration Configuration
	var conf map[string]interface{}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		logging("Conf loader", fmt.Sprintf("Cannot parse the config file: %v", err))
		configuration = Configuration {hostname: "138.231.144.134", port: 80}
	} else {
		configuration = Configuration {
			hostname: conf["hostname"].(string),
			port: int(conf["port"].(float64))}
	}

	// Parse the command line
	override_addr_parsed := flag.String("n", configuration.hostname, "Hostname of the server (ip address or name)")
	override_port_parsed := flag.Int("p", configuration.port, "Port of the remote server")
	flag.Parse()

	configuration.hostname = *override_addr_parsed
	configuration.port = *override_port_parsed

	logging("Conf loader",
		fmt.Sprintf("Hostname: %s, port: %d", configuration.hostname, configuration.port))

	return configuration
}

func main(){
	// Load the configuration
	config := loadConfig("conf/conf.json")

	// starting the client
	go startClient(&client_chan, config)

	// initializing the game
	g := &Game{}
	Init(g)
	//initializing ebiten
	ebiten.SetWindowSize(screenWidth,screenHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("EHO: Elves, humans and orks")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
