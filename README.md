UpTrack - Uptime Tracker
========================

UpTrack is a service to track down and monitor HTTP Endpoints conditions





## Configuration
Configuration of the service and setting the job Descriptor is possible via params (see [Documentation](docs/uptrack.md)	) and environment variables.

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
Even though you may configure the uptrack server via `UPTRACK_CONFIG`, you can overwrite most 

```
UPTRACK_JOBS_CONFIG
location of the jobs descriptor, see /config/jobs.yaml as template
   
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


docker run -v "$(pwd)/config":"/go/config" \
-e UPTRACK_CONFIG=/go/config/uptrack.yaml \
-e UPTRACK_JOBS_CONFIG=/go/config/jobs.yaml \
uptrack server start



#### For a more details, check out the full [Documentation](docs/uptrack.md)	

