package main

import (
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"os"
	"strconv"
)


func store(sdk *fabsdk.FabricSDK, channelID string, userName string, ccID string) {
	context := sdk.ChannelContext(channelID, fabsdk.WithUser(userName))
	// 通道客户端
	cClient, _ := channel.New(context)
	// 账本客户端
	lClient, _ := ledger.New(context)
	// 事件客户端
	eClient, _ := event.New(context)

	reg, notifier, _ := eClient.RegisterChaincodeEvent(ccID, "event1")
	defer eClient.Unregister(reg)

	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data :add key-value")
	response, _ := cClient.Execute(channel.Request{
		ChaincodeID:     ccID,
		Fcn:             "invoke",
		Args:            [][]byte{[]byte("store"), []byte("{\"id\":\"1\",\"name\":\"zs\"}")},
		TransientMap: transientDataMap,
	})

	select {
	case ccEvent := <-notifier:
		fmt.Println("|#############################received cc event:",ccEvent)
	}
	txid := response.TransactionID
	block, _ := lClient.QueryBlockByTxID(txid)
	blockHeight := strconv.Itoa(int(block.Header.GetNumber()))
	blockHash := hex.EncodeToString(block.Header.GetDataHash())
	fmt.Println("|#############################", txid, blockHeight, blockHash)
}

func query(sdk *fabsdk.FabricSDK, channelID string, userName string, chaincodeID string) {

	context := sdk.ChannelContext(channelID, fabsdk.WithUser(userName))
	// 通道客户端
	cClient, _ := channel.New(context)

	response, _ := cClient.Query(channel.Request{
		ChaincodeID: chaincodeID,
		Fcn: "invoke",
		Args: [][]byte{[]byte("query"), []byte("keyid")},
	})
	fmt.Println("|#############################", response.Payload)
}

func xx() {
	sdk, _ := fabsdk.New(config.FromFile("config1.yaml"))
	rClient, _ := resmgmt.New(sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1")))


	// 创建通道
	mspClient, _ := mspclient.New(sdk.Context(), mspclient.WithOrg("org1"))
	adminIdentity, _ := mspClient.GetSigningIdentity("Admin")
	req := resmgmt.SaveChannelRequest{ChannelID: "ch1", ChannelConfigPath: "fabric/channel/channel.tx", SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	cresp, _ := rClient.SaveChannel(req, resmgmt.WithOrdererEndpoint("ord1.oorg1.com"))
	fmt.Println("|#############################create channel txid:", cresp.TransactionID)
	// 创建通道交易id:294129f2c804fd170f6596ff09d630f600a1441a99a312e701061587710a3983

	// 加入通道
	rClient.JoinChannel("ch1", resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("ord1.oorg1.com"))
	//block height=1

	// 安装链码
	// sdk安装all peer
	ccPkg, _ := packager.NewCCPackage("myfabric/chaincode/", os.Getenv("GOPATH"))
	ccReq := resmgmt.InstallCCRequest{
		Name: "cc1",
		Path: "myfabric/chaincode/",
		Version: "0",
		Package: ccPkg,
	}
	rClient.InstallCC(ccReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//ccid:b03fc0570396d008c5b97ba1bade048c70faa0ecc5b8778675e097052384ee08
	//just installed for peers of an org
	//must install for all orgs


	// 实例化链码
	ccPolicy := policydsl.SignedByAnyMember([]string{"org1MSP"})
	resp, _ :=rClient.InstantiateCC("ch1", resmgmt.InstantiateCCRequest{
		Name: "cc1",
		Path: os.Getenv("GOPATH"),
		Version: "0",
		Args: [][]byte{[]byte("init")},
		Policy: ccPolicy,
	})
	fmt.Println("|#############################实例化链码txid:", resp.TransactionID)
	// 实例化链码txid:18960d43615f2e0b8861f53b370a5f9c31cc5eabcad3a688cc58533ffc2d4ed5
	// block height=2


}


func main() {

	sdk, _ := fabsdk.New(config.FromFile("config1.yaml"))
	_, _ = resmgmt.New(sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1")))

	context := sdk.ChannelContext("ch1", fabsdk.WithUser("User1"))
	// 通道客户端
	_, err := channel.New(context)

	if err != nil {
		fmt.Println(err)
	}

	defer sdk.Close()

	// controllers.Init(&sdk)
	// beego.Run()
}
