package main 

import (
	"net/http"
	"os"
	"strings"
	"bytes"
	"time"
)

/**
* Triggers are built from the triggers.json file. They are used to determine how the requests the server generates 
* should be handled. ENV('') can be used to pull things from the environment file and use them in the triggers.json file.
*/
type Trigger struct {
	Name 					string `json:"name"`		// Name of the trigger for debugging purposes and human readability. Ex "My_Super_Cool_Trigger"
	Endpoint 				string `json:"endpoint"`	// The endpoint to hit for the trigger. Ex "https://api.example.com"
	Method 					string `json:"method"`		// The HTTP method we will use for this. Ex "GET", "POST", "PUT", "PATCH", "DELETE"
	Route 					string `json:"route"`		// The route on the endpoint. Ex "/events/doSomething"
	Duration 				string `json:"duration"`	// How often formatted for https://cs.opensource.google/go/go/+/go1.23.4:src/time/format.go;l=1617. Ex "1s", "1m", "1h", "1d"
	
	AdditionalHeaders 		[]struct {					// Additional headers that may be needed for the request. Could pass API keys here. 
								Key 	string `json:"key"` 		// The key of the header.
								Value 	string `json:"value"` 		// The value of the header. 
							} `json:"additional_headers"`
}

/**
* Adds support for pulling from .env file when the trigger configuration calls for it. 
* This checks to see if the value is prefixed with ENV(' and suffixed with ') and if it is,
* it will pull the value from the .env file. If it is not, it will return the value as is.
* @param value (string) The value inside of triggers.json
*/
func buildValue(value string) string {
	ret := value; 
	if strings.HasPrefix(value, "ENV('") && strings.HasSuffix(value, "')") {
		envValue := strings.TrimSuffix(strings.TrimPrefix(value, "ENV('"), "')")
		ret = os.Getenv(envValue);
	}
	return ret;
}

/**
* Build the http request for the trigger and return it. 
* @param trigger (struct) Trigger
*/
func buildRequest(trigger Trigger) *http.Request {
	// Build empty body request. 
	body := []byte(`{}`);

	request, err := http.NewRequest(buildValue(trigger.Method), (buildValue(trigger.Endpoint)+buildValue(trigger.Route)), bytes.NewBuffer(body));
	
	if(err != nil){
		logText("There was an error building the request for trigger: "+buildValue(trigger.Name), err);
		return nil;
	}

	if(len(trigger.AdditionalHeaders) > 0){
		// Build headers for the request.
		for _, header := range(trigger.AdditionalHeaders) {
			request.Header.Add(buildValue(header.Key), buildValue(header.Value));
		}
	}

	// Set User-Agent for the request to be the sync server. 
	request.Header.Add("User-Agent", "Web-Sync-Server");

	return request;
}

/**
* @param trigger (struct) Trigger
* @param processNow (bool) If false, the trigger will be scheduled to be processed and not do so immediately.
* Process the trigger struct that was passed in.
*/
func processTrigger(trigger Trigger, processNow bool){
	procType := "process"
	if(!processNow){
		procType = "schedule";
	}

	logText(
		"Attempting to "+procType+" trigger: "+buildValue(trigger.Name)+
		" REQUEST TYPE: "+buildValue(trigger.Method)+ 
		" ON: "+buildValue(trigger.Endpoint)+buildValue(trigger.Route),
	nil);

	/**
	* If we need to process this now and not just schedule it, do so now. 
	*/
	if(processNow){
		req := buildRequest(trigger);
		_, err := client.Do(req);

		if(err != nil){
			logText("There was an error processing trigger: "+buildValue(trigger.Name), err);
		}else {
			logText("Successfully processed trigger: "+buildValue(trigger.Name), nil);
		}
	}

	/**
	* Schedule the trigger to be processed. 
	*/
	gDuration, err := time.ParseDuration(buildValue(trigger.Duration));
	if(err != nil){
		logText("There was an error parsing the duration for trigger: "+buildValue(trigger.Name), err);
		return;
	}

	time.AfterFunc(gDuration, func(){ processTrigger(trigger, false)});
}


func validateTrigger(trigger Trigger) bool{
	// TODO: Implement this function.
	return true;
}


