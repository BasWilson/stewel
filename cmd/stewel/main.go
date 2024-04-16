package main

import (
	reverse_proxy "github.com/baswilson/adraptor/tools/updater/lib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		println("Skipping .env file")
	}

	config := reverse_proxy.Config{
		Hosts: []reverse_proxy.Host{
			{
				Host: "shadow.adraptor.network",
				Targets: []reverse_proxy.Target{
					"http://localhost:4000",
					"http://localhost:4001",
				},
				LoadBalancer: reverse_proxy.LoadBalancer{
					Method: reverse_proxy.RoundRobin,
				},
			},
		},
	}

	reverse_proxy.Create(":80", config)
}
