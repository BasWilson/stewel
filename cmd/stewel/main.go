package main

import (
	"encoding/json"
	"os"

	reverse_proxy "github.com/baswilson/stewel/lib"
)

func handleLocalFile(jsonFile string) reverse_proxy.Config {
	file, err := os.Open(jsonFile)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	var config reverse_proxy.Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err.Error())
	}

	return config
}

func main() {

	configPath := "./stewel-config.json"

	if len(os.Args) == 2 {
		configPath = os.Args[1]
	}

	config := handleLocalFile(configPath)

	reverse_proxy.Create(":80", config)
}
