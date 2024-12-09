# GS Sync Server

This is a server written in Golang to handle the scheduling of certain calls to an endpoint defined in the `triggers.json` configuration.

The point of this is to allow for poling requests without having to depend on a cron job that may vanish during an infrastructure change. The only thing that would need to be done is running the SyncServer like a Daemon and well... Go!

## Building
Simply run `go build` within the project directory. This should compile a binary for you to use. 

## Example .env file configuration
At the minimum you need these. It's worth noting that variables from your .env file can be called in your `triggers.json` configuration.

.env
```
ERROR_LOG="./var/log/error.log"
COMMAND_LOG="./var/log/command.log"
```

`ERROR_LOG` - Where error logs will be stored. By default it is in /var/log of the binary. The package should include some examples of errors. 

`COMMAND_LOG` - A log of what is to be logged from this request. 

## The triggers.json file. 
This file is the bread and butter of your configuration. It is a .json file that contains a list of objects that define how the server should handle your requests to certain end-points. We have enclosed a `triggers.example.json` to give you an idea of how you could set your configuration up. Here it is below: 

```
[
	{
		"name":"asset_position_update",
		"endpoint":"http://example.com",
		"method":"GET",
		"route":"/events/asset-location-updates",
		"duration":"1m"
	}, 
	{
		"name":"useia_fuel_price_update",
		"endpoint":"ENV('END_POINT')",
		"method":"GET",
		"route":"/events/weeklyEIADieselPrices",
		"duration":"128h",
		"additional_headers":[{
			"key":"Authorization",
			"value":"ENV('AUTH_TOKEN')"
		}]	
	}
]
```
### Options Explained: 
Options marked with * are required. 

| Option | Description |
|--------|-------------|
|*name    |The name of this specific trigger for logging purposes.|
|*endpoint|The base URL of what we are calling. In the example above. Should include http:// or https:// in the URL. Don't put a trailing slash.|
|*method  |What method is to be used in the request.|
|*route   |The route on the endpoint that we are accessing.|
|*duration|A string defining how often we should run this trigger. Valid strings should follow the inputs of this function - https://cs.opensource.google/go/go/+/go1.23.4:src/time/format.go;l=1617 Ex: "1s", "1m", "1h", "1d"|
|additional_headers| An array of objects defining any additional headers that should be sent with your request.|

### A note on ENV()

All parameters within your triggers.json configuration can pull values from the .env file. This is to stop you from repeating yourself and allow for easy testing on development if need be. 


## Flags

These are flags that can be passed in when launching the binary. 

| Flag | Type | Values | Description |
|------|------|--------|-------------|
|`--debugMode`|*Bool*|(true/false default:false)| This flag will output logs to the terminal for viewing along with the current set environment variables to ensure they are being seen by the application. |
|`--warmUp`|*Bool*|(true/false default:true)| When set to true, this flag will cause everything within your triggers.json configuration to fire once before scheduling as configured.

### Author Notes
This is my first project I'd written in Golang outside of some practice here and there. I found it to be pretty insightful into the pros and cons of using the language and I think I'll be using it more in the future. :) 

I am open to forks and issues opened on this. Constructive criticism is always welcome. Feel free to review the code.