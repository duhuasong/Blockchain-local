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
	Device_Person string
	Event_type string
}

func main() {
	err := shim.Start(new(myChaincode))
	if err != nil {
		fmt.Printf("Error at initialization")
	}
}

func (t *myChaincode) Init (stub shim.ChaincodeStubInterface) peer.Response {	
	err := stub.PutState("metric_count",[]byte("0"))
	if err != nil {
		return shim.Error("Failed saving start state of metric_count")
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
	if fn == "add_metric" {
		result, err = add_metric(stub,args)		
	} else if fn == "metric_count" {
		result, err = metric_count(stub)
	} else if fn == "raw_metrics" {
		result, err = raw_metrics(stub, args)
	} else if fn == "search_by_user" {
		result, err = search_by_user(stub, args)
	} else {
		return shim.Error("Unknown command")
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}

func add_metric (stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("Missing required parameters")
	}

	idAsBytes,err1 := stub.GetState("metric_count")
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
		Device_Person: args[1],
		Event_type: args[2]
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
	err8 := stub.PutState("metric_count",[]byte(idstring))
	if err8 != nil {
		return nil, fmt.Errorf("Failed updating metric_count")
	}
	
	return nil,nil
}

func metric_count(stub shim.ChaincodeStubInterface)([]byte,error) {
	bytes,err := stub.GetState("metric_count")
	if err != nil {
		return nil, fmt.Errorf("Failed reading metric_count")
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

func search_by_user (stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Person and number of records required")
	}

	number, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("Number is in invalid format")
	}

	bytes, err2 := stub.GetState("metrics")
	if err2 != nil {
		return nil, fmt.Errorf("Failed reading metrics")
	}

	var metrics []Metric
	json.Unmarshal(bytes,&metrics)

	//reverse array
	for i, j := 0, len(metrics)-1; i < j; i, j = i+1, j-1 {
        metrics[i], metrics[j] = metrics[j], metrics[i]
	}
	
	//find specific 
	var ret []Metric
	for index, element := range metrics {
		if element.Device_Person == args[0] {
			ret = ret.append(element)
			if len(ret) == number {
				break;
			}
		}
	}

	var final, err3 = json.Marshal(ret)
	if err3 != nil {
		return nil, fmt.Errorf("Failed marshaling array")
	}

	return final, nil
}