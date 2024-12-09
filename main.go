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
*	Configuration 
*
********************/


/**
* Initial variables for both the environment and the flags. 
* Also went ahead and added the client and triggers since they'll
* be used throughout the program.
*/
var (
    errorLog               string
    commandLog             string
    client                 = &http.Client{}
	debugMode			   bool
	warmUp				   bool
	triggers			   []Trigger
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
		logText("There was an error opening the triggers file. Does it exist?", tErr);
		log.Fatal(tErr);
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
	/**
	* Start the trigger processing loop. The triggers themselves will 
	* handle their own intervals.
	*/
	for _, trigger := range triggers {
		processTrigger(trigger, warmUp);
	}

	select {};
}