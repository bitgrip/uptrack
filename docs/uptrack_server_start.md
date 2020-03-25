## uptrack server start

Start the uptrack server

### Synopsis

Start the uptrack server

```
uptrack server start [flags]
```

### Options

```
  -h, --help   help for start
```

### Options inherited from parent commands

```
      --check-frequency int          Default interval to execute job (default 10)
      --datadog-apiKey string        Datadog API-Key
      --datadog-appKey string        Datadog APP-key
      --datadog-endpoint string      Datadog Endpoint (default "https://api.datadoghq.com/api/v1/series")
      --datadog-interval int         Interval for sending metrics to Datadog (default 5)
      --jobs-config string           Descriptor file defining all checks (default "./config/jobs.yaml")
      --log-json                     if to log using json format
      --prometheus-endpoint string   Prometheus Endpoint (default "/metrics")
      --prometheus-port string       Port exposed by prometheus (default "9001")
      --uptrack-config string        Configuration file for uptrack application  (default "config/uptrack.yaml")
  -v, --verbosity int                verbosity level to use
```

### SEE ALSO

* [uptrack server](uptrack_server.md)	 - Interact with the uptrack server

