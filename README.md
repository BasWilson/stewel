## Stewel Reverse Proxy 🥾

Dead simple reverse proxy based on [httputil.ReverseProxy](https://golang.org/pkg/net/http/httputil/#ReverseProxy).

Just set a list of hosts and targets and you're good to go.

Host is the domain you want to proxy, targets are the servers you want to proxy to.

Example usage:

```go
config := reverse_proxy.Config{
    Hosts: []reverse_proxy.Host{
        {
            Host: "yourhost.com",
            Targets: []reverse_proxy.Target{
                "http://localhost:4000",
                "http://localhost:4001",
            },
            LoadBalancer: reverse_proxy.LoadBalancer{
                Method: reverse_proxy.RoundRobin,
            },
        },
        {
            Host: "anotherhost.com",
            Targets: []reverse_proxy.Target{
                "http://192.168.178.20:4000",
            },
            LoadBalancer: reverse_proxy.LoadBalancer{
                Method: reverse_proxy.RoundRobin,
            },
        },
    },
}

reverse_proxy.Create(":80", config)
```

This will create a reverse proxy that listens on port 80 and proxies requests to `yourhost.com` to `localhost:4000` and `localhost:4001` in a round robin fashion.
