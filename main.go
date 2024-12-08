package main 

import (
	"net/http"
	"bytes"
	"time"
	"strconv"
	"os"
	"log"
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
    trackingUpdateInterval int
    eiaUpdateIntervalDays  int
    client                 = &http.Client{}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort = os.Getenv("SERVER_PORT");
	endPoint = os.Getenv("ENDPOINT");
	erpKey = os.Getenv("ERP_KEY");
	errorLog = os.Getenv("ERROR_LOG");
	commandLog = os.Getenv("COMMAND_LOG");
    trackingUpdateInterval, _ = strconv.Atoi(os.Getenv("TRACKING_UPDATE_INTERVAL"));
    eiaUpdateIntervalDays, _ = strconv.Atoi(os.Getenv("EIA_UPDATE_INTERVAL_DAYS"));

	if(eiaUpdateIntervalDays > 102200 || trackingUpdateInterval > 102200){
		logText("ERROR: EIA Update Interval or tracking update interval is too large. Please set it to a reasonable value.", true);
		panic("ERROR: EIA Update Interval or tracking update interval is too large. Please set it to a reasonable value.");
	}
}

func main() {
	triggerLocationUpdates();
	triggerEIAFuelPrices();
}


/******************************************************************************************
* 
*	Triggers
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

	response, err := client.Do(request);

	if(err != nil){
		logText("There was an error gathering location updates.", true);
		panic(err)
	}else {
		logText("Location updates were successful. \nNext update in "+strconv.Itoa(trackingUpdateInterval)+" minutes...", false);
		defer response.Body.Close();
		time.AfterFunc(time.Duration(trackingUpdateInterval) * time.Minute, triggerLocationUpdates);
	}
}

/**
* This function is used to trigger updates for fuel market price
* averages from the EIA. (https://www.eia.gov/). It is updated every monday.
*/
func triggerEIAFuelPrices(){
	logText("Attempting to process EIA fuel market price updates....", false);
	// Empty body request
	body := []byte(`{}`);

	request, err := http.NewRequest("GET", endPoint+"events/weeklyEIADieselPrices", bytes.NewBuffer(body));
	request.Header.Add("Authorization", erpKey);

	response, err := client.Do(request);


	if(err != nil){
		logText("There was an error gathering new fuel prices", true);
		panic(err)
	}else {
		logText("Fuel Updates were successful. \nNext update in "+strconv.Itoa(eiaUpdateIntervalDays)+" days...", false);
		defer response.Body.Close();
		time.AfterFunc((time.Duration(eiaUpdateIntervalDays)*24*time.Hour), triggerEIAFuelPrices);
	}
}

/******************************************************************************************
* 
*	Utility Methods 
*
*******************************************************************************************/

func logText(sToLog string, isError bool) bool{
    var logFile string
    if isError {
        logFile = errorLog
    } else {
        logFile = commandLog
    }

    file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Printf("ERROR: Could not write to log file: %v", err)
        return false
    }
    defer file.Close()

    logger := log.New(file, "", log.LstdFlags)
    logger.Println(sToLog)
    return true
}