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
if settrue, logs output is in JSON format
```

### Server Configuration
Even though you may configure the uptrack server via `UPTRACK_CONFIG`, you can overwrite most 


docker run -v "$(pwd)/config":"/go/config" \
-e UPTRACK_CONFIG=/go/config/uptrack.yaml \
-e UPTRACK_JOBS_CONFIG=/go/config/jobs.yaml \
uptrack server start



#### For a more details, check out the full [Documentation](docs/uptrack.md)	

