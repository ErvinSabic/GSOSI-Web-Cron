package main

import (
	"log"
	"os"
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
func outputDebugInfo(){
	debugInfo := 
		"_______________________________________ \n" +
		yellowColor+"Debug Mode Enabled\n" +
		"Environment Variables Gathered from .env \n" +resetColor+
		"_______________________________________ \n" +
		"Server Port: " + serverPort + "\n" +
		"End Point: " + endPoint + "\n" +
		"ERP Key: " + erpKey + "\n" +
		"Error Log: " + errorLog + "\n" +
		"Command Log: " + commandLog + "\n" +
		"Tracking Update Interval: " + trackingUpdateInterval + "\n" +
		"EIA Update Interval: " + eiaUpdateIntervalDays + "\n" +
		"________________________________________\n";
	logText(debugInfo, nil);	
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