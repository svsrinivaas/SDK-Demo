/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	//"fmt"
	//"strconv"
	//"encoding/json"
	//"bytes"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("mylogger")

// SimpleChaincode aduti trail Chaincode implementation
type SimpleChaincode struct {
}

//log data model
type LogInfo struct {
	userId 		string `json:"userId"`
	operation  	string `json:"operation"`
	desc     	string `json:"desc"`
	createdTime	string `json:"time"`
}

var EVENT_COUNTER = "event_counter"
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Info("Successfully initialized log chain") 
	return nil, nil
}

// Transaction makes an entry of audit log 
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2 (userId, log info) ")
	}

	// Initialize the chaincode
	var user = args[0]
	var logInfo = args[1]
	
	err := stub.PutState(user, []byte(logInfo)) 
	if err != nil { 
		logger.Error("Could not save log details to ledger", err) 
		return nil, err 
	} 
		
	logger.Info("Successfully saved log details") 
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	logger.Debug("Entering Get Log Details")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments, expected userId")
		return nil, errors.New("Missing user Id")
	}

	var userId = args[0]
	bytes, err := stub.GetState(userId)
	if err != nil {
		logger.Error("Could not fetch log details with id " + userId + " from ledger", err)
		return nil, err
	}
	return bytes, nil
}
	
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Error("Could not start SampleChaincode")
	} else {
		logger.Info("SampleChaincode successfully started")
	}
}