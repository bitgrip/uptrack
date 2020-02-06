## uptrack server

Interact to the server

### Synopsis

Interaction to the UpTrack server

### Options

```
      --datadog-credentials string   File containing datadog credentials (default "/etc/uptrack/datadog/credentials")
      --default-interval int         Default interval to execute job (default 10)
  -h, --help                         help for server
      --jobs-config string              Directory to find job descriptors (default "/etc/uptrack/jobs.d")
      --prometheus-scrape string     Endpoint prometheus can scrape from (default ":9001/metrics")
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.uptrack.yaml)
      --log-json        if to log using json format
  -v, --verbosity int   verbosity level to use
```

### SEE ALSO

* [uptrack](uptrack.md)	 - track down your uptime
* [uptrack server start](uptrack_server_start.md)	 - Start the server

