package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SampleChaincode struct {
}

var logger = shim.NewLogger("chaincode")

type DataDigest struct {
	KeyId              string   `json:"keyId"`
	Vserion            string   `json:"version"`
	UserName           string   `json:"username"`
	OperationType      string   `json:"operationType"`
	DataType           string   `json:"dataType"`
	ServiceType        string   `json:"serviceType"`
	FileName           string   `json:"fileName"`
	FileSize           string   `json:"fileSize"`
	FileHash           string   `json:"fileHash"`
	URI                string   `json:"uri"`
	ParentKeyId        string   `json:"parentKeyId"`
	AttachmentFileUris []string `json:"attchmentFileNames"`
	AttachmentTotalHash string `json:"attachmentTotalHash"`
}
type VerifyRequest struct {
	KeyId string `json:"keyId""`
	Hash  string `json:"hash""`
}

type Response struct {
	Status string `json:"status""`
	Err    string `json:"err""`
}

type BlockInfoResp struct {
	Status      string `json:"status""`
	Err         string `json:"err""`
	TxId        string `json:"txId""`
	BlockHeight string `json:"blockHeight""`
	BlockHash   string `json:"blockHash""`
}

// So the goal is to prepare the ledger to handle future requests.
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### Chaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke 匹配方法字符串
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### Chaincode Invoke ###########")

	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	//根据KeyId查询
	if args[0] == "query" {
		return t.query(stub, args)
	}

	//结构体存至链
	if args[0] == "storeDigest" {
		return t.storeDigest(stub, args)
	}

	/*
		//根据结构体 验证是否在链上
		if args[0] == "verifyFile" {
			return t.verifyFile(stub,args)
		}
	*/

	return shim.Error("Unknown action, check the first argument")
}

//根据结构体 验证
func (t *SampleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Info("########### HeroesServiceChaincode 查询数据 ###########")
	var argumentKeyId string //存储传过来的 待验证的KeyId
	argumentKeyId = args[1]
	if argumentKeyId != "" {
		byteArrDatadigest, errGetDigest := stub.GetState(argumentKeyId)

		if errGetDigest != nil {
			logger.Info("###########查询数据 - byteArrDataDigest 结构为空 ###########")
			return shim.Error("###########查询数据 - byteArrDataDigest为空 ###########")

		}

		return shim.Success(byteArrDatadigest)

	} else {
		logger.Info("###########查询数据 - KeyId 为空 ###########")
		return shim.Error("###########查询数据 - KeyId 为空 ###########")
	}
}

// 存证数据
func (t *SampleChaincode) storeDigest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Info("########### HeroesServiceChaincode 存证数据 ###########")

	var argumentDataDigest DataDigest
	errUnmarshal := json.Unmarshal([]byte(args[1]), &argumentDataDigest)
	if errUnmarshal != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	byteArrDataDigest, erradd := json.Marshal(argumentDataDigest)
	if erradd != nil {
		return shim.Error("序列化数据失败")
	}
	errPutState := stub.PutState(argumentDataDigest.KeyId, byteArrDataDigest)
	if errPutState != nil {
		return shim.Error("更新数据失败")
	}
	// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
	err := stub.SetEvent("eventInvoke", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return this value in response
	return shim.Success(nil)
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

