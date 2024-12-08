package main

import (
	"log"
	"os"
    "github.com/joho/godotenv"

//	"strconv"
)
/******************************************************************************************
* 
*	Utility Methods 
*
*******************************************************************************************/

/**
* This function is used to output debug information to the console and log file.
* When the --debugMode=true flag is set.
*/
func outputEnvInfo(){
	debugInfo := 
		"_______________________________________ \n" +
		yellowColor+"Debug Mode Enabled\n" +
		"Environment Variables Gathered from .env \n" +resetColor+
		"_______________________________________ \n";
    
    // Read the .env file. This will return a map of key value pairs.
    envVars, err := godotenv.Read(".env");
    if(err != nil){
        logText("There was an error reading the .env file", err);
        return
    }
    // Append every key value pair within the .env file. 
    for key, value := range envVars {
        debugInfo += key + ": " + value + "\n";
    }

    // Output to logger.
	logText(debugInfo, nil);	
}

/**
* This function will output the value of the triggers slice to the console and log file.
*/
func outputTriggers(triggers []Trigger){
    for _, trigger := range triggers {
        triggerInfo := 
            "_______________________________________ \n" +
            yellowColor+"Trigger Info \n" + resetColor +
            "Name: " + trigger.Name + "\n" +
            "Endpoint: " + trigger.Endpoint + "\n" +
            "Method: " + trigger.Method + "\n" +
            "Route: " + trigger.Route + "\n" +
            "Duration: " + trigger.Duration + "\n" +
            "Additional Headers: \n";
        for _, header := range trigger.AdditionalHeaders {
            triggerInfo += "Key: " + header.Key + " Value: " + header.Value + "\n";
        }
        triggerInfo += "_______________________________________ \n";
        logText(triggerInfo, nil);
    }
}

/**
* This function is used to log text to the console and log file. Console 
* output is only available when --debugMode=true is set.
*/
func logText(sToLog string, err error) bool {
    if debugMode {
        if err != nil {
            print(redColor + sToLog + "\n" + err.Error() + resetColor + "\n")
        } else {
            print(sToLog + "\n")
        }
    }

    var logFile string
    if err != nil {
        logFile = errorLog
    } else {
        logFile = commandLog
    }

    file, fErr := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if fErr != nil {
        log.Printf("ERROR: Could not write to log file: %v", fErr)
        return false
    }
    defer file.Close()

    logger := log.New(file, "", log.LstdFlags)
    logger.Println(sToLog)
	if(err != nil){
		logger.Println(err.Error());
	}
    return true
}