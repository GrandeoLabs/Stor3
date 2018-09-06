package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct{}

// share = asset : buildings , director's contract etc
type ChainKey struct {
	DocumentID     string `json: dID`  //internal
	Glue           string `json: glue` // array of objects
	TypeofDocument string `json: typeofdocument`
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Init Method is called")
	s.initLedger(stub)
	return shim.Success(nil)
}

// initialising ledger  -  only for the purpose for testing
func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Initalising ledger")
	chainkeys := []ChainKey{
		ChainKey{DocumentID: "doc1", Glue: "glue1", TypeofDocument: "type1"},
		ChainKey{DocumentID: "doc2", Glue: "glue2", TypeofDocument: "type2"},
	}

	m := 0
	for m < len(chainkeys) {
		var chainkeysoutput = chainkeys[m]
		chainkeyAsBytes, _ := json.Marshal(&chainkeysoutput)
		stub.PutState(strconv.Itoa(m+1), chainkeyAsBytes)
		fmt.Println("Added", chainkeys[m])
		m = m + 1

	}
	return shim.Success(nil)
}

// Invoke function - every time it is called, will check for the function name and arguments
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	functionname, args := stub.GetFunctionAndParameters()
	if functionname == "queryChainKey" {
		return s.queryChainKey(stub, args)
	} else if functionname == "listChainKeys" {
		return s.listChainKeys(stub)
	} else if functionname == "recordChainKeys" {
		fmt.Println("Recording hash...")
		return s.recordChainKey(stub, args)
	} else {
		fmt.Println("Invalid function name")
	}
	return shim.Success(nil)
}

func (s *SmartContract) recordChainKey(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error(" Incorrect number of arguments - please enter 10 values")
	}

	fmt.Println("we are trying to print the list of contracts .....")
	var chainkeys = ChainKey{DocumentID: args[0], Glue: args[1], TypeofDocument: args[2]}
	fmt.Println(chainkeys)
	chainkeyAsBytes, _ := json.Marshal(chainkeys)

	err := stub.PutState(args[0], chainkeyAsBytes)

	if err != nil {
		fmt.Println(" Error Occured whilst saving the record", err)

	}
	return shim.Success(chainkeyAsBytes)
}

// ListAssets -  list all the asset info
func (s *SmartContract) listChainKeys(stub shim.ChaincodeStubInterface) sc.Response {
	startKey := "0"
	endKey := "999999"
	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		fmt.Printf("We are trying to print raw data")
		fmt.Printf(string(queryResponse.Value))
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members, supress it for first array memeber
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(",\"Record\":")
		// Record is a JSON object so we can re-write it as

		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("List of Contracts: \n%s \n", buffer.String())
	return shim.Success(buffer.Bytes())
}

// Query Contract Info - query for specifc constract based on "contract ID"
func (s *SmartContract) queryChainKey(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		shim.Error("Incorrect number of arguments, Expecting 1")
	}
	chainkeyAsBytes, _ := stub.GetState(args[0])
	data := string(chainkeyAsBytes)
	fmt.Println("Returning list of Asset %s args[0]", data)

	if chainkeyAsBytes == nil {
		return shim.Error("could not locate Asset")
	}
	return shim.Success(chainkeyAsBytes)
}

// main function
func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating a new smart contract %s", err)
	}
}
