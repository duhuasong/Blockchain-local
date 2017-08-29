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
	Image string
	Type int
	Timestamp int
}

func main() {
	err := shim.Start(new(myChaincode))
	if err != nil {
		fmt.Printf("Error at initialization")
	}
}

func (t *myChaincode) Init (stub shim.ChaincodeStubInterface) peer.Response {
	err1 := stub.PutState("test",[]byte("It's alive!"))
	if err1 != nil {
		return shim.Error("Network test failed")
	}
	
	err3 := stub.PutState("pictures_taken",[]byte("0"))
	if err3 != nil {
		return shim.Error("Failed saving start state of picture_taken")
	}
	
	var met []Metric
	var bytes,err4 = json.Marshal(met)
	if err4 != nil {
		return shim.Error("Failed converting to JSON")
	}
	
	err5 := stub.PutState("metrics",bytes)
	if err5 != nil {
		return shim.Error("Failed saving start state of metrics")
	}
	
	return shim.Success(nil)
}

func (t *myChaincode) Invoke (stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result []byte
	var err error
	if fn == "add_picture" {
		result, err = add_picture(stub,args)		
	} else if fn == "picture_count" {
		result, err = picture_count(stub)
	} else if fn == "raw_metrics" {
		result, err = raw_metrics(stub, args)
	} else {
		return shim.Error("Unknown command")
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func add_picture (stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("Missing required parameters")
	}

	idAsBytes,err1 := stub.GetState("pictures_taken")
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
	
	timestamp, err5 := strconv.Atoi(args[1])
	if err5 != nil {
		return nil, fmt.Errorf("Provided timestamp is invalid")
	}

	typea, err10 := strconv.Atoi(args[2])
	if err10 != nil {
		return nil, fmt.Errorf("Provided type is not int")
	}
	
	metric := Metric {
		Id: id,
		Image: args[0],
		Type: typea,
		Timestamp: timestamp,
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
	err8 := stub.PutState("pictures_taken",[]byte(idstring))
	if err8 != nil {
		return nil, fmt.Errorf("Failed updating pictures_taken")
	}
	
	return nil,nil
}

func picture_count(stub shim.ChaincodeStubInterface)([]byte,error) {
	bytes,err := stub.GetState("pictures_taken")
	if err != nil {
		return nil, fmt.Errorf("Failed reading pictures_taken")
	}
	return bytes,nil
}

func raw_metrics (stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
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
