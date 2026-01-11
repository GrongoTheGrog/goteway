# Goteway

Goteway is a simple, easy and powerful api gateway written in Golang, as the name might suggest. 


## How to use it?

In order to use it, one must decide between the declarative or the imperative configuration. 

### Imperative

The imperative way consists in setting and configuring your routes in the code level. 

- First, initialize the gateway and routes

```go
gw := gateway.NewGateway()

route1 := gateway.NewRoute("http://service"). 
            PathPattern("/service1/*", "/v1/service1/*"). 
            Header("X-Service", "service1").             
            Hosts("frontend.com", "www.frontend.com").   
            Methods("GET", "OPTIONS")


```
The benefit of that approch is that you can attatch your own custom filters to each route or in a global context as folllows: 

```go
customFilter := filter.NewBasicFilter(func (ctx *filter.Context) *http.Response{
  ctx.Log("Logging method from custom filter: %s", ctx.Request.Method)
  return ctx.RunNextFilter()
})

service1.Filter(customFilter)	// Binds the filter to the route
gw.AddFilter(customFilter)		// Binds the filter to the entire gateway
```

- Then attach the route to the gateway

```go
gw.AddRoute(service1)
```

- Then start it

```go
gw.Start(":9000")
```

### Declarative

The declarative way is done by editing a yml file called goteway.yml in the root directory of where the program is running.

```Yaml
gateway: 
  log_filter:                       // logs the following informations about the request/response
    - method
    - path
    - latency
    - full_url
    - status


routes:
  - name: service1
    enabled: true
    endpoint: "http://service1"
    paths:
      - "/service1/*"

    rate_limiting:
      enabled: true                  
      interval_seconds: 60           
      max_requests: 10               // max requests within the interval
      resource: "user"               // this decides what resouce the filter is limiting; includes "user", "gateway" and "route"

    strip_prefix: 1                  // removes x prefixes from the path
    strip_suffix: 0                  // removes x suffixes from the path
    status: 201                      // changes response's status
    header_config:                   // add headers to response
      - name: something
        value: newValue


```
