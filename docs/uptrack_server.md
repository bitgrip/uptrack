## uptrack server

Interact with the uptrack server

### Synopsis

Configure and interact with the uptrack server

### Options

```
      --check-frequency int          Default interval to execute job (default 10)
      --datadog-apiKey string        Datadog API-Key
      --datadog-appKey string        Datadog APP-key
      --datadog-endpoint string      Datadog Endpoint (default "https://api.datadoghq.com/api/v1/series")
      --datadog-interval int         Interval for sending metrics to Datadog (default 5)
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
* [uptrack server start](uptrack_server_start.md)	 - Start the uptrack server

