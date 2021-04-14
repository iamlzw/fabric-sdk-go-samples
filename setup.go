package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/common/tools/protolator"
	//"github.com/hyperledger/fabric/protos/utils"
	cb "github.com/hyperledger/fabric/protos/common"
	"io/ioutil"
	"log"
	"strconv"
)


var defaultQueryArgs = [][]byte{[]byte("b")}

var defaultInitCCArgs = [][]byte{[]byte("Init"),[]byte("a"),[]byte("100"),[]byte("b"),[]byte("200")}

//init the sdk
func initSDK() *fabsdk.FabricSDK {
	//// Initialize the SDK with the configuration file
	configProvider := config.FromFile("config_e2e.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Errorf("failed to create sdk: %v", err)
	}

	return sdk
}

func createChannel(sdk *fabsdk.FabricSDK) {
	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("OrdererOrg"))
	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		fmt.Println("Failed to create channel management client: %s", err)
	}
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg("Org1"))
	if err != nil {
		fmt.Println(err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity("Admin")
	if err != nil {
		fmt.Println(err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: "mychannel",
		ChannelConfigPath: "/home/www/go/src/github.com/hyperledger/fabric-samples/first-network/channel-artifacts/channel.tx",
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))
	fmt.Println(txID)
}

func joinChannel(sdk *fabsdk.FabricSDK){
	//prepare context
	adminOrg1Context := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	// Org resource management client
	org1ResMgmt, err := resmgmt.New(adminOrg1Context)
	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}

	// Org peers join channel
	if err = org1ResMgmt.JoinChannel("mychannel", resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com")); err != nil {
		fmt.Println("Org peers failed to JoinChannel: %s", err)
	}

	adminOrg2Context := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org2"))

	// Org resource management client
	org2ResMgmt, err := resmgmt.New(adminOrg2Context)
	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}

	// Org peers join channel
	if err = org2ResMgmt.JoinChannel("mychannel", resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com")); err != nil {
		fmt.Println("Org peers failed to JoinChannel: %s", err)
	}
}

func createCC(sdk *fabsdk.FabricSDK) {

	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}
	ccPkg, err := packager.NewCCPackage("github.com/example_cc", "/home/www/go")
	if err != nil {
		fmt.Println(err)
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: "mycc", Path: "github.com/example_cc", Version: "1.0", Package: ccPkg}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}
	// Set up chaincode policy
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := orgResMgmt.InstantiateCC(
		"mychannel",
		resmgmt.InstantiateCCRequest{Name: "mycc", Path: "github.com/example_cc", Version: "1.0", Args: defaultInitCCArgs, Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	fmt.Println(resp.TransactionID)
}

func queryLedger(sdk *fabsdk.FabricSDK){
	ccp := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}
	log.Println("Initialized channel client")

	ledgerClient, err := ledger.New(ccp)

	// Test Query Info - retrieve values before transaction
	chainInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		fmt.Println("QueryInfo return error: %s", err)
	}
	height := chainInfo.BCI.Height
	for i := int64(0); i < int64(height); i++ {
		block ,err := ledgerClient.QueryBlock(uint64(i))
		if err != nil{
			fmt.Println(err)
		}
		buf := new (bytes.Buffer)
		err = protolator.DeepMarshalJSON(buf, block)
		if err != nil {
			fmt.Errorf("malformed block contents: %s", err)
		}
		filename := "blockfiles/mychannel_"+strconv.FormatInt(i,10)+".json"
		err = ioutil.WriteFile(filename,buf.Bytes(),0644)
		if err != nil{
			fmt.Println("write to file failure:",err)
		}

	}
	response, err := cc.Query(channel.Request{ChaincodeID: "mycc", Fcn: "query", Args: [][]byte{[]byte("a")}},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(),
	)
	if err != nil {
		fmt.Println("Failed to query funds: %s", err)
	}

	txResponse := response.Responses
	endorsement := txResponse[0].Endorsement
	endorser := endorsement.Endorser
	signature := endorsement.Signature

	fmt.Println(string(endorser))
	fmt.Println(string(signature))
}

func downloadBlock(sdk *fabsdk.FabricSDK){
	ccp := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))

	ledgerClient, err := ledger.New(ccp)

	// Test Query Info - retrieve values before transaction
	chainInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		fmt.Println("QueryInfo return error: %s", err)
	}
	height := chainInfo.BCI.Height
	for i := int64(0); i < int64(height); i++ {
		block ,err := ledgerClient.QueryBlock(uint64(i))
		if err != nil{
			fmt.Println(err)
		}
		b,err := proto.Marshal(block)
		//buf := new (bytes.Buffer)
		//err = protolator.DeepMarshalJSON(buf, block)
		if err != nil {
			fmt.Errorf("malformed block contents: %s", err)
		}
		filename := "blockfiles/mychannel_"+strconv.FormatInt(i,10)+".block"
		err = ioutil.WriteFile(filename,b,0644)
		if err != nil{
			fmt.Println("write to file failure:",err)
		}
	}
}

func parseBlock(){
	file := "blockfiles/mychannel_3.block"
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Errorf("Could not read block")
	}

	block := &cb.Block{}
	err = proto.Unmarshal(data, block)
	if err != nil {
		fmt.Errorf("error unmarshaling to block: %s", err)
	}
	buf := new (bytes.Buffer)
	err = protolator.DeepMarshalJSON(buf, block)
	if err != nil {
		fmt.Errorf("malformed block contents: %s", err)
	}

	if err != nil {
		fmt.Errorf("malformed block contents: %s", err)
	}
	filename := "mychannel_4"+".json"
	err = ioutil.WriteFile(filename,buf.Bytes(),0644)
	if err != nil{
		fmt.Println("write to file failure:",err)
	}
}




