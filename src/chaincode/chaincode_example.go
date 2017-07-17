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
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
var EVENT_COUNTER = "event_counter"
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var user, operation, desc, time string    // Entities
	//var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	user = args[0]
	operation = args[1]
	desc = args[2]
	time = args[3]
	
	fmt.Printf("user = %s, operation = %s, desc = %s, time = %s \n", user, operation, desc, time)

	// Write the state to the ledger
	/*ok, err = stub.InsertRow(user, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: user}},
			&shim.Column{Value: &shim.Column_String{String_: operation}},
			&shim.Column{Value: &shim.Column_String_{String_: desc}},
			&shim.Column{Value: &shim.Column_String{String_: time}}},
	})
    
	err = stub.PutState(EVENT_COUNTER, []byte("1"))
	if err != nil {
		return nil, err
	}*/
	return nil, nil
}

// Transaction makes an entry of audit log 
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var user, operation, desc, time string    // Entities
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	user = args[0]
	operation = args[1]
	desc = args[2]
	time = args[3]

	fmt.Printf("user = %s, operation = %s, desc = %s, time = %s \n", user, operation, desc, time)
	
	// Write the state to the ledger
	ok, err := stub.InsertRow(user, shim.Row{
		Columns: []*shim.Column{
			&shim.Column{Value: &shim.Column_String_{String_: user}},
			&shim.Column{Value: &shim.Column_String_{String_: operation}},
			&shim.Column{Value: &shim.Column_String_{String_: desc}},
			&shim.Column{Value: &shim.Column_String_{String_: time}}},
	})

	if !ok && err == nil {
		return nil, errors.New("User log already exists.")
	}

	fmt.Printf("Invoke...done!")
	
	//Event based
    b, err := stub.GetState(EVENT_COUNTER)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	noevts, _ := strconv.Atoi(string(b))

	tosend := "Event Counter is " + string(b)

	err = stub.PutState(EVENT_COUNTER, []byte(strconv.Itoa(noevts+1)))
	if err != nil {
		return nil, err
	}

	err = stub.SetEvent("evtsender", []byte(tosend))
	if err != nil {
		return nil, err
        }
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var userName string
	//var user, opr, desc, time string // Entities
	var err error
	userName = args[0]

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	var columns []shim.Column
	row, err := stub.GetRow(userName, columns)
	if err != nil {
		return nil, fmt.Errorf("Failed retrieving log for [%s]: [%s]", userName, err)
	}

	user := row.Columns[0].GetBytes()
	opr := row.Columns[1].GetBytes()
	desc := row.Columns[2].GetBytes()
	time := row.Columns[3].GetBytes()
	
	fmt.Printf("row value : %$, %$, %$, %$", user, opr, desc, time)

	//jsonResp := "{\"user\":\"" + user + "\",\"operation\":\"" + opr + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	return user, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}