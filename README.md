# Valiant Sync Server

This is a server written in Golang to handle the scheduling of certain calls to an endpoint defined in the .env. Currently what those endpoint are is hard coded but it's ultimately up to the end point to decide what it will do with the requests send by this sync server.

The point of this is to allow for poling requests without having to depend on a cron job that may vanish during an infrastructure change. The only thing that would need to be done is running the SyncServer like a Deamon and well... Go!


## Example .env file configuration**

.env
```
SERVER_PORT="3333"
END_POINT="
ERP_KEY=""
ERROR_LOG="./var/log/error.log"
COMMAND_LOG="./var/log/command.log"
TRACKING_UPDATE_INTERVAL="5"
EIA_UPDATE_INTERVAL_DAYS="8"
```

`END_POINT` - The endpoint where the requests that are within the script will be sent.

`ERP_KEY` - The key that will be added in the request's authorization header. 

`ERROR_LOG` - Where error logs will be stored. By default it is in /var/log of the binary. The package should include some examples of errors. 

`COMMAND_LOG` - A log of what is to be logged from this request. 

## Keys to be deprecated in the future:

`TRACKING_UPDATE_INTERVAL` - Currently used to trigger truck tracking updates every 5 minutes, in this example.

`EIA_UPDATE_INTERVAL_DAYS` - Currently used to trigger EIA Fuel price updates every 8 days, in this example. 

`END_POINT` - The endpoint where the requests that are within the script will be sent. In the future we will simply define URLS. 

`SERVER_PORT` - The port that will be used by the Sync Server. **Not yet implemented** however it will accept connections in the future via web sockets. 

## Flags
| Flag | Type | Values | Description |
|------|------|--------|-------------|
|`--debugMode`|*Bool*|(true/false)| This command will output logs to the terminal for viewing along with the current set environment variables to ensure they are being seen by the application. 
