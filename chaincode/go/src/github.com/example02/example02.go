package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		return t.invoke(stub, args)
	}  else if function == "query" {
		return t.query(stub, args)
	}  else if function == "put" {
		return t.putStandard(stub, args)
	}  else if function == "get" {
		return t.getStandard(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string
	var opA,opB string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]
	opA = "-"
	opB = "+"

	_, err = strconv.ParseFloat(args[3], 64)
	if err != nil {
		return shim.Error("Provided value was not a number")
	}

	txid := stub.GetTxID()
	compositeIndexName := "varName~op~value~txID"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, []string{A, opA, args[3], txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", A, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", A, compositePutErr.Error()))
	}

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey1, compositeErr1 := stub.CreateCompositeKey(compositeIndexName, []string{B, opB, args[3], txid})
	if compositeErr1 != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", B, compositeErr.Error()))
	}

	compositePutErr1 := stub.PutState(compositeKey1, []byte{0x00})
	if compositePutErr1 != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", B, compositePutErr.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Successfully invoke %s from %s to %s", args[3], A, B)))

	
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check we have a valid number of args
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}

	name := args[0]
	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := stub.GetStateByPartialCompositeKey("varName~op~value~txID", []string{name})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", name, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the name %s exists", name))
	}
	v , err := stub.GetState(name)

	// Iterate through result set and compute final value
	var finalVal float64
	finalVal, err = strconv.ParseFloat(string(v), 64)
	if err != nil {
		return shim.Error("Provided value was not a number")
	}
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := stub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}

		// Retrieve the delta value and operation
		operation := keyParts[1]
		valueStr := keyParts[2]

		// Convert the value string and perform the operation
		value, convErr := strconv.ParseFloat(valueStr, 64)
		if convErr != nil {
			return shim.Error(convErr.Error())
		}

		switch operation {
		case "+":
			finalVal += value
		case "-":
			finalVal -= value
		default:
			return shim.Error(fmt.Sprintf("Unrecognized operation %s", operation))
		}
	}

	return shim.Success([]byte(strconv.FormatFloat(finalVal, 'f', -1, 64)))
}

func (s *SimpleChaincode) putStandard(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	name := args[0]
	valStr := args[1]

	_, getErr := stub.GetState(name)
	if getErr != nil {
		return shim.Error(fmt.Sprintf("Failed to retrieve the statr of %s: %s", name, getErr.Error()))
	}

	putErr := stub.PutState(name, []byte(valStr))
	if putErr != nil {
		return shim.Error(fmt.Sprintf("Failed to put state: %s", putErr.Error()))
	}

	return shim.Success(nil)
}

func (s *SimpleChaincode) getStandard(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	name := args[0]

	val, getErr := stub.GetState(name)
	if getErr != nil {
		return shim.Error(fmt.Sprintf("Failed to get state: %s", getErr.Error()))
	}

	return shim.Success(val)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

