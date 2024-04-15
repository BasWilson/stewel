package reverse_proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Config struct {
	Hosts []Host `json:"hosts" validate:"required"`
}

type Host struct {
	Host         string       `json:"host" validate:"required"`
	Targets      []Target     `json:"targets" validate:"required"`
	LoadBalancer LoadBalancer `json:"loadBalancer"`
}

type Target string

type LoadBalancer struct {
	Method LoadBalancingMethod `json:"method" validate:"required"`
}

type LoadBalancingMethod string

const (
	RoundRobin LoadBalancingMethod = "round-robin"
)

type Instance struct {
	Config Config
}

var instance Instance

func handleConnection(req *http.Request) {
	hostHeader := req.Host

	var foundHost Host

	for _, host := range instance.Config.Hosts {
		println("Comparing", host.Host, hostHeader)
		if host.Host == hostHeader {
			foundHost = host
			break
		}
	}

	if foundHost.Host != "" {
		target, err := url.Parse(string(foundHost.Targets[0]))

		if err != nil {
			req.Response.StatusCode = http.StatusBadGateway
		}

		fmt.Println("proxied to target: ", target)
		req.Header.Set("X-stewelVersion", "1")

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
	} else {
		req.Response.StatusCode = http.StatusNotFound
	}

}

func Create(config Config) error {

	instance = Instance{
		Config: config,
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{})
	proxy.Director = handleConnection
	err := http.ListenAndServe(":1337", proxy)

	return err
}
