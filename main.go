package main

import (
	"net/http"
	"encoding/json"
	"flag"
	"bytes"
	"time"
	//"strconv"
	"os"
	"log"
	//"errors"
    "github.com/joho/godotenv"
)

/******************************************************************************************
* 
	Configuration 
*
*******************************************************************************************/

var (
    serverPort             string
    endPoint               string
    erpKey                 string
    errorLog               string
    commandLog             string
    trackingUpdateInterval string
    eiaUpdateIntervalDays  string
    client                 = &http.Client{}
	debugMode			   bool
	triggers			   []byte
)

const (
	redColor    = "\033[31m"
	yellowColor = "\033[33m"
	resetColor  = "\033[0m"
)

func init() {
	flag.BoolVar(&debugMode, "debugMode", false, "Enable Debug Mode");
	flag.Parse();

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort = os.Getenv("SERVER_PORT");
	endPoint = os.Getenv("END_POINT");
	erpKey = os.Getenv("ERP_KEY");
	errorLog = os.Getenv("ERROR_LOG");
	commandLog = os.Getenv("COMMAND_LOG");
    trackingUpdateInterval = os.Getenv("TRACKING_UPDATE_INTERVAL");
    eiaUpdateIntervalDays = os.Getenv("EIA_UPDATE_INTERVAL_DAYS");

	tFile, tErr := os.OpenFile("triggers.json", os.O_RDONLY, 0444);
	if tErr != nil {
		panic(tErr);
	}
	decoder := json.NewDecoder(tFile)
	token, dErr := decoder.Token();
	if(dErr != nil){
		panic(dErr);
	}
	print(token);
	if debugMode {
		outputDebugInfo();
	}
}

func main() {
	triggerLocationUpdates();
	triggerEIAFuelPrices();

	select {};
}


/******************************************************************************************
* 
*	Triggers
*
*******************************************************************************************/

func processTrigger(triggerName string){

}

/**
* This function is used to trigger location updates for 
* the ERoad ELD Integration built into the Symfony API. 
*/
func triggerLocationUpdates(){
	logText("Attempting to process ERoad location updates....", nil);
	// Create an empty request.
	body := []byte(`{}`);

	// Create a new post request for the eRoadLocationEndpoint and give it the API Key.  
	request, err := http.NewRequest("GET", endPoint+"events/eRoadLocationUpdates", bytes.NewBuffer(body));
	request.Header.Add("Authorization", erpKey);

	response, err := client.Do(request);

	if(err != nil){
		logText("There was an error gathering location updates.", err);
	}else {
		logText("Location updates were successful. \nNext update in "+trackingUpdateInterval+" minutes...", nil);
		defer response.Body.Close();
		duration, tError := time.ParseDuration(trackingUpdateInterval+"m");
		if(tError != nil){
			logText("There was an error parsing the Tracking Update Interval. It will not run again until server reset.", tError);
		} else {
			time.AfterFunc(duration, triggerLocationUpdates);
		}
	}
}

/**
* This function is used to trigger updates for fuel market price
* averages from the EIA. (https://www.eia.gov/). It is updated every 8 days.
*/
func triggerEIAFuelPrices(){
	logText("Attempting to process EIA fuel market price updates....", nil);
	// Empty body request
	body := []byte(`{}`);

	request, err := http.NewRequest("GET", endPoint+"events/weeklyEIADieselPrices", bytes.NewBuffer(body));
	request.Header.Add("Authorization", erpKey);

	response, err := client.Do(request);

	if(err != nil){
		logText("There was an error gathering new fuel prices", err);
	}else {
		logText("Fuel Updates were successful. \nNext update in "+eiaUpdateIntervalDays+" days...", nil);
		defer response.Body.Close();
		duration, tError := time.ParseDuration(eiaUpdateIntervalDays+"h");
		if(tError != nil){
			logText("There was an error parsing the EIA Update Interval. It will not run again until server reset.", tError);
		}else {
			time.AfterFunc(duration, triggerEIAFuelPrices);
		}
	}
}