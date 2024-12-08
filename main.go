package main

import (
	"net/http"
	"encoding/json"
	"flag"
	"os"
	"log"
	"io/ioutil"
    "github.com/joho/godotenv"
)

/*******************
* 
	Configuration 
*
********************/


var (
    erpKey                 string
    errorLog               string
    commandLog             string
    trackingUpdateInterval string
    eiaUpdateIntervalDays  string
    client                 = &http.Client{}
	debugMode			   bool
	warmUp				   bool
	triggers			   []Trigger
)

const (
	redColor    = "\033[31m"
	yellowColor = "\033[33m"
	resetColor  = "\033[0m"
)

func init() {
	/**
	* Command line flags. 
	*/
	flag.BoolVar(&debugMode, "debugMode", false, "Enable Debug Mode");
	flag.BoolVar(&warmUp, "warmUp", true, "Fire triggers on startup when true. Otherwise, wait for the first interval.");
	
	flag.Parse();

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file");
	}

	errorLog = os.Getenv("ERROR_LOG");
	commandLog = os.Getenv("COMMAND_LOG");

	/**
	* Load the triggers from the triggers.json file.
	*/
	tFile, tErr := os.OpenFile("triggers.json", os.O_RDONLY, 0444);
	if tErr != nil {
		panic(tErr);
	}
	defer tFile.Close();

	/**
	* Decode the JSON contents within the file. 
	*/

	byteValue, jErr := ioutil.ReadAll(tFile);
	if(jErr != nil){
		logText("There was an error reading the triggers file. Are permissions set properly?", jErr);
		log.Fatal(jErr);
	}

	/**
	* Put these triggers into the triggers slice.
	*/
	juErr := json.Unmarshal(byteValue, &triggers);
	if(juErr != nil){
		logText("The triggers file could not be processed. Is the format correct?", juErr);
	}

	/**
	* Output debugging information using functions in the util file.
	*/
	if debugMode {
		outputEnvInfo();
		outputTriggers(triggers);
	}
}

func main() {
	for _, trigger := range triggers {
		processTrigger(trigger, warmUp);
	}

	select {};
}