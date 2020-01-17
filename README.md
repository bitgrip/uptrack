UpTrack - Uptime Tracker
========================

UpTrack is a service to track down and monitor HTTP Endpoints conditions



## Configuration
Configuration of the service and setting the job Descriptor is possible via params 
(see [Documentation](docs/uptrack.md)	) and environment variables.

## Configuation via Environment Variables
### Basic Configuration
```
UPTRACK_CONFIG:
main configuration file for uptrack 

VERBOSITY
define verbosity level 0, 1, 2,

LOG_JSON
if set true, logs output is in JSON format
```

### Server Configuration
Even though you may configure the uptrack server via `UPTRACK_CONFIG`, you can overwrite the parameters via environment variables. 

```
UPTRACK_JOBS_CONFIG
location of the jobs descriptor, see below
   
UPTRACK_CHECK_FREQUENCY
how often checks are performed
   
UPTRACK_PROMETHEUS_PORT
port where prometheus server is exposed

UPTRACK_PROMETHEUS_ENDPOINT
availabel prometheus endpoint for scraping
   
UPTRACK_DATADOG_ENDPOINT
url to send the datadog metrics to
   
UPTRACK_DATADOG_APPKEY
Datadog App-Key
   
UPTRACK_DATADOG_APIKEY
DataDog API-Key   

UPTRACK_DATADOG_INTERVAL
how often the metrics are sent to DataDog

```

### Jobs Configuration
There are two kinds of jobs to configure:
#### Up-Job
- checks, if a website/service is online(set expected status code)
- gauge connection time, request time, etc...

#### DNS-Job
- checks if DNS-Resolution works as expected



### jobs.yaml
Configure your jobs via a 'jobs.yaml' as follows

```
project: Project Name
datadog_enabled: false
prometheus_enabled: true
up_jobs:
    up_job_1:
        description: 'job description here'
        host: 'example.com'
        url: https://www.example.com/auth
        method: GET
        header:
            Authorization:
                - 'Basic abc123=='
        PlainBody: '{"name":"stuff"}'
        Base64Body: eyJuYW1lIjoic3R1ZmYifQ==
        CheckSSL: true
dns_jobs:
    bitgrip_dns:
        host: 'example.com'
        fqdn: 'www.example.com'
        ips:
            - '123.223.47.12'
            - '123.45.67.18'

```


### Start Locally
#### as Docker
``` 
UPTRACK_IMAGE_TAG=master
 docker run -v "$(pwd)/config":"/go/config" \
 -e UPTRACK_CONFIG=/go/config/uptrack.yaml \
 -e UPTRACK_JOBS_CONFIG=/go/config/jobs.yaml \
 bitgrip/uptrack:${UPTRACK_IMAGE_TAG} server start
```




#### For a more details, check out the full [Documentation](docs/uptrack.md)	

