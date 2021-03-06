## uptrack

track down your uptime

### Synopsis

**uptrack** is a service to steadily check the uptime of your HTTP services.
You define **Jobs** (UpJobs/DNSJobs), consisting of checks currently available checks are:

For UpJobs:
* Connection to given url possible
* Tme to certification expiry
* Time to connect
* Time to first byte
* Request time
* Received bytes

For DNSJobs
* Ratio of found ips

Find a fully filled job configuration here: [Jobs configuration](../config/jobs.yaml)



### Options

```
  -h, --help                    help for uptrack
      --log-json                if to log using json format
      --uptrack-config string   Configuration file for uptrack application  (default "config/uptrack.yaml")
  -v, --verbosity int           verbosity level to use
```

### SEE ALSO

* [uptrack completion](uptrack_completion.md)	 - Generates bash completion scripts
* [uptrack gen-doc](uptrack_gen-doc.md)	 - genrates the markdown documentation
* [uptrack server](uptrack_server.md)	 - Interact with the uptrack server
* [uptrack version](uptrack_version.md)	 - version of uptrack

