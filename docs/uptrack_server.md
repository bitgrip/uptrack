## uptrack server

Interact to the server

### Synopsis

Start, configure and interact with the uptrack server

### Options

```
      --datadog-apiKey string        Datadog API-Key
      --datadog-appKey string        Datadog APP-key
      --datadog-interval int         Interval for sending metrics to Datadog (default 5)
      --default-interval int         Default interval to execute job (default 10)
  -h, --help                         help for server
      --jobs-config string           Descriptor file defining all checks (default "./config/jobs.yaml")
      --prometheus-endpoint string   Prometheus Endpoint (default "/metrics")
      --prometheus-port string       Port exposed by prometheus (default "9001")
```

### Options inherited from parent commands

```
      --log-json                if to log using json format
      --uptrack-config string   Configuration file for uptrack application  (default "config/uptrack.yaml")
  -v, --verbosity int           verbosity level to use
```

### SEE ALSO

* [uptrack](uptrack.md)	 - track down your uptime
* [uptrack server start](uptrack_server_start.md)	 - Start the server

