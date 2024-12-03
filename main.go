package main 

import (
	"fmt"
	"net/http"
	"bytes"
	"time"
	"strconv"
	"os"
	"log"
)

/******************************************************************************************
* 
																	Configuration 
*
*******************************************************************************************/

const serverPort = 3333; 
const endPoint = "http://localhost:8000/";
const erpKey = "Dwuj334Na2JMu61QBsQzj3KBd4uKuVZs-internal";
const errorLog = "./var/log/error.log";
const commandLog = "./var/log/command.log";
const trackingUpdateInterval = 5; 

func main() {
	triggerLocationUpdates();
}


/******************************************************************************************
* 
	Triggers
*
*******************************************************************************************/


/**
* This function is used to trigger location updates for 
* the ERoad ELD Integration built into the Symfony API. 
*/
func triggerLocationUpdates(){
	logText("Attempting to process ERoad location updates....", false);
	// Create an empty request.
	body := []byte(`{}`);

	// Create a new post request for the eRoadLocationEndpoint and give it the API Key.  
	request, err := http.NewRequest("GET", endPoint+"events/eRoadLocationUpdates", bytes.NewBuffer(body));
	request.Header.Add("Authorization", erpKey);

	client := &http.Client{};

	response, err := client.Do(request);

	if(err != nil){
		logText("There was an error gathering location updates.", true);
		panic(err)
	}else {
		logText("Location updates were successful. \nNext update in "+strconv.Itoa(trackingUpdateInterval)+" minutes...", false);
		defer response.Body.Close();
		time.Sleep(trackingUpdateInterval * time.Minute);
		triggerLocationUpdates();
	}
}

/**
* This function is used to trigger updates for fuel market price
* averages from the EIA. (https://www.eia.gov/). It is updated every monday.
*/
func triggerEIAPrices(){
}

/******************************************************************************************
* 
	Utility Methods 
*
*******************************************************************************************/

func logText(sToLog string, isError bool) bool{
	if(isError){
		file, err := os.OpenFile(errorLog, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666);
	}else {
		file, err := os.OpenFile(commandLog, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666);
	}
	if(err != nil){
		logText("ERROR: Could not write to command log.", true);
		panic(err);
	}
	defer file.Close(); 

	log.setOutput(file);
	log.Println(sToLog);
	return true;
}