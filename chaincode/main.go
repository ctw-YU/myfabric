package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Data struct {
	Id     		string   `json:"id"`
	Name       	string   `json:"name"`
}
type MyChaincode struct {

}

func main() {
	shim.Start(new(MyChaincode))
}

func (cc *MyChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("chaincode instantiated succefully"))
}

func (cc *MyChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	if function != "invoke" {
		shim.Error("unknown function!")
	}

	if args[0] == "query" {
		keyId := args[1]
		data, _ := stub.GetState(keyId)
		return shim.Success(data)
	}

	if args[0] == "store" {
		var reqData Data
		json.Unmarshal([]byte(args[1]), &reqData)
		data, _ := json.Marshal(reqData)
		stub.PutState(reqData.Id, data)
		stub.SetEvent("event1", []byte("store success"))
		return shim.Success(nil)
	}

	return shim.Error("unknown function!")
}
