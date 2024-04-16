## Stewel Reverse Proxy 🥾

Dead simple reverse proxy based on [httputil.ReverseProxy](https://golang.org/pkg/net/http/httputil/#ReverseProxy).

Just set a list of hosts and targets and you're good to go.

Host is the domain you want to proxy, targets are the servers you want to proxy to.

Example usage as a library:

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
        },
    },
}

reverse_proxy.Create(":80", config)
```

This will create a reverse proxy that listens on port 80 and proxies requests to `yourhost.com` to `localhost:4000` and `localhost:4001` in a round robin fashion.

You can use it as a library or as a standalone binary.

To use it as a standalone binary, just run:

```bash
make compile
./bin/stewel yourconfig.json
```

You can also use the provided Dockerfile to build a Docker image.
Pass the path to your config file as a volume.
You might encounter issues with network routing when using the Docker image, so make sure to set the `--network=host` flag when running the container.

Sample Dockerfile usage:

```bash
docker build -t yourtag .
docker run -p 80:80 --network=host -v /path/to/stewel-config.json:/app/stewel-config.yaml yourtag
```

### Todo

-   [ ] Add tests
-   [ ] Add more load balancing methods
-   [ ] Add automatic SSL certificate generation with Let's Encrypt [lego](https://github.com/go-acme/lego)
