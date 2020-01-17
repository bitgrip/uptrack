## uptrack server start

Start the server

### Synopsis

Start the UpTrack server

```
uptrack server start [flags]
```

### Options

```
  -h, --help   help for start
```

### Options inherited from parent commands

```
      --config string                config file (default is $HOME/.uptrack.yaml)
      --datadog-credentials string   File containing datadog credentials (default "/etc/uptrack/datadog/credentials")
      --default-interval int         Default interval to execute job (default 10)
      --jobs-dir string              Directory to find job descriptors (default "/etc/uptrack/jobs.d")
      --log-json                     if to log using json format
      --prometheus-scrape string     Endpoint prometheus can scrape from (default ":9001/metrics")
  -v, --verbosity int                verbosity level to use
```

### SEE ALSO

* [uptrack server](uptrack_server.md)	 - Interact to the server

