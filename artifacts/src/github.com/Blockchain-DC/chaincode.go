package main

import (
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type myChaincode struct {}

type Metric struct {
	Id int
	Timestamp int
	Id_slike string
	Event_type string
}

func main() {
	err := shim.Start(new(myChaincode))
	if err != nil {
		fmt.Printf("Error at initialization")
	}
}

func (t *myChaincode) Init (stub shim.ChaincodeStubInterface) peer.Response {	
	err := stub.PutState("event_count",[]byte("0"))
	if err != nil {
		return shim.Error("Failed saving start state of event_count")
	}
	
	var met []Metric
	var bytes,err2 = json.Marshal(met)
	if err2 != nil {
		return shim.Error("Failed converting to JSON")
	}
	
	err3 := stub.PutState("metrics",bytes)
	if err3 != nil {
		return shim.Error("Failed saving start state of metrics")
	}
	
	return shim.Success(nil)
}

func (t *myChaincode) Invoke (stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result []byte
	var err error
	if fn == "add_event" {
		result, err = add_event(stub,args)		
	} else if fn == "event_count" {
		result, err = event_count(stub)
	} else if fn == "raw_events" {
		result, err = raw_events(stub, args)
	} else {
		return shim.Error("Unknown command")
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func add_event (stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("Missing required parameters")
	}

	idAsBytes,err1 := stub.GetState("event_count")
	if err1 != nil {
		return nil, fmt.Errorf("Failed getting id")
	}
	
	id,err2 := strconv.Atoi(string(idAsBytes))
	if err2 != nil {
		return nil, fmt.Errorf("Failed converting from string to integer")
	}
	
	metricsAsBytes,err3 := stub.GetState("metrics")
	if err3 != nil {
		return nil, fmt.Errorf("Failed getting metrics")
	}
	
	var met []Metric
	err4 := json.Unmarshal(metricsAsBytes, &met)
	if err4 != nil {
		return nil, fmt.Errorf("Failed unmarshaling JSON")
	}
	
	timestamp, err5 := strconv.Atoi(args[0])
	if err5 != nil {
		return nil, fmt.Errorf("Provided timestamp is invalid")
	}
	
	metric := Metric {
		Id: id,
		Timestamp: timestamp,
		Id_slike: args[1],
		Event_type: args[2],
	}

	met = append(met,metric)
	bytes,err5 := json.Marshal(met)
	if err5 != nil {
		return nil, fmt.Errorf("Failed marshaling JSON")
	}
	
	err6 := stub.PutState("metrics",bytes)
	if err6 != nil {
		return nil, fmt.Errorf("Failed updating metrics")
	}
	
	id++
	idstring := strconv.Itoa(id)
	err8 := stub.PutState("event_count",[]byte(idstring))
	if err8 != nil {
		return nil, fmt.Errorf("Failed updating event_count")
	}
	
	return nil,nil
}

func event_count(stub shim.ChaincodeStubInterface)([]byte,error) {
	bytes,err := stub.GetState("event_count")
	if err != nil {
		return nil, fmt.Errorf("Failed reading event_count")
	}
	return bytes,nil
}

func raw_events (stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Missing number of results")
	}
	
	nr, er := strconv.Atoi(args[0])
	if er != nil {
		return nil, fmt.Errorf("Number required")
	}
	
	bytes,err := stub.GetState("metrics")
	if err != nil {
		return nil, fmt.Errorf("Failed reading metrics")
	}
	
	var metrics []Metric
	json.Unmarshal(bytes,&metrics)
	
	var count = len(metrics)
	if count < nr {
		//return nil, fmt.Errorf("You can't get more entries that there are")
		//Now returns all entries if requested number is greater then number of entries
	} else {
		metrics = metrics[count-nr:]		
	}	
	
	var final, err1 = json.Marshal(metrics)
	if err1 != nil {
		return nil,err1
	}
	
	return final,nil
}
